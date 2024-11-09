package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/imthewolverine/schoolPolice_backend/models"
	"github.com/imthewolverine/schoolPolice_backend/services"
)

// CreateUser - Example handler for creating a new user
func CreateUser(c *gin.Context) {
    var user models.User
    if err := c.ShouldBindJSON(&user); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    createdUser, err := services.CreateUser(user)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
        return
    }
    c.JSON(http.StatusOK, createdUser)
}
