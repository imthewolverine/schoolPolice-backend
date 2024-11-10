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

    // Define the /register route
    router.POST("/register", func(c *gin.Context) {
        controllers.RegisterUser(c, userService)
    })

    // Existing login route
    router.POST("/login", func(c *gin.Context) {
        controllers.Login(c, authService)
    })
}
