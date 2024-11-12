package controllers

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/imthewolverine/schoolPolice-backend/models"
	"github.com/imthewolverine/schoolPolice-backend/services"
)

// SuccessResponse represents a generic success message
type SuccessResponse struct {
    Message string `json:"message"`
}

// RegisterRequest represents the expected JSON input for user registration
type RegisterRequest struct {
    Name           string  `json:"name" binding:"required"`
    Email          string  `json:"email" binding:"required"`
    Password       string  `json:"password" binding:"required"`
    Address        string  `json:"address"`
    PhoneNumber    string  `json:"phoneNumber"`
    Rating         float64 `json:"rating"`
    TotalWorkCount int     `json:"totalWorkCount"`
    UserID         int     `json:"userid"`
}


// RegisterUser creates a new user.
// @Summary Register a new user
// @Description Register a new user with required details
// @Tags User
// @Accept json
// @Produce json
// @Param user body controllers.RegisterRequest true "Register Request"
// @Success 200 {object} controllers.SuccessResponse
// @Failure 400 {object} controllers.ErrorResponse
// @Failure 500 {object} controllers.ErrorResponse
// @Router /register [post]
func RegisterUser(c *gin.Context, userService *services.UserService) {
    var req RegisterRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid request format"})
        return
    }

    user := models.User{
        Name:           req.Name,
        Email:          req.Email,
        Password:       req.Password, // In production, hash the password
        Address:        req.Address,
        PhoneNumber:    req.PhoneNumber,
        Rating:         req.Rating,
        TotalWorkCount: req.TotalWorkCount,
        UserID:         req.UserID,
    }

    // Create user in Firestore
    ctx := context.Background()
    if err := userService.CreateUser(ctx, user); err != nil {
        c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "User created successfully"})
}
