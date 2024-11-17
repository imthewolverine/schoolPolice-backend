package routes

import (
	"net/http"

	"cloud.google.com/go/firestore"
	"github.com/gin-gonic/gin"
	"github.com/imthewolverine/schoolPolice-backend/controllers"
	"github.com/imthewolverine/schoolPolice-backend/services"
	swagFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func RegisterRoutes(router *gin.Engine, firestoreClient *firestore.Client) {
	authService := services.NewAuthService(firestoreClient)
	userService := services.NewUserService(firestoreClient)
	adService := services.NewAdService(firestoreClient)
	attendanceService := &services.AttendanceService{FirestoreClient: firestoreClient}

	// User routes
	router.POST("/register", func(c *gin.Context) {
		controllers.RegisterUser(c, userService)
	})
	router.POST("/login", func(c *gin.Context) {
		controllers.Login(c, authService)
	})

	// Send notification route
	router.POST("/send-notification", func(c *gin.Context) {
		var req struct {
			UserID string `json:"userId"` // Accept userId instead of token
			Title  string `json:"title"`
			Body   string `json:"body"`
		}

		// Bind JSON request data to req
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
			return
		}

		// Fetch the FCM token from Firestore using userId
		token, err := services.FetchTokenFromDatabase(req.UserID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch token from database"})
			return
		}

		// Send the push notification using the retrieved token
		err = services.SendPushNotification(token, req.Title, req.Body)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send notification"})
			return
		}

		// Return success response
		c.JSON(http.StatusOK, gin.H{"message": "Notification sent successfully"})
	})

	// FCM token route
	router.POST("/save-fcm-token", func(c *gin.Context) {
		services.SaveFCMToken(c, firestoreClient)
	})

	// Ad routes
	router.GET("/ads", func(c *gin.Context) {
		controllers.GetAllAds(c, adService)
	})
	router.GET("/ads/:id", func(c *gin.Context) {
		controllers.GetAdByID(c, adService)
	})
	router.POST("/ads", func(c *gin.Context) {
		controllers.AddAd(c, adService)
	})
	router.DELETE("/ads/:id", func(c *gin.Context) {
		controllers.DeleteAdByID(c, adService)
	})
	router.PUT("/ads/:id", func(c *gin.Context) {
		controllers.UpdateAdByID(c, adService)
	})
	router.POST("/ads/:adID/requests", func(c *gin.Context) {
		controllers.AddRequestToAd(c, adService)
	})
	router.GET("/users/:userID/ads/requests", func(c *gin.Context) {
		controllers.GetUserAdsRequests(c, adService)
	})
	router.POST("/requests/accept", func(c *gin.Context) {
		controllers.AcceptRequest(c, adService)
	})
	router.GET("/requests/accepted/:workerID", func(c *gin.Context) {
		controllers.GetAcceptedRequestsByWorkerID(c, adService)
	})

	// Attendance routes
	router.PATCH("/attendance/:attendanceID/start", func(c *gin.Context) {
		controllers.UpdateStartTime(c, attendanceService)
	})
	router.PATCH("/attendance/:attendanceID/end", func(c *gin.Context) {
		controllers.UpdateEndTime(c, attendanceService)
	})

	// Swagger route
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swagFiles.Handler))
}
