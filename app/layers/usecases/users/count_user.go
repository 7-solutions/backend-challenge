package users

import (
	"context"
	"log/slog"
	"time"
)

func (u *useCase) CountUser(ctx context.Context) (int64, error) {
	start := time.Now().UnixMilli()
	u.logCountUserEvent(ctx, "request_received", "", time.UnixMilli(start))

	count, err := u.usersRepo.CountUser(ctx)
	if err != nil {

		u.logCountUserEvent(ctx, "request_failed", "db_error", time.UnixMilli(start))
		slog.ErrorContext(ctx, "failed to count user on command count_user", slog.String("error", err.Error()))
		return 0, err
	}
	u.logCountUserEvent(ctx, "request_success", "", time.UnixMilli(start))
	return count, nil
}

func (u *useCase) logCountUserEvent(ctx context.Context, event, reason string, start time.Time) {
	latency := time.Since(start).Milliseconds()

	u.logger.InfoContext(ctx, "count_user_stat",
		slog.String("cmd", "count_user"),
		slog.String("event", event),
		slog.String("reason", reason),
		slog.Int64("latency_ms", latency),
		slog.Time("timestamp", time.Now()),
	)
}
