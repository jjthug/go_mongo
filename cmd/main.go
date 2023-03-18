package main

import (
	"context"
	"log"
	"net/http"

	"go_mongo/internal/config"
	"go_mongo/internal/delivery/rest"
	"go_mongo/internal/repository"
	"go_mongo/pkg/db/mongo"

	"github.com/caarlos0/env/v7"
)

func main() {
	ctx := context.Background()
	cfg := config.Config{}
	if err := env.Parse(&cfg); err != nil {
		log.Fatalf("parsing config %+v", err)
	}

	db, err := mongo.Connect(ctx, cfg.Mongo)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close(ctx)

	repo := repository.NewRepo(db.Database)
	h := rest.NewHandler(repo)

	log.Printf("start listening: %s\n", cfg.Rest.Address)
	if err := http.ListenAndServe(cfg.Rest.Address, h.Router()); err != nil {
		log.Printf("stop listening: %s\n", cfg.Rest.Address)
	}
}
