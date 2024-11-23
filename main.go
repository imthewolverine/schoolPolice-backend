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
	log.Println("Environment variables loaded")

	// Initialize context
	ctx := context.Background()

	// Initialize Firestore
	firestoreClient, err := setupFirestore(context.Background())
	if err != nil {
		log.Fatalf("Failed to set up Firestore: %v", err)
	}
	defer firestoreClient.Close()
	log.Println("Firestore client successfully initialized")

	// Set up the Gin router
	r := gin.Default()

	// Register all routes, passing the Firestore client to the router
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
	credentialsPath := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")
	if credentialsPath == "" {
		return nil, fmt.Errorf("GOOGLE_APPLICATION_CREDENTIALS is not set")
	}

	sa := option.WithCredentialsFile(credentialsPath)
	app, err := firebase.NewApp(ctx, nil, sa)
	if err != nil {
		return nil, err
	}

	client, err := app.Firestore(ctx)
	if err != nil {
		return nil, err
	}
	return client, nil
}
