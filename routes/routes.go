package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/imthewolverine/schoolPolice_backend/controllers"
)

func RegisterRoutes(r *gin.Engine) {
    // User routes
    r.POST("/users", controllers.CreateUser)
    // Additional routes can be added here
}
