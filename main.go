package main

import (
	"context"
	"fmt"
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
    defer firestoreClient.Close()

    // Initialize Gin router
    r := gin.Default()

    // Pass Firestore client to routes
    routes.RegisterRoutes(r, firestoreClient)

    // Start server
    r.Run(":8080") // Run on port 8080
}

func setupFirestore(ctx context.Context) (*firestore.Client, error) {
    credentialsPath := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")
    
    if credentialsPath == "" {
        return nil, fmt.Errorf("GOOGLE_APPLICATION_CREDENTIALS is not set")
    }
    
    // Initialize Firebase app with credentials
    sa := option.WithCredentialsFile(credentialsPath)
    app, err := firebase.NewApp(ctx, nil, sa) // Omit ProjectID if credentials contain it
    if err != nil {
        return nil, err
    }

    client, err := app.Firestore(ctx)
    if err != nil {
        return nil, err
    }

    return client, nil
}

