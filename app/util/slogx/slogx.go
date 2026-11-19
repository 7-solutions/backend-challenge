package slogx

import (
	"context"
	"io"
	"log/slog"
	"os"

	"github.com/chaihanij/backend-challenge/app/constants"
)

// Config holds the configuration for creating a logger.
type Config struct {
	// Level sets the minimum log level (default: slog.LevelInfo).
	Level slog.Level

	// Format selects the output encoding: "json" or "text" (default).
	Format string

	// Writer is the output destination. Defaults to os.Stdout when nil.
	Writer io.Writer
}

// New creates a new [slog.Logger] wrapped with the slogx handler that
// enriches records with context IDs and OpenTelemetry / Datadog trace info.
func New(cfg *Config) *slog.Logger {
	w := cfg.Writer
	if w == nil {
		w = os.Stdout
	}

	opts := &slog.HandlerOptions{Level: cfg.Level}

	var base slog.Handler
	if cfg.Format == "json" {
		base = slog.NewJSONHandler(w, opts)
	} else {
		base = slog.NewTextHandler(w, opts)
	}

	return slog.New(NewHandler(base))
}

// SetDefault calls [New] and registers the result as the global default logger.
func SetDefault(cfg *Config) *slog.Logger {
	logger := New(cfg)
	slog.SetDefault(logger)
	return logger
}

// WithRequestID returns a new context that carries the given request ID.
func WithRequestID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, constants.CtxKeyRequestID, id)
}

// WithCorrelationID returns a new context that carries the given correlation ID.
func WithCorrelationID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, constants.CtxKeyCorrelationID, id)
}

// WithTraceID returns a new context that carries the given trace ID.
func WithTraceID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, constants.CtxKeyTraceID, id)
}
