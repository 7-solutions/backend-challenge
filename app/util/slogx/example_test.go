package slogx_test

import (
	"bytes"
	"context"
	"log/slog"

	"github.com/chaihanij/backend-challenge/app/util/slogx"
)

func ExampleNew() {
	var buf bytes.Buffer
	logger := slogx.New(&slogx.Config{
		Format: "json",
		Level:  slog.LevelInfo,
		Writer: &buf,
	})

	ctx := slogx.WithRequestID(context.Background(), "req-123")
	ctx = slogx.WithCorrelationID(ctx, "corr-456")

	// requestId and correlationId are added automatically from ctx.
	logger.InfoContext(ctx, "user logged in", slog.String("userId", "u-1"))
}

func ExampleStructToAny() {
	type User struct {
		Name string `json:"name"`
	}

	logger := slog.Default()
	ctx := context.Background()

	// StructToAny converts a struct to []any suitable for slog variadic methods.
	user := User{Name: "Alice"}
	logger.InfoContext(ctx, "created user", slogx.StructToAny(user)...)
}
