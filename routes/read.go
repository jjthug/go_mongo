package routes

import (
	"context"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go_mongo/collection"
	"go_mongo/database"
	"go_mongo/model"
	"log"
	"net/http"
	"time"
)

func ReadOneUser(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var DB = database.ConnectDB()
	var userCollection = collection.GetCollection(DB, "Users")

	userId := c.Param("userId")
	var result model.User

	defer cancel()

	objId, _ := primitive.ObjectIDFromHex(userId)

	err := userCollection.FindOne(ctx, bson.M{"id": objId}).Decode(&result)

	res := map[string]interface{}{"data": result}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "success!", "Data": res})
}

type TagsRequest struct {
	QueryTags []string `json:"queryTags" binding:"required"`
}

func GetUsersFromTags(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Get a handle to the Users collection
	var DB = database.ConnectDB()
	var userCollection = collection.GetCollection(DB, "Users")

	// Bind the JSON payload to a TagsRequest struct
	var tagsRequest TagsRequest
	if err := c.ShouldBindJSON(&tagsRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Extract the queryTags parameter from the struct
	queryTags := tagsRequest.QueryTags

	// Build the MongoDB aggregation pipeline
	pipeline := []bson.M{
		bson.M{"$match": bson.M{"tags": bson.M{"$in": queryTags}}},
		bson.M{"$unwind": "$tags"},
		bson.M{"$match": bson.M{"tags": bson.M{"$in": queryTags}}},
		bson.M{"$group": bson.M{
			"_id":   "$_id",
			"name":  bson.M{"$first": "$name"},
			"count": bson.M{"$sum": 1},
		}},
		bson.M{"$sort": bson.M{"count": -1}},
		bson.M{"$limit": 10},
	}

	// Execute the aggregation pipeline
	cursor, err := userCollection.Aggregate(ctx, pipeline)
	if err != nil {
		log.Fatal(err)
	}
	defer cursor.Close(ctx)

	// Iterate over the cursor to populate the result variable
	var result []model.User
	for cursor.Next(ctx) {
		var user model.User
		if err := cursor.Decode(&user); err != nil {
			log.Fatal(err)
		}
		result = append(result, user)
	}

	res := map[string]interface{}{"data": result}
	c.JSON(http.StatusOK, gin.H{"message": "success!", "Data": res})
}
