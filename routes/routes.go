package routes

import (
	"cloud.google.com/go/firestore"
	"github.com/gin-gonic/gin"
	"github.com/imthewolverine/schoolPolice-backend/controllers"

	// Import your new route
	"github.com/imthewolverine/schoolPolice-backend/services"
)

func RegisterRoutes(router *gin.Engine, firestoreClient *firestore.Client) {
	authService := services.NewAuthService(firestoreClient)
	userService := services.NewUserService(firestoreClient)
	adService := services.NewAdService(firestoreClient)

	// User routes
	router.POST("/register", func(c *gin.Context) {
		controllers.RegisterUser(c, userService)
	})
	router.POST("/login", func(c *gin.Context) {
		controllers.Login(c, authService)
	})

	// Ad routes
	router.GET("/ads", func(c *gin.Context) {
		controllers.GetAllAds(c, adService)
	})

	// Test notification route
	router.POST("/test-notification", controllers.TestNotificationHandler)

	// FCM token route
	router.POST("/save-fcm-token", func(c *gin.Context) {
		services.SaveFCMToken(c, firestoreClient) // Pass Firestore client to the handler
	})
}
