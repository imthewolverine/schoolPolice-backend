package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/imthewolverine/schoolPolice-backend/services"
)
type AdRequest struct {
    Date          string `json:"date"`
    Description   string `json:"description"`
    ParentId      string `json:"parentId,omitempty"` // JSON input as string
    Salary        int    `json:"salary"`
    School        string `json:"school"`
    SchoolAddress string `json:"schoolAddress"`
    Status        string `json:"status"`
    Time          string `json:"time"` // Use string for time in request to simplify parsing
}

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

func GetAdByID(c *gin.Context, adService *services.AdService) {
    id := c.Param("id") // Retrieve the 'id' parameter from the URL

    // Fetch the ad by ID
    ctx := context.Background()
    ad, err := adService.FetchAdByID(ctx, id)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
        return
    }

    // Return the ad as JSON response
    c.JSON(http.StatusOK, ad)
}

func AddAd(c *gin.Context, adService *services.AdService) {
    var adRequest AdRequest
    if err := c.ShouldBindJSON(&adRequest); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format", "details": err.Error()})
        return
    }

    // Convert AdRequest to Ad
    ctx := context.Background()
    var ad services.Ad
    ad.Date = adRequest.Date
    ad.Description = adRequest.Description
    ad.Salary = adRequest.Salary
    ad.School = adRequest.School
    ad.SchoolAddress = adRequest.SchoolAddress
    ad.Status = adRequest.Status

    // Parse the Time field from string to time.Time
    parsedTime, err := time.Parse(time.RFC3339, adRequest.Time)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid time format", "details": err.Error()})
        return
    }
    ad.Time = parsedTime

    // Convert ParentId to DocumentRef if provided
    if adRequest.ParentId != "" {
        ad.ParentId = adService.FirestoreClient.Doc(adRequest.ParentId)
    }

    // Call the service to add the ad to Firestore
    id, err := adService.AddAd(ctx, ad)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add ad"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Ad added successfully", "id": id})
}




func DeleteAdByID(c *gin.Context, adService *services.AdService) {
    id := c.Param("id") // Get the ad ID from the URL

    ctx := context.Background()
    if err := adService.DeleteAdByID(ctx, id); err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Ad not found"})
        return
    }
	c.JSON(http.StatusOK, gin.H{"message": "Ad deleted successfully"})
}

func UpdateAdByID(c *gin.Context, adService *services.AdService) {
    id := c.Param("id") // Get the ad ID from the URL

    var adRequest AdRequest
    if err := c.ShouldBindJSON(&adRequest); err != nil {
        // Provide detailed error response if JSON parsing fails
        c.JSON(http.StatusBadRequest, gin.H{
            "error":   "Invalid request format",
            "details": err.Error(),
        })
        return
    }

    // Convert AdRequest to Ad struct for Firestore operations
    ctx := context.Background()
    var ad services.Ad
    ad.Date = adRequest.Date
    ad.Description = adRequest.Description
    ad.Salary = adRequest.Salary
    ad.School = adRequest.School
    ad.SchoolAddress = adRequest.SchoolAddress
    ad.Status = adRequest.Status

    // Parse the Time field from string to time.Time
    parsedTime, err := time.Parse(time.RFC3339, adRequest.Time)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error":   "Invalid time format",
            "details": err.Error(),
        })
        return
    }
    ad.Time = parsedTime

    // Convert ParentId to DocumentRef if provided
    if adRequest.ParentId != "" {
        ad.ParentId = adService.FirestoreClient.Doc(adRequest.ParentId)
    }

    // Call the service to update the ad in Firestore
    if err := adService.UpdateAdByID(ctx, id, ad); err != nil {
        c.JSON(http.StatusNotFound, gin.H{
            "error":   "Failed to update ad",
            "details": err.Error(),
        })
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Ad updated successfully"})
}

func AddRequestToAd(c *gin.Context, adService *services.AdService) {
    adID := c.Param("adID") // Get the ad ID from the URL

    var requestData struct {
        WorkerID string `json:"workerID"`
    }
    if err := c.ShouldBindJSON(&requestData); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error":   "Invalid request format",
            "details": err.Error(),
        })
        return
    }

    ctx := context.Background()

    // Step 1: Create a new request document
    requestRef, err := adService.CreateRequest(ctx, requestData.WorkerID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create request"})
        return
    }

    // Step 2: Add the request reference to the ad's requests array
    if err := adService.AddRequestToAd(ctx, adID, requestRef); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update ad with request"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Request added to ad successfully"})
}

func GetUserAdsRequests(c *gin.Context, adService *services.AdService) {
    userID := c.Param("userID") // Get the user ID from the URL

    ctx := context.Background()

    // Get all requests associated with the user's ads
    requests, err := adService.GetUserAdsRequests(ctx, userID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error":   "Failed to retrieve requests",
            "details": err.Error(),
        })
        return
    }

    c.JSON(http.StatusOK, gin.H{"requests": requests})
}
