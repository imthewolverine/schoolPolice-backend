package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

// Define the FCMMessage structure to match the JSON payload
type FCMMessage struct {
	Message Message `json:"message"`
}

type Message struct {
	Token        string       `json:"token"`
	Notification Notification `json:"notification"`
}

type Notification struct {
	Title string `json:"title"`
	Body  string `json:"body"`
}

// SendNotification sends an FCM notification to a specific device
func SendNotification(ctx context.Context, fcmToken, title, body string) error {
	// Construct the message payload
	message := FCMMessage{
		Message: Message{
			Token: fcmToken,
			Notification: Notification{
				Title: title,
				Body:  body,
			},
		},
	}

	// Convert the message to JSON
	payload, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal FCM message: %v", err)
	}

	// Firebase Cloud Messaging endpoint for V1 API
	url := "https://fcm.googleapis.com/v1/projects/school-police-c59de/messages:send"
	//url := "https://fcm.googleapis.com/projects/school-police-c59de/messages:send"

	// Create the HTTP request
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(payload))
	if err != nil {
		return fmt.Errorf("failed to create HTTP request: %v", err)
	}

	// Get a dynamic auth token
	authToken, err := GetFirebaseAuthToken()
	if err != nil {
		return fmt.Errorf("failed to get auth token: %v", err)
	}
	log.Println("Generated Auth Token:", authToken)

	// Set headers
	req.Header.Set("Authorization", "Bearer "+authToken)
	req.Header.Set("Content-Type", "application/json")

	// Create an HTTP client with a timeout
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send FCM request: %v", err)
	}
	defer resp.Body.Close()

	// Handle non-200 response
	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		log.Printf("Failed to send notification, status: %s, response: %s", resp.Status, string(bodyBytes))
		return fmt.Errorf("failed to send notification, status: %s", resp.Status)
	}

	log.Println("Notification sent successfully!")
	return nil
}
