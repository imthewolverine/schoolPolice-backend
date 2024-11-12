package controllers

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/imthewolverine/schoolPolice-backend/services"
)

// TestNotificationHandler handles the test notification request
func TestNotificationHandler(c *gin.Context) {
	var req struct {
		Token string `json:"token"`
		Title string `json:"title"`
		Body  string `json:"body"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	// Send the notification using the services package
	ctx := context.Background()
	if err := services.SendNotification(ctx, req.Token, req.Title, req.Body); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Notification sent successfully!"})
}
