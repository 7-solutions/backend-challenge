package users

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/chaihanij/backend-challenge/app/entities"
	customerrors "github.com/chaihanij/backend-challenge/app/errors"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func (u *useCase) FindOneUser(ctx context.Context, id string) (*entities.User, error) {
	start := time.Now().UnixMilli()

	u.logFindOneUserEvent(ctx, "request_received", "", start)
	result, err := u.usersRepo.FindOneUser(ctx, id)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {

			u.logFindOneUserEvent(ctx, "request_success", "user_not_found", start)
			return nil, customerrors.NewNotFoundError("users", fmt.Sprintf("user with id: %s not found", id))
		}
		u.logFindOneUserEvent(ctx, "request_failed", "db_error", start)
		slog.ErrorContext(ctx, "failed to find user on command find_one_user", slog.Any("id", id))
		return nil, err
	}
	u.logFindOneUserEvent(ctx, "request_success", "", start)
	return result, nil
}

func (u *useCase) logFindOneUserEvent(ctx context.Context, event, reason string, start int64) {
	latency := time.Since(time.UnixMilli(start)).Milliseconds()
	u.logger.InfoContext(ctx, "find_one_user_stat",
		slog.String("cmd", "find_one_user"),
		slog.String("event", event),
		slog.String("reason", reason),
		slog.Int64("latency_ms", latency),
		slog.Time("timestamp", time.Now()),
	)
}
