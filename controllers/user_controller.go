package controllers

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/imthewolverine/schoolPolice-backend/models"
	"github.com/imthewolverine/schoolPolice-backend/services"
)

type RegisterRequest struct {
    Name         string `json:"name" binding:"required"`
    Email        string `json:"email" binding:"required"`
    Password     string `json:"password" binding:"required"`
    Address      string `json:"address"`
    PhoneNumber  string `json:"phoneNumber"`
    Rating       float64 `json:"rating"`
    TotalWorkCount int   `json:"totalWorkCount"`
    UserID       int     `json:"userid"`
}

func RegisterUser(c *gin.Context, userService *services.UserService) {
    var req RegisterRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
        return
    }

    user := models.User{
        Name:          req.Name,
        Email:         req.Email,
        Password:      req.Password, // In production, hash the password
        Address:       req.Address,
        PhoneNumber:   req.PhoneNumber,
        Rating:        req.Rating,
        TotalWorkCount: req.TotalWorkCount,
        UserID:        req.UserID,
    }

    // Create user in Firestore
    ctx := context.Background()
    if err := userService.CreateUser(ctx, user); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "User created successfully"})
}
