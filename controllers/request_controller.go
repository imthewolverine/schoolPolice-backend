package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/imthewolverine/schoolPolice-backend/services"
)

// AcceptRequestInput represents the JSON input for accepting a request
type AcceptRequestInput struct {
    RequestID string `json:"requestId" binding:"required"`
    AdID      string `json:"adId" binding:"required"`
}

// AcceptRequest handles accepting a request and creating an attendance record.
// @Summary Accept a request and create an attendance record
// @Description Accept a request, update its status to "accepted", and add a new document in the Attendance collection
// @Tags Requests
// @Accept json
// @Produce json
// @Param acceptRequest body AcceptRequestInput true "Accept Request Input"
// @Success 200 {object} map[string]string{"message": "Request accepted and attendance created successfully"}
// @Failure 400 {object} map[string]string{"error": "Invalid request format"}
// @Failure 500 {object} map[string]string{"error": "Failed to process request"}
// @Router /requests/accept [post]
func AcceptRequest(c *gin.Context, adService *services.AdService) {
    var input AcceptRequestInput
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format", "details": err.Error()})
        return
    }

    ctx := context.Background()

    // Step 1: Update the request status to "accepted"
    if err := adService.UpdateRequestStatus(ctx, input.RequestID, "accepted"); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update request status"})
        return
    }

    // Step 2: Create a new Attendance record
    attendance := services.Attendance{
        EndTime:          "",
        LocationVerified: true,
        StartTime:        time.Now().Format(time.RFC3339),
        TotalTime:        "",
        Ad:               adService.FirestoreClient.Doc("ad/" + input.AdID),
        Request:          adService.FirestoreClient.Doc("request/" + input.RequestID),
    }

    if err := adService.CreateAttendance(ctx, attendance); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create attendance record"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Request accepted and attendance created successfully"})
}



// GetAcceptedRequestsByWorkerID retrieves all accepted requests for a specific workerID.
// @Summary Get accepted requests by workerID
// @Description Retrieve all requests with status "accepted" for a given workerID
// @Tags Requests
// @Produce json
// @Param workerID path string true "Worker ID"
// @Success 200 {array} services.Request
// @Failure 500 {object} map[string]string{"error": "Failed to retrieve accepted requests"}
// @Router /requests/accepted/{workerID} [get]
func GetAcceptedRequestsByWorkerID(c *gin.Context, adService *services.AdService) {
    workerID := c.Param("workerID") // Get workerID from URL parameter

    ctx := context.Background()

    // Fetch all attendance records for accepted requests for the workerID
    attendances, err := adService.FetchAcceptedRequestsByWorkerID(ctx, workerID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve accepted requests"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"attendances": attendances})
}
