package workers

import (
	"context"
	"log/slog"
	"time"
)

const userCountLogInterval = 10 * time.Second

// UserCounter supplies the total number of users to the background logger.
type UserCounter interface {
	CountUser(context.Context) (int64, error)
}

// StartUserCountLogger logs the current user count every 10 seconds until ctx
// is canceled.
func StartUserCountLogger(ctx context.Context, counter UserCounter, logger *slog.Logger) {
	go func() {
		ticker := time.NewTicker(userCountLogInterval)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				count, err := counter.CountUser(ctx)
				if err != nil {
					logger.ErrorContext(ctx, "failed to log total user count", slog.Any("error", err))
					continue
				}
				logger.InfoContext(ctx, "total users", slog.Int64("count", count))
			}
		}
	}()
}
