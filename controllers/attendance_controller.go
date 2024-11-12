package controllers

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/imthewolverine/schoolPolice-backend/services"
)

// UpdateStartTime updates the StartTime field of an Attendance document with the current time.
// @Summary Update StartTime of an Attendance document
// @Description Set the StartTime field of an Attendance document to the current time
// @Tags Attendance
// @Param attendanceID path string true "Attendance ID"
// @Success 200 {object} map[string]string{"message": "StartTime updated successfully"}
// @Failure 400 {object} map[string]string{"error": "Invalid attendance ID"}
// @Failure 500 {object} map[string]string{"error": "Failed to update StartTime"}
// @Router /attendance/{attendanceID}/start [patch]
func UpdateStartTime(c *gin.Context, attendanceService *services.AttendanceService) {
    attendanceID := c.Param("attendanceID") // Retrieve the attendance ID from the URL

    ctx := context.Background()

    // Call the service function to update StartTime
    err := attendanceService.UpdateStartTime(ctx, attendanceID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update StartTime"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "StartTime updated successfully"})
}

// UpdateEndTime updates the EndTime field to the current time and calculates the total time spent.
// @Summary Update EndTime of an Attendance document
// @Description Set the EndTime field of an Attendance document to the current time and calculate TotalTime
// @Tags Attendance
// @Param attendanceID path string true "Attendance ID"
// @Success 200 {object} map[string]string{"message": "EndTime and TotalTime updated successfully"}
// @Failure 400 {object} map[string]string{"error": "Invalid attendance ID"}
// @Failure 500 {object} map[string]string{"error": "Failed to update EndTime and TotalTime"}
// @Router /attendance/{attendanceID}/end [patch]
func UpdateEndTime(c *gin.Context, attendanceService *services.AttendanceService) {
    attendanceID := c.Param("attendanceID") // Retrieve the attendance ID from the URL

    ctx := context.Background()

    // Call the service function to update EndTime and calculate TotalTime
    err := attendanceService.UpdateEndTime(ctx, attendanceID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update EndTime and TotalTime"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "EndTime and TotalTime updated successfully"})
}