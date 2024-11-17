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

/*package main

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

	// Initialize Firestore
	firestoreClient, err := setupFirestore(context.Background())
	if err != nil {
		log.Fatalf("Failed to initialize Firestore: %v", err)
	}

	// Set up the Gin router
	r := gin.Default()

	// Register all routes, passing the Firestore client to the router
	routes.RegisterRoutes(r, firestoreClient)

	// Start the server
	r.Run("0.0.0.0:8080")
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
}*/
