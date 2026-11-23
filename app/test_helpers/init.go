package test_helpers

import (
	"log/slog"
	"os"

	"github.com/chaihanij/backend-challenge/app/config"
	"github.com/chaihanij/backend-challenge/app/util/slogx"
	"github.com/gin-gonic/gin"
)

var Config *config.Config
var Logger *slog.Logger

func Init() {
	gin.SetMode(gin.TestMode)
	_ = os.Setenv("SERVICE_NAME", "backend-challenge")
	_ = os.Setenv("DEBUG_MODE", "true")
	_ = os.Setenv("MONGODB_URL", "mongodb://localhost:27017")
	_ = os.Setenv("MONGODB_NAME", "backend-challenge")
	_ = os.Setenv("MONGODB_REQUEST_TIMEOUT", "10")
	_ = os.Setenv("JWT_SECRET", "secret")
	_ = os.Setenv("JWT_EXPIRE", "3600")

	Config, _ = config.NewConfig()
	l := slogx.New(&slogx.Config{
		Level:  slog.LevelDebug,
		Format: "json",
		Writer: os.Stdout,
	})
	slog.SetDefault(l)
	Logger = l
}
