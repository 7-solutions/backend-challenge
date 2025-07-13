// @title        Backend Golang User API
// @version      1.0
// @description  A simple RESTful user-management API in Go + MongoDB + JWT
// @host         localhost:3000
// @BasePath     /api/v1
// @schemes      http
package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/testGolang/backend-challenge/config"
	"github.com/testGolang/backend-challenge/delivery/router"
)

func main() {
	cfg, client := config.LoadConfig()
	defer func() {
		_ = client.Disconnect(context.Background())
	}()

	app := router.Setup(cfg, client)

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-quit
		log.Println("Shutting down server...")
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = app.Shutdown()
		client.Disconnect(ctx)
	}()

	log.Printf("Server listening on %s\n", cfg.ServerAddr)
	log.Fatal(app.Listen(cfg.ServerAddr))
}
