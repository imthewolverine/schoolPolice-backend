package services

import (
	"context"
	"fmt"
	"os"

	"golang.org/x/oauth2/google"
	"google.golang.org/api/idtoken"
	"google.golang.org/api/option"
)

func GetFirebaseAuthToken() (string, error) {
	// Path to your service account JSON file
	serviceAccountFile := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")
	fmt.Println("GOOGLE_APPLICATION_CREDENTIALS:", serviceAccountFile)
	fmt.Println("GOOGLE_CLOUD_PROJECT:", os.Getenv("GOOGLE_CLOUD_PROJECT"))

	// Read the service account key file
	jsonKey, err := os.ReadFile(serviceAccountFile)
	if err != nil {
		return "", fmt.Errorf("failed to read service account file: %v", err)
	}

	// Create credentials from the JSON key
	//creds, err := google.CredentialsFromJSON(context.Background(), jsonKey, "https://www.googleapis.com/auth/cloud-platform")
	creds, err := google.CredentialsFromJSON(context.Background(), jsonKey, "https://www.googleapis.com/auth/firebase.messaging")
	if err != nil {
		return "", fmt.Errorf("failed to create credentials from JSON: %v", err)
	}

	// Create token source
	tokenSource, err := idtoken.NewTokenSource(context.Background(), "https://fcm.googleapis.com", option.WithCredentials(creds))
	if err != nil {
		return "", fmt.Errorf("failed to create token source: %v", err)
	}

	// Retrieve OAuth token
	token, err := tokenSource.Token()
	if err != nil {
		return "", fmt.Errorf("failed to retrieve token: %v", err)
	}
	fmt.Println("Generated Auth Token:", token.AccessToken) // Print the generated token
	return token.AccessToken, nil
}
