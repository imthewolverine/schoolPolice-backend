package services

import (
	"context"
	"log"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
)

// Attendance represents the structure of the Attendance document
type Attendance struct {
    ID               string                 `json:"id,omitempty"`
    EndTime          string                 `firestore:"EndTime"`
    LocationVerified bool                   `firestore:"LocationVerified"`
    StartTime        string                 `firestore:"StartTime"`
    TotalTime        string                 `firestore:"TotalTime"`
    Ad               *firestore.DocumentRef `firestore:"ad"`
    Request          *firestore.DocumentRef `firestore:"request"`
}


// UpdateRequestStatus updates the status of a request in Firestore
func (s *AdService) UpdateRequestStatus(ctx context.Context, requestID, status string) error {
    // Log the requestID and status for debugging
    log.Printf("Attempting to update request with ID: %s to status: %s\n", requestID, status)

    requestRef := s.FirestoreClient.Collection("requests").Doc(requestID)
    _, err := requestRef.Update(ctx, []firestore.Update{
        {Path: "Status", Value: status},
    })

    if err != nil {
        // Log the error message for further debugging
        log.Printf("Error updating request status: %v\n", err)
    }

    return err
}

// CreateAttendance adds a new document to the Attendance collection
func (s *AdService) CreateAttendance(ctx context.Context, attendance Attendance) error {
    _, _, err := s.FirestoreClient.Collection("Attendance").Add(ctx, attendance)
    return err
}

// Request represents the structure of a request document
type Request struct {
    WorkerID string                 `firestore:"WorkerID"`
    Status   string                 `firestore:"Status"`
    // Add other relevant fields from your request document
}

// FetchAcceptedRequestsByWorkerID fetches all accepted requests and their associated attendance records for a specific workerID
func (s *AdService) FetchAcceptedRequestsByWorkerID(ctx context.Context, workerID string) ([]Attendance, error) {
    var attendances []Attendance

    // Create a DocumentRef for the workerID path
    workerRef := s.FirestoreClient.Doc("user/" + workerID)
    log.Printf("Fetching accepted requests for WorkerID: %v with Status: accepted\n", workerRef)

    // Step 1: Query "request" collection where WorkerID matches and Status is "accepted"
    requestIter := s.FirestoreClient.Collection("requests").
        Where("WorkerID", "==", workerRef).
        Where("Status", "==", "accepted").
        Documents(ctx)

    // Step 2: For each accepted request, fetch the associated Attendance document
    for {
        requestDoc, err := requestIter.Next()
        if err == iterator.Done {
            break
        }
        if err != nil {
            log.Printf("Error fetching request: %v\n", err)
            return nil, err
        }

        requestID := requestDoc.Ref.ID
        log.Printf("Accepted request found with ID: %s\n", requestID)

        // Create a full path reference to match in Attendance collection
        requestRef := s.FirestoreClient.Collection("requests").Doc(requestID)
        
        // Query Attendance collection for documents where `request` field matches the full path reference
        attendanceIter := s.FirestoreClient.Collection("Attendance").
            Where("request", "==", requestRef).
            Documents(ctx)

        for {
            attendanceDoc, err := attendanceIter.Next()
            if err == iterator.Done {
                break
            }
            if err != nil {
                log.Printf("Error fetching attendance for request %s: %v\n", requestID, err)
                return nil, err
            }

            // Map Firestore document to Attendance struct and include the document ID
            var attendance Attendance
            if err := attendanceDoc.DataTo(&attendance); err != nil {
                log.Printf("Error converting attendance document %s: %v\n", attendanceDoc.Ref.ID, err)
                return nil, err
            }
            attendance.ID = attendanceDoc.Ref.ID // Set the attendance document ID

            attendances = append(attendances, attendance)
            log.Printf("Attendance record added for request %s with attendance ID %s\n", requestID, attendance.ID)
        }
    }

    log.Printf("Total attendance records found: %d\n", len(attendances))
    return attendances, nil
}
