package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/chaihanij/backend-challenge/app/config"
	"github.com/chaihanij/backend-challenge/app/database"
	"github.com/chaihanij/backend-challenge/app/docs"
	_authhttp "github.com/chaihanij/backend-challenge/app/layers/deliveries/http/api/auth"
	_healthhttp "github.com/chaihanij/backend-challenge/app/layers/deliveries/http/api/health"
	_usershttp "github.com/chaihanij/backend-challenge/app/layers/deliveries/http/api/users"
	"github.com/chaihanij/backend-challenge/app/layers/deliveries/http/middlewares"
	"github.com/chaihanij/backend-challenge/app/layers/deliveries/workers"
	_usersrepo "github.com/chaihanij/backend-challenge/app/layers/repositories/users"
	_authusecase "github.com/chaihanij/backend-challenge/app/layers/usecases/auth"
	_usersusecase "github.com/chaihanij/backend-challenge/app/layers/usecases/users"
	"github.com/chaihanij/backend-challenge/app/util/slogx"
	"github.com/gin-gonic/gin"

	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title			Backend Challenge API
// @version			1.0.0
// @description		This is the API documentation for Backend Challenge API.
// @contact.name	chaihanij
// @contact.email	chaihanij@gmail.com
// @host			{{backendBaseUrl}}
func main() {
	logger := slogx.New(&slogx.Config{Level: slog.LevelInfo, Format: "json"})
	if err := run(logger); err != nil {
		logger.Error("server stopped with error", slog.Any("error", err))
		os.Exit(1)
	}
}

func run(logger *slog.Logger) error {
	cfg, err := config.NewConfig()
	if err != nil {
		return err
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	client, err := database.NewMongoClient(ctx, database.MongoConfig{
		URL:      cfg.MongoDBUrl,
		DB:       cfg.MongoDBName,
		User:     cfg.MongoDBUser,
		Password: cfg.MongoDBPassword,
		Timeout:  cfg.MongoDBRequestTimeout,
	})
	if err != nil {
		return err
	}
	defer func() {
		shutdownCtx, cancel := context.WithTimeout(context.Background(), cfg.MongoDBRequestTimeout)
		defer cancel()
		if err := client.Disconnect(shutdownCtx); err != nil {
			logger.Error("failed to disconnect from MongoDB", slog.Any("error", err))
		}
	}()

	userRepo := _usersrepo.New(client, cfg.MongoDBName, cfg.MongoDBRequestTimeout)
	userUseCase := _usersusecase.New(userRepo, logger)
	authUseCase := _authusecase.New(userRepo, logger, cfg.JWTSecret, cfg.ServiceName, cfg.ServiceName, cfg.JWTExpire)

	gin.SetMode(gin.ReleaseMode)
	ginEngine := gin.New()
	ginEngine.Use(gin.Recovery(), middlewares.ErrorHandler())
	ginEngine.Use(middlewares.AccessLogWithConfig(nil))
	ginEngine.Use(middlewares.CorrelationId())
	ginEngine.Use(middlewares.RequestId())
	ginEngine.Use(middlewares.TraceId())

	_healthhttp.NewHandler(ginEngine)

	_authhttp.New(ginEngine, authUseCase)

	_usershttp.New(ginEngine, userUseCase, middlewares.NewAuth(authUseCase))

	if cfg.DebugMode {
		gin.SetMode(gin.DebugMode)
		docs.SwaggerInfo.Schemes = []string{"http"}
		ginEngine.GET("/api/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	}

	server := &http.Server{
		Addr:              ":8080",
		Handler:           ginEngine,
		ReadHeaderTimeout: 5 * time.Second,
	}

	workers.StartUserCountLogger(ctx, userUseCase, logger)

	serveErr := make(chan error, 1)
	go func() {
		logger.Info("HTTP server listening", slog.String("address", server.Addr))
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			serveErr <- err
		}
	}()

	select {
	case err := <-serveErr:
		return err
	case <-ctx.Done():
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := server.Shutdown(shutdownCtx); err != nil {
		return err
	}
	return nil
}
