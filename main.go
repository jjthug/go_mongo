package main

import (
	"github.com/gin-gonic/gin"
	"go_mongo/routes"
	"log"
)

func main() {
	router := gin.Default()

	router.POST("/", routes.CreateUser)

	// called as localhost:3000/getOne/{id}
	router.GET("getOne/:postId", routes.ReadOneUser)
	router.GET("/getUserFromTags", routes.GetUsersFromTags)

	// called as localhost:3000/update/{id}
	router.PUT("/update/:postId", routes.UpdateUser)

	// called as localhost:3000/delete/{id}
	router.DELETE("/delete/:postId", routes.DeleteUser)

	err := router.Run("localhost:8080")
	if err != nil {
		log.Panic("failed to run router:", err)
	}
}
