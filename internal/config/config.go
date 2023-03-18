package config

import (
	"time"

	_ "github.com/caarlos0/env/v7"
)

type (
	Config struct {
		Mongo MongoDBConfig `envPrefix:"MONGO_DB_"`
		Rest  RestConfig    `envPrefix:"REST_"`
	}
	MongoDBConfig struct {
		Address        string        `env:"ADDRESS" envDefault:"127.0.0.1:27017"`
		Name           string        `env:"NAME" envDefault:"myGoappDB"`
		ConnectTimeout time.Duration `env:"CONNECT_TIMEOUT" envDefault:"2s"`
	}
	RestConfig struct {
		Address string `env:"ADDRESS" envDefault:"localhost:8080"`
	}
)
