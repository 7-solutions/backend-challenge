package mongo

import (
	"context"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func NewLoggingMiddleware(client *mongo.Client, dbName string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()
		err := c.Next()
		duration := time.Since(start)

		coll := client.Database(dbName).Collection("logs")

		method := c.Method()
		path := c.Path()
		status := c.Response().StatusCode()
		timestamp := time.Now()
		durationStr := duration.String()

		go func() {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			entry := bson.M{
				"method":    method,
				"path":      path,
				"status":    status,
				"duration":  durationStr,
				"timestamp": timestamp,
			}
			if _, err := coll.InsertOne(ctx, entry); err != nil {
				log.Printf("logger insert error: %v", err)
			}
		}()

		log.Printf("%s %s → %d (%s)", method, path, status, duration)

		return err
	}
}
