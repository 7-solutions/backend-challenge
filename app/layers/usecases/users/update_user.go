package users

import (
	"context"
	"log/slog"
	"time"

	"github.com/chaihanij/backend-challenge/app/entities"
)

func (u *useCase) UpdateUser(ctx context.Context, input *entities.User) (*entities.User, error) {
	start := time.Now().UnixMilli()

	u.logUpdateUserEvent(ctx, "request_received", "", start)
	_, err := u.usersRepo.FindOneUser(ctx, input.ID)
	if err != nil {

		u.logUpdateUserEvent(ctx, "request_failed", "db_error", start)
		slog.ErrorContext(ctx, "failed to find one user before update user on command update_user", slog.Any("id", input.ID))
		return nil, err
	}

	err = u.usersRepo.UpdateUser(ctx, input)
	if err != nil {
		u.logUpdateUserEvent(ctx, "request_failed", "db_error", start)
		slog.ErrorContext(ctx, "failed to update user on command update_user", slog.Any("input", input))
		return nil, err
	}

	result, err := u.usersRepo.FindOneUser(ctx, input.ID)
	if err != nil {

		u.logUpdateUserEvent(ctx, "request_failed", "db_error", start)
		slog.ErrorContext(ctx, "failed to find one user after update user on command update_user", slog.Any("id", input.ID))
		return nil, err
	}

	u.logUpdateUserEvent(ctx, "request_success", "", start)
	return result, nil
}

func (u *useCase) logUpdateUserEvent(ctx context.Context, event, reason string, start int64) {
	latency := time.Since(time.UnixMilli(start)).Milliseconds()
	u.logger.InfoContext(ctx, "update_user_stat",
		slog.String("cmd", "update_user"),
		slog.String("event", event),
		slog.String("reason", reason),
		slog.Int64("latency_ms", latency),
		slog.Time("timestamp", time.Now()),
	)
}
