package main

import (
	"context"
	"log"
	"os"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"github.com/gin-gonic/gin"
	"github.com/imthewolverine/schoolPolice-backend/config"
	"github.com/imthewolverine/schoolPolice-backend/routes"
	"google.golang.org/api/option"
)

func main() {
    // Load environment variables
    config.LoadEnv()

    // Initialize context
    ctx := context.Background()

    // Setup Firestore client
    firestoreClient, err := setupFirestore(ctx)
    if err != nil {
        log.Fatalf("failed to set up Firestore: %v", err)
    }
    defer firestoreClient.Close() // Close client when done

    // Initialize Gin router
    r := gin.Default()

    // Pass Firestore client to routes
    routes.RegisterRoutes(r, firestoreClient)

    // Start server
    r.Run(":8080") // Run on port 8080
}

func setupFirestore(ctx context.Context) (*firestore.Client, error) {
    projectID := "your-firebase-project-id" // Replace with your Firebase project ID

    // Use service account key specified by GOOGLE_APPLICATION_CREDENTIALS
    sa := option.WithCredentialsFile(os.Getenv("GOOGLE_APPLICATION_CREDENTIALS"))
    app, err := firebase.NewApp(ctx, &firebase.Config{ProjectID: projectID}, sa)
    if err != nil {
        return nil, err
    }

    client, err := app.Firestore(ctx)
    if err != nil {
        return nil, err
    }

    return client, nil
}
