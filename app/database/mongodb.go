package database

import (
	"context"
	"log/slog"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type MongoConfig struct {
	URL      string
	DB       string
	User     string
	Password string
	Timeout  time.Duration
}

func NewMongoClient(ctx context.Context, cfg MongoConfig) (*mongo.Client, error) {

	clientOptions := options.Client().
		ApplyURI(cfg.URL).
		SetServerSelectionTimeout(cfg.Timeout).
		SetConnectTimeout(cfg.Timeout)

	if cfg.User != "" {
		clientOptions.SetAuth(options.Credential{
			Username:   cfg.User,
			Password:   cfg.Password,
			AuthSource: cfg.DB,
		})
	}

	client, err := mongo.Connect(clientOptions)
	if err != nil {
		slog.Error("MongoDB connection failed",
			slog.String("url", cfg.URL),
			slog.String("error", err.Error()))
		return nil, err
	}

	pingCtx, cancel := context.WithTimeout(ctx, cfg.Timeout)
	defer cancel()

	if err := client.Ping(pingCtx, nil); err != nil {
		slog.Error("MongoDB ping failed",
			slog.String("url", cfg.URL),
			slog.String("error", err.Error()))
		return nil, err
	}

	slog.Info("Connected to MongoDB", slog.String("url", cfg.URL))
	return client, nil
}
