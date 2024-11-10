package routes

import (
	"cloud.google.com/go/firestore"
	"github.com/gin-gonic/gin"
	"github.com/imthewolverine/schoolPolice-backend/controllers"
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
}
