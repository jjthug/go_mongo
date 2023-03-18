package routes

import (
	"context"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go_mongo/collection"
	"go_mongo/database"
	"go_mongo/model"
	"log"
	"net/http"
	"time"
)

func CreateUser(c *gin.Context) {
	var DB = database.ConnectDB()
	var postCollection = collection.GetCollection(DB, "Users")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	user := new(model.User)
	defer cancel()

	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err})
		log.Fatal(err)
		return
	}

	postPayload := model.User{
		ID:   primitive.NewObjectID(),
		Name: user.Name,
		Tags: user.Tags,
	}

	result, err := postCollection.InsertOne(ctx, postPayload)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Added user successfully", "Data": map[string]interface{}{"data": result}})
}
