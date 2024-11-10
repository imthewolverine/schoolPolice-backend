package controllers

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/imthewolverine/schoolPolice-backend/services"
)

type LoginRequest struct {
    UsernameOrEmail string `json:"usernameOrEmail"`
    Password        string `json:"password"`
}

type LoginResponse struct {
    Token string `json:"token"`
    Error string `json:"error,omitempty"`
}

// Login handles the login request, authenticates the user, and returns a JWT token
func Login(c *gin.Context, authService *services.AuthService) {
    var req LoginRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
        return
    }

    // Authenticate user
    ctx := context.Background()
    token, err := authService.AuthenticateUser(ctx, req.UsernameOrEmail, req.Password)
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
        return
    }

    // Return JWT token
    c.JSON(http.StatusOK, LoginResponse{Token: token})
}
