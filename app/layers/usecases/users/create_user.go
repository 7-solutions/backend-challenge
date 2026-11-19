package users

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/chaihanij/backend-challenge/app/entities"
	customerrors "github.com/chaihanij/backend-challenge/app/errors"
	"github.com/chaihanij/backend-challenge/app/util/pointer"
)

func (u *useCase) CreateUser(ctx context.Context, input *entities.User) (*entities.User, error) {
	start := time.Now()

	u.logCreateUserEvent(ctx, "request_received", input.Email, "", start)

	exitingUser, err := u.usersRepo.FindUser(ctx, &entities.UserFilter{
		Email: pointer.ToPointer(input.Email),
	})

	if err != nil {
		u.logCreateUserEvent(ctx, "request_failed", input.Email, "db_error", start)
		slog.ErrorContext(ctx, "failed to find user on command create_user", slog.Any("input", input))
		return nil, err
	}

	if len(exitingUser) > 0 {
		u.logCreateUserEvent(ctx, "request_failed", input.Email, "duplicate_email", start)
		return nil, customerrors.NewConflictError(
			"users",
			fmt.Sprintf("email: %s user with this email already exists", input.Email),
		)
	}

	result, err := u.usersRepo.CreateUser(ctx, input)
	if err != nil {
		u.logCreateUserEvent(ctx, "request_failed", input.Email, "db_error", start)
		slog.ErrorContext(ctx, "failed to create user on command create_user", slog.Any("input", input))
		return nil, err
	}

	u.logCreateUserEvent(ctx, "request_success", input.Email, "", start)

	return result, nil
}

func (u *useCase) logCreateUserEvent(ctx context.Context, event, email, reason string, start time.Time) {
	latency := time.Since(start).Milliseconds()

	u.logger.InfoContext(ctx, "create_user_stat",
		slog.String("cmd", "create_user"),
		slog.String("event", event), // request_received | request_failed | request_success
		slog.String("reason", reason),
		slog.Int64("latency_ms", latency),
		slog.Time("timestamp", time.Now()),
	)
}
