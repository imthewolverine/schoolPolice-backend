package services

import (
	"context"
	"log"
	"net/http"

	"cloud.google.com/go/firestore"
	"github.com/gin-gonic/gin"
)

// TokenRequest represents the structure of the token data
type TokenRequest struct {
	Token  string `json:"token"`
	UserID string `json:"userId"`
}

// SaveFCMToken handles requests to save FCM tokens
func SaveFCMToken(c *gin.Context, firestoreClient *firestore.Client) {
	var tokenReq TokenRequest

	// Bind JSON request to TokenRequest struct
	if err := c.BindJSON(&tokenReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	// Save the token and user ID to Firestore
	_, _, err := firestoreClient.Collection("userTokens").Add(context.Background(), map[string]interface{}{
		"userId": tokenReq.UserID,
		"token":  tokenReq.Token,
	})
	if err != nil {
		log.Printf("Failed to save FCM token: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save token"})
		return
	}

	// Respond with success
	c.JSON(http.StatusOK, gin.H{"message": "Token saved successfully"})
}
