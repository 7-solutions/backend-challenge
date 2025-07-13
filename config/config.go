package config

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Config struct {
	MongoURI   string
	DBName     string
	JWTSecret  string
	ServerAddr string
}

func LoadConfig() (*Config, *mongo.Client) {
	_ = godotenv.Load()

	cfg := &Config{
		MongoURI:   os.Getenv("MONGO_URI"),
		DBName:     os.Getenv("MONGO_DB"),
		JWTSecret:  os.Getenv("JWT_SECRET"),
		ServerAddr: os.Getenv("SERVER_ADDR"),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.MongoURI))
	if err != nil {
		log.Fatalf("MongoDB connect error: %v", err)
	}

	return cfg, client
}
