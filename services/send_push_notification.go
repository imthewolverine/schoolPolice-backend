package services

import (
	"context"
	"fmt"
	"log"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"firebase.google.com/go/messaging"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

var (
	firestoreClient *firestore.Client
	messagingClient *messaging.Client
)

// InitializeFirebase initializes Firebase app, Firestore, and FCM clients
func InitializeFirebase() error {
	opt := option.WithCredentialsFile("C:/Users/Dell/Desktop/schoolPolice-backend/config/secrets/school-police-c59de-firebase-adminsdk-45dsj-817b321849.json")
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		return fmt.Errorf("error initializing Firebase app: %v", err)
	}

	// Initialize Firestore client
	firestoreClient, err = app.Firestore(context.Background())
	if err != nil {
		return fmt.Errorf("error initializing Firestore client: %v", err)
	}

	// Initialize FCM messaging client
	messagingClient, err = app.Messaging(context.Background())
	if err != nil {
		return fmt.Errorf("error initializing FCM client: %v", err)
	}

	log.Println("Firebase initialized successfully.")
	return nil
}

// GetFirestoreClient provides access to the initialized Firestore client
func GetFirestoreClient() *firestore.Client {
	return firestoreClient
}

// SendPushNotification sends a push notification using Firebase Cloud Messaging
func SendPushNotification(token string, title string, body string) error {
	if messagingClient == nil {
		return fmt.Errorf("FCM client is not initialized")
	}

	// Log the details of the notification being sent
	log.Printf("Attempting to send push notification - Token: %s, Title: %s, Body: %s", token, title, body)

	message := &messaging.Message{
		Token: token,
		Notification: &messaging.Notification{
			Title: title,
			Body:  body,
		},
	}

	// Send the message
	response, err := messagingClient.Send(context.Background(), message)
	if err != nil {
		log.Printf("Failed to send push notification: %v", err)
		return err
	}

	log.Printf("Successfully sent message: %s", response)
	return nil
}

// FetchTokenFromDatabase retrieves the FCM token for a specific user from Firestore
func FetchTokenFromDatabase(userID string) (string, error) {
	if firestoreClient == nil {
		return "", fmt.Errorf("Firestore client is not initialized")
	}

	log.Printf("Fetching token for userID: %s", userID)

	iter := firestoreClient.Collection("userTokens").Where("userId", "==", userID).Documents(context.Background())
	defer iter.Stop()

	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			log.Printf("No token found for userID: %s", userID)
			return "", fmt.Errorf("no token found for userID: %s", userID)
		}
		if err != nil {
			log.Printf("Error fetching token for userID %s: %v", userID, err)
			return "", err
		}

		token, ok := doc.Data()["token"].(string)
		if !ok {
			log.Printf("Token not found in document for userID: %s", userID)
			return "", fmt.Errorf("token not found in document for userID: %s", userID)
		}

		log.Printf("Successfully fetched token: %s for userID: %s", token, userID)
		return token, nil
	}
}

// CloseFirestoreClient closes the Firestore client connection
func CloseFirestoreClient() {
	if firestoreClient != nil {
		err := firestoreClient.Close()
		if err != nil {
			log.Printf("Error closing Firestore client: %v", err)
		} else {
			log.Println("Firestore client closed successfully.")
		}
	}
}
