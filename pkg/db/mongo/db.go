package mongo

import (
	"context"
	"fmt"
	"log"

	"go_mongo/internal/config"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DB struct {
	*mongo.Database
	Close func(ctx context.Context) error
}

func Connect(ctx context.Context, cfg config.MongoDBConfig) (db DB, err error) {
	connectionString := "mongodb://" + cfg.Address
	client, err := mongo.NewClient(options.Client().ApplyURI(connectionString))
	if err != nil {
		return db, fmt.Errorf("create mongo client err: %w", err)
	}

	if err = client.Connect(ctx); err != nil {
		return db, fmt.Errorf("connect to mongo err: %w", err)
	}

	connectDBCtx, cancel := context.WithTimeout(ctx, cfg.ConnectTimeout)
	defer cancel()

	if err := client.Ping(connectDBCtx, nil); err != nil {
		return db, fmt.Errorf("ping error: %w", err)
	}

	db = DB{client.Database(cfg.Name), client.Disconnect}
	log.Println("Connected to mongoDB")

	return db, err
}
