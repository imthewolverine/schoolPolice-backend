package controllers

import (
	"context"
	"net/http"

	"cloud.google.com/go/firestore"
	"github.com/gin-gonic/gin"
)

// Login handles user login by retrieving data from Firestore
func Login(c *gin.Context, firestoreClient *firestore.Client) {
    // Access Firestore client here
    ctx := context.Background()

    // Example Firestore query to get all users
    users := firestoreClient.Collection("users")
    docs, err := users.Documents(ctx).GetAll()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve users"})
        return
    }

    // Process documents and store user data
    var userData []map[string]interface{}
    for _, doc := range docs {
        userData = append(userData, doc.Data())
    }

    // Return the user data as a JSON response
    c.JSON(http.StatusOK, gin.H{"users": userData})
}
