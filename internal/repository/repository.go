package repository

import (
	"context"
	"errors"
	"fmt"

	"go_mongo/internal/repository/model"
	service "go_mongo/pkg"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	UserCollection = "Users"
)

type Repo struct {
	client *mongo.Database
}

func NewRepo(client *mongo.Database) *Repo {
	return &Repo{client: client}
}

func (r Repo) CreateUser(ctx context.Context, name string, tags []string) (ID string, err error) {
	payload := model.User{
		ID:   primitive.NewObjectID(),
		Name: name,
		Tags: tags,
	}
	res, err := r.client.Collection(UserCollection).InsertOne(ctx, payload)
	if err != nil {
		return "", fmt.Errorf("insert err: %w", err)
	}
	id := res.InsertedID.(primitive.ObjectID).Hex()
	return id, nil
}

func (r Repo) GetUserByID(ctx context.Context, id string) (*service.User, error) {
	objId, _ := primitive.ObjectIDFromHex(id)
	var u model.User
	err := r.client.Collection(UserCollection).FindOne(ctx, bson.M{"id": objId}).Decode(&u)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, nil
	}
	return &service.User{
		ID:   u.ID.Hex(),
		Name: u.Name,
		Tags: u.Tags,
	}, nil
}

func (r Repo) GetUsersByTags(ctx context.Context, tags []string) ([]service.User, error) {
	// Build the MongoDB aggregation pipeline
	pipeline := []bson.M{
		{"$match": bson.M{"tags": bson.M{"$in": tags}}},
		{"$unwind": "$tags"},
		{"$match": bson.M{"tags": bson.M{"$in": tags}}},
		{"$group": bson.M{
			"_id":   "$_id",
			"name":  bson.M{"$first": "$name"},
			"count": bson.M{"$sum": 1},
		}},
		{"$sort": bson.M{"count": -1}},
		{"$limit": 10},
	}

	// Execute the aggregation pipeline
	cursor, err := r.client.Collection(UserCollection).Aggregate(ctx, pipeline)
	if err != nil {
		return nil, fmt.Errorf("aggregation pipeline err: %w", err)
	}
	defer cursor.Close(ctx)

	// Iterate over the cursor to populate the result variable
	var (
		result []service.User
		user   model.User
	)
	for cursor.Next(ctx) {
		if err := cursor.Decode(&user); err != nil {
			return nil, fmt.Errorf("decoding user err: %w", err)
		}
		result = append(result, service.User{
			ID:   user.ID.Hex(),
			Name: user.Name,
			Tags: user.Tags,
		})
	}
	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("cursor err: %w", err)
	}

	return result, nil
}

func (r Repo) UpdateUser(ctx context.Context, svcUser service.User) (bool, error) {
	objId, err := primitive.ObjectIDFromHex(svcUser.ID)
	if err != nil {
		return false, fmt.Errorf("unacceptable user id: %s", svcUser.ID)
	}
	user := model.User{
		ID:   objId,
		Name: svcUser.Name,
		Tags: svcUser.Tags,
	}
	result, err := r.client.Collection(UserCollection).UpdateOne(ctx, bson.M{"id": objId}, bson.M{"$set": user})
	if err != nil {
		return false, fmt.Errorf("update user err: %w", err)
	}
	return result.MatchedCount > 0, nil
}

func (r Repo) DeleteUser(ctx context.Context, id string) (bool, error) {
	result, err := r.client.Collection(UserCollection).DeleteOne(ctx, bson.M{"id": id})
	if err != nil {
		return false, fmt.Errorf("delete user err: %w", err)
	}

	return result.DeletedCount >= 1, nil
}
