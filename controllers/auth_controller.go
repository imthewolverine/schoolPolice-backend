package controllers

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/imthewolverine/schoolPolice-backend/services"
)

// LoginRequest represents the expected JSON input for login
type LoginRequest struct {
    UsernameOrEmail string `json:"usernameOrEmail"`
    Password        string `json:"password"`
}

// LoginResponse represents the JSON output for a successful login
type LoginResponse struct {
    Token string `json:"token"`
}

// Login logs in a user and returns a token.
// @Summary Login user
// @Description Log in a user with username and password
// @Tags Auth
// @Accept json
// @Produce json
// @Param login body controllers.LoginRequest true "Login Request"
// @Success 200 {object} controllers.LoginResponse
// @Failure 400 {object} controllers.ErrorResponse
// @Failure 401 {object} controllers.ErrorResponse
// @Failure 500 {object} controllers.ErrorResponse
// @Router /login [post]
func Login(c *gin.Context, authService *services.AuthService) {
    var req LoginRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid request format"})
        return
    }

    // Authenticate user
    ctx := context.Background()
    token, err := authService.AuthenticateUser(ctx, req.UsernameOrEmail, req.Password)
    if err != nil {
        c.JSON(http.StatusUnauthorized, ErrorResponse{Error: err.Error()})
        return
    }

    // Return JWT token
    c.JSON(http.StatusOK, LoginResponse{Token: token})
}
