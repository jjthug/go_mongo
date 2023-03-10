package main

import (
	"github.com/gin-gonic/gin"
	"go_mongo/routes"
	"os"
)

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		port = "8000"
	}
	router := gin.New()
	router.Use(gin.Logger())

	routes.AuthRouter(router)
	routes.UserRoutes(router)

	router.GET("/user", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "success"})
	})

	router.Run(":" + port)
}
