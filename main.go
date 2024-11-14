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
)

func main() {
	log.Println("Starting application...")

	// Load environment variables
	config.LoadEnv()
	log.Println("Environment variables loaded")

	// Initialize context
	ctx := context.Background()

	// Setup Firestore client
	firestoreClient, err := setupFirestore(ctx)
	if err != nil {
		log.Fatalf("Failed to set up Firestore: %v", err)
	}
	defer firestoreClient.Close()
	log.Println("Firestore client successfully initialized")

	// Initialize Gin router
	r := gin.Default()

	// Pass Firestore client to routes
	routes.RegisterRoutes(r, firestoreClient)
	log.Println("Routes registered successfully")

	// Use PORT environment variable from Cloud Run
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default to 8080 if PORT is not set
	}
	log.Printf("Starting server on port %s", port)

	// Listen and serve on the specified port
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func setupFirestore(ctx context.Context) (*firestore.Client, error) {
	log.Println("Initializing Firestore...")

	// Initialize Firebase app with default credentials on Cloud Run
	app, err := firebase.NewApp(ctx, nil) // No credentials file needed for default account
	if err != nil {
		return nil, fmt.Errorf("failed to initialize Firebase app: %v", err)
	}

	client, err := app.Firestore(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create Firestore client: %v", err)
	}

	log.Println("Firestore client created successfully")
	return client, nil
}
