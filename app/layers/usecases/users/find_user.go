package users

import (
	"context"
	"log/slog"
	"time"

	"github.com/chaihanij/backend-challenge/app/entities"
)

func (u *useCase) FindUser(ctx context.Context, input *entities.UserFilter) ([]entities.User, error) {
	start := time.Now().UnixMilli()
	u.logFindUserEvent(ctx, "request_received", "", start)
	results, err := u.usersRepo.FindUser(ctx, input)
	if err != nil {

		u.logFindUserEvent(ctx, "request_failed", "db_error", start)
		slog.ErrorContext(ctx, "failed to find user on command find_user", slog.Any("input", input))
		return nil, err
	}
	u.logFindUserEvent(ctx, "request_success", "", start)
	return results, nil
}

func (u *useCase) logFindUserEvent(ctx context.Context, event, reason string, start int64) {
	latency := time.Since(time.UnixMilli(start)).Milliseconds()
	u.logger.InfoContext(ctx, "find_user_stat",
		slog.String("cmd", "find_user"),
		slog.String("event", event),
		slog.String("reason", reason),
		slog.Int64("latency_ms", latency),
		slog.Time("timestamp", time.Now()),
	)
}
