package router

import (
	"log"
	"time"

	swagger "github.com/arsmn/fiber-swagger/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/testGolang/backend-challenge/application/usecases"
	"github.com/testGolang/backend-challenge/config"
	"github.com/testGolang/backend-challenge/delivery/handlers"
	"github.com/testGolang/backend-challenge/delivery/middleware"
	_ "github.com/testGolang/backend-challenge/docs"
	logmw "github.com/testGolang/backend-challenge/infrastructure/mongo"
	"go.mongodb.org/mongo-driver/mongo"
)

func Setup(cfg *config.Config, client *mongo.Client) *fiber.App {
	app := fiber.New()
	app.Use(logmw.NewLoggingMiddleware(client, cfg.DBName))
	// serve Swagger UI at /swagger/index.html
	app.Get("/swagger/*", swagger.HandlerDefault)

	// Mongo repo + use case + handler
	userRepo := logmw.NewUserRepo(client, cfg.DBName)
	userUC := usecases.NewUserUseCase(userRepo)
	userH := handlers.NewUserHandler(userUC)

	api := app.Group("/api/v1")
	api.Post("/users/register", userH.Register)
	api.Post("/users/login", userH.Login)
	protected := api.Group("/", middleware.JWTMiddleware())
	protected.Get("/users", userH.List)
	protected.Get("/users/:id", userH.Get)
	protected.Put("/users/:id", userH.Update)
	protected.Delete("/users/:id", userH.Delete)

	//start background goroutine
	go func() {
		ticker := time.NewTicker(10 * time.Second)
		for range ticker.C {
			count, _ := userRepo.Count()
			log.Printf("Total users: %d\n", count)
		}
	}()

	return app
}
