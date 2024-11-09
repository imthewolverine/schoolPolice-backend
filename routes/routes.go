package routes

import (
	"cloud.google.com/go/firestore"
	"github.com/gin-gonic/gin"
	"github.com/imthewolverine/schoolPolice-backend/controllers"
)

func RegisterRoutes(router *gin.Engine, firestoreClient *firestore.Client) {
    // Define the /login route and pass the Firestore client to the Login handler
    router.POST("/login", func(c *gin.Context) {
        controllers.Login(c, firestoreClient)
    })
    // Add other routes here and pass firestoreClient as needed
}
