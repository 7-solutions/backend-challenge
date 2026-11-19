package users

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	customerrors "github.com/chaihanij/backend-challenge/app/errors"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func (u *useCase) DeleteUser(ctx context.Context, id string) error {
	start := time.Now().UnixMilli()
	u.logDeleteUserEvent(ctx, "request_received", "", start)

	_, err := u.usersRepo.FindOneUser(ctx, id)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {

			u.logDeleteUserEvent(ctx, "request_failed", "user_not_found", start)
			return customerrors.NewNotFoundError("users", fmt.Sprintf("user with id: %s not found", id))
		}
		u.logDeleteUserEvent(ctx, "request_failed", "db_error", start)
		slog.ErrorContext(ctx, "failed to find user on command delete_user", slog.Any("id", id))
		return err
	}

	err = u.usersRepo.DeleteUser(ctx, id)
	if err != nil {
		u.logDeleteUserEvent(ctx, "request_failed", "db_error", start)
		slog.ErrorContext(ctx, "failed to delete user on command delete_user", slog.Any("id", id))
		return err
	}

	u.logDeleteUserEvent(ctx, "request_success", "", start)
	return nil
}

func (u *useCase) logDeleteUserEvent(ctx context.Context, event, reason string, start int64) {
	latency := time.Since(time.UnixMilli(start)).Milliseconds()
	u.logger.InfoContext(ctx, "delete_user_stat",
		slog.String("cmd", "delete_user"),
		slog.String("event", event),
		slog.String("reason", reason),
		slog.Int64("latency_ms", latency),
		slog.Time("timestamp", time.Now()),
	)
}
