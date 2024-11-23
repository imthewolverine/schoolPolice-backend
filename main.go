package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/imthewolverine/schoolPolice-backend/config"
	"github.com/imthewolverine/schoolPolice-backend/routes"
	"github.com/imthewolverine/schoolPolice-backend/services"
)

func main() {
	// Load environment variables
	config.LoadEnv()

	// Check that GOOGLE_APPLICATION_CREDENTIALS is set
	credentialsPath := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")
	if credentialsPath == "" {
		log.Fatalf("GOOGLE_APPLICATION_CREDENTIALS is not set")
	}

	// Initialize Firebase (both Firestore and FCM)
	err := services.InitializeFirebase()
	if err != nil {
		log.Fatalf("Failed to initialize Firebase: %v", err)
	}
	defer services.CloseFirestoreClient() // Close Firestore client on exit

	// Set up the Gin router
	r := gin.Default()

	// Register routes and pass Firestore client from services
	routes.RegisterRoutes(r, services.GetFirestoreClient())

	// Start the server
	r.Run("0.0.0.0:8080") // replace with localhost:8080 if testing on the same device
}
