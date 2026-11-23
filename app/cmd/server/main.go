package main

import (
	"context"
	"errors"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/chaihanij/backend-challenge/app/config"
	"github.com/chaihanij/backend-challenge/app/database"
	"github.com/chaihanij/backend-challenge/app/docs"
	usersv1 "github.com/chaihanij/backend-challenge/app/gen/proto/users/v1"
	_usersgrpc "github.com/chaihanij/backend-challenge/app/layers/deliveries/grpc/users"
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
	"google.golang.org/grpc"

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

func startHTTPServer(server *http.Server, logger *slog.Logger, serveErr chan<- error) {
	go func() {
		logger.Info("HTTP server listening", slog.String("address", server.Addr))
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			serveErr <- err
		}
	}()
}

func startGRPCServer(address string, service usersv1.UserServiceServer, logger *slog.Logger, serveErr chan<- error) (*grpc.Server, error) {
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return nil, err
	}
	server := grpc.NewServer()
	usersv1.RegisterUserServiceServer(server, service)
	go func() {
		logger.Info("gRPC server listening", slog.String("address", listener.Addr().String()))
		if err := server.Serve(listener); err != nil {
			serveErr <- err
		}
	}()
	return server, nil
}

func startWorkers(ctx context.Context, userUseCase _usersusecase.UseCase, logger *slog.Logger) {
	workers.StartUserCountLogger(ctx, userUseCase, logger)
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

	serveErr := make(chan error, 2)
	startHTTPServer(server, logger, serveErr)
	grpcServer, err := startGRPCServer(":50051", _usersgrpc.NewService(userUseCase), logger, serveErr)
	if err != nil {
		return err
	}
	startWorkers(ctx, userUseCase, logger)

	select {
	case err := <-serveErr:
		return err
	case <-ctx.Done():
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	grpcServer.GracefulStop()
	if err := server.Shutdown(shutdownCtx); err != nil {
		return err
	}
	return nil
}
