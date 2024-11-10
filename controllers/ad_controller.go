package controllers

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/imthewolverine/schoolPolice-backend/services"
)

func GetAllAds(c *gin.Context, adService *services.AdService) {
    // Fetch all ads from Firestore
    ctx := context.Background()
    ads, err := adService.FetchAllAds(ctx)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve ads"})
        return
    }

    // Return ads as JSON response
    c.JSON(http.StatusOK, gin.H{"ads": ads})
}
