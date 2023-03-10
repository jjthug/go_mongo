package routes

import (
	"github.com/gin-gonic/gin"
	controller "go_mongo/controllers"
)

func AuthRouter(incomingRoutes *gin.Engine) {
	incomingRoutes.POST("/users/signup", controller.Signup())
	incomingRoutes.POST("/users/login", controller.Login())

}
