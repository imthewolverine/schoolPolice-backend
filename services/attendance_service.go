package services

import (
	"context"
	"fmt"
	"log"
	"time"

	"cloud.google.com/go/firestore"
)

// AttendanceService provides functions to work with Attendance documents
type AttendanceService struct {
    FirestoreClient *firestore.Client
}
type Account struct {
    Balance int `firestore:"balance"`
}

// UpdateStartTime updates the StartTime of an Attendance document to the current time
func (s *AttendanceService) UpdateStartTime(ctx context.Context, attendanceID string) error {
    attendanceRef := s.FirestoreClient.Collection("Attendance").Doc(attendanceID)
    
    _, err := attendanceRef.Update(ctx, []firestore.Update{
        {Path: "StartTime", Value: time.Now().Format(time.RFC3339)}, // Use current time
    })
    
    return err
}

// UpdateEndTime updates the EndTime to the current time, calculates the total time spent,
// and transfers salary from the parent's account to the worker's account.
func (s *AttendanceService) UpdateEndTime(ctx context.Context, attendanceID string) error {
    log.Println("Starting UpdateEndTime function")
    
    attendanceRef := s.FirestoreClient.Collection("Attendance").Doc(attendanceID)

    // Get the current time for EndTime
    endTime := time.Now()

    // Fetch the existing StartTime, Ad reference, and Request reference from the document
    doc, err := attendanceRef.Get(ctx)
    if err != nil {
        log.Printf("Failed to retrieve attendance document: %v\n", err)
        return fmt.Errorf("failed to retrieve attendance document: %v", err)
    }

    var attendance Attendance
    if err := doc.DataTo(&attendance); err != nil {
        log.Printf("Failed to convert attendance document: %v\n", err)
        return fmt.Errorf("failed to convert attendance document: %v", err)
    }

    // Parse the StartTime to calculate the difference
    startTime, err := time.Parse(time.RFC3339, attendance.StartTime)
    if err != nil {
        log.Printf("Invalid StartTime format: %v\n", err)
        return fmt.Errorf("invalid StartTime format: %v", err)
    }

    // Calculate the total time spent as the difference between EndTime and StartTime
    totalTime := endTime.Sub(startTime)
    log.Printf("TotalTime calculated: %s\n", totalTime)

    // Step 1: Retrieve Ad data (Salary and ParentId) from the referenced Ad document
    adDoc, err := attendance.Ad.Get(ctx)
    if err != nil {
        log.Printf("Failed to retrieve ad document: %v\n", err)
        return fmt.Errorf("failed to retrieve ad document: %v", err)
    }

    var adData struct {
        Salary   int                    `firestore:"Salary"`
        ParentId *firestore.DocumentRef `firestore:"ParentId"`
    }
    if err := adDoc.DataTo(&adData); err != nil {
        log.Printf("Failed to convert ad document: %v\n", err)
        return fmt.Errorf("failed to convert ad document: %v", err)
    }
    salary := adData.Salary
    log.Printf("Salary retrieved from Ad document: %d\n", salary)

    // Step 2: Retrieve WorkerID from the referenced Request document in "requests" collection
    requestDoc, err := attendance.Request.Get(ctx)
    if err != nil {
        log.Printf("Failed to retrieve request document: %v\n", err)
        return fmt.Errorf("failed to retrieve request document: %v", err)
    }

    var requestData struct {
        WorkerID *firestore.DocumentRef `firestore:"WorkerID"`
    }
    if err := requestDoc.DataTo(&requestData); err != nil {
        log.Printf("Failed to convert request document: %v\n", err)
        return fmt.Errorf("failed to convert request document: %v", err)
    }

    // Step 3: Retrieve account reference for Parent and Worker from their user documents
    // Fetch Parent's account reference
    parentUserDoc, err := adData.ParentId.Get(ctx)
    if err != nil {
        log.Printf("Failed to retrieve parent user document: %v\n", err)
        return fmt.Errorf("failed to retrieve parent user document: %v", err)
    }
    var parentUserData struct {
        Account *firestore.DocumentRef `firestore:"account"`
    }
    if err := parentUserDoc.DataTo(&parentUserData); err != nil {
        log.Printf("Failed to convert parent user data: %v\n", err)
        return fmt.Errorf("failed to convert parent user data: %v", err)
    }
    parentAccountRef := parentUserData.Account

    // Fetch Worker's account reference
    workerUserDoc, err := requestData.WorkerID.Get(ctx)
    if err != nil {
        log.Printf("Failed to retrieve worker user document: %v\n", err)
        return fmt.Errorf("failed to retrieve worker user document: %v", err)
    }
    var workerUserData struct {
        Account *firestore.DocumentRef `firestore:"account"`
    }
    if err := workerUserDoc.DataTo(&workerUserData); err != nil {
        log.Printf("Failed to convert worker user data: %v\n", err)
        return fmt.Errorf("failed to convert worker user data: %v", err)
    }
    workerAccountRef := workerUserData.Account

    // Step 4: Retrieve and update Parent and Worker account balances
    // Fetch Parent's Account
    parentAccountDoc, err := parentAccountRef.Get(ctx)
    if err != nil {
        log.Printf("Failed to retrieve parent account: %v\n", err)
        return fmt.Errorf("failed to retrieve parent account: %v", err)
    }
    var parentAccount Account
    if err := parentAccountDoc.DataTo(&parentAccount); err != nil {
        log.Printf("Failed to convert parent account data: %v\n", err)
        return fmt.Errorf("failed to convert parent account data: %v", err)
    }
    log.Printf("Parent account balance retrieved: %d\n", parentAccount.Balance)

    // Fetch Worker's Account
    workerAccountDoc, err := workerAccountRef.Get(ctx)
    if err != nil {
        log.Printf("Failed to retrieve worker account: %v\n", err)
        return fmt.Errorf("failed to retrieve worker account: %v", err)
    }
    var workerAccount Account
    if err := workerAccountDoc.DataTo(&workerAccount); err != nil {
        log.Printf("Failed to convert worker account data: %v\n", err)
        return fmt.Errorf("failed to convert worker account data: %v", err)
    }

    // Step 5: Perform balance transfer
    if parentAccount.Balance < salary {
        log.Println("Insufficient balance in parent's account")
        return fmt.Errorf("insufficient balance in parent's account")
    }
    parentAccount.Balance -= salary
    workerAccount.Balance += salary
    log.Printf("Balance transferred. Parent balance after transfer: %d, Worker balance: %d\n", parentAccount.Balance, workerAccount.Balance)

    // Step 6: Update the accounts with the new balances
    _, err = parentAccountRef.Set(ctx, map[string]interface{}{
        "balance": parentAccount.Balance,
    }, firestore.MergeAll)
    if err != nil {
        log.Printf("Failed to update parent account balance: %v\n", err)
        return fmt.Errorf("failed to update parent account balance: %v", err)
    }

    _, err = workerAccountRef.Set(ctx, map[string]interface{}{
        "balance": workerAccount.Balance,
    }, firestore.MergeAll)
    if err != nil {
        log.Printf("Failed to update worker account balance: %v\n", err)
        return fmt.Errorf("failed to update worker account balance: %v", err)
    }

    // Step 7: Update EndTime and TotalTime in the Attendance document
    _, err = attendanceRef.Update(ctx, []firestore.Update{
        {Path: "EndTime", Value: endTime.Format(time.RFC3339)},       // Set EndTime
        {Path: "TotalTime", Value: totalTime.String()},               // Store the duration as a string (e.g., "1h30m")
    })
    if err != nil {
        log.Printf("Failed to update attendance document: %v\n", err)
        return fmt.Errorf("failed to update attendance document: %v", err)
    }

    log.Println("UpdateEndTime function completed successfully")
    return nil
}