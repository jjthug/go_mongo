package routes

import (
	"context"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go_mongo/collection"
	"go_mongo/database"
	"net/http"
	"time"
)

func DeleteUser(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var DB = database.ConnectDB()
	userId := c.Param("userId")
	var postCollection = collection.GetCollection(DB, "Users")
	defer cancel()
	objId, _ := primitive.ObjectIDFromHex(userId)
	result, err := postCollection.DeleteOne(ctx, bson.M{"id": objId})
	res := map[string]interface{}{"data": result}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err})
		return
	}
	if result.DeletedCount < 1 {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "No data to delete"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "User deleted successfully", "Data": res})
}
