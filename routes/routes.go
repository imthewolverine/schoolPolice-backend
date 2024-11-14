package routes

import (
	"cloud.google.com/go/firestore"
	"github.com/gin-gonic/gin"
	"github.com/imthewolverine/schoolPolice-backend/controllers" // Import the docs package to initialize the generated Swagger docs
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

	// Test notification route
	router.POST("/test-notification", controllers.TestNotificationHandler)

	// FCM token route
	router.POST("/save-fcm-token", func(c *gin.Context) {
		services.SaveFCMToken(c, firestoreClient) // Pass Firestore client to the handler
	})
	// Ad routes
	router.GET("/ads", func(c *gin.Context) {
		controllers.GetAllAds(c, adService)
	})
	router.GET("/ads/:id", func(c *gin.Context) { // Route to get a single ad by ID
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
	router.PATCH("/attendance/:attendanceID/start", func(c *gin.Context) {
		controllers.UpdateStartTime(c, attendanceService)
	})
	router.PATCH("/attendance/:attendanceID/end", func(c *gin.Context) {
		controllers.UpdateEndTime(c, attendanceService)
	})
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swagFiles.Handler))
}
