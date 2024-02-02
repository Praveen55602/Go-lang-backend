package routes

import (
	controller "auth/controllers"
	"auth/middleware"

	"github.com/gin-gonic/gin"
)

func UserRoutes(incomingRoutes *gin.Engine) {
	//using this middleware we are protecting the routes
	incomingRoutes.Use(middleware.Authenticate)
	incomingRoutes.GET("/users", controller.GetUsers)

	//getUser function will be called which in turn will return a handler function, that function will take in the context send by this get request and will perform some operation on it since that context is of type pointer we can directly manupulate it
	incomingRoutes.GET("/users/:user_id", controller.GetUser())
}
