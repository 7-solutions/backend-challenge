package auth

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/chaihanij/backend-challenge/app/entities"
	customerrors "github.com/chaihanij/backend-challenge/app/errors"
	utilpassword "github.com/chaihanij/backend-challenge/app/util/password"
	"github.com/chaihanij/backend-challenge/app/util/pointer"
	"github.com/golang-jwt/jwt/v5"
)

func (u *useCase) Login(ctx context.Context, email, password string) (*string, error) {
	start := time.Now()
	u.logLoginEvent(ctx, "request_received", "", start)

	users, err := u.usersRepo.FindUser(ctx, &entities.UserFilter{Email: pointer.ToPointer(email)})
	if err != nil {

		u.logLoginEvent(ctx, "request_failed", "db_error", start)
		slog.ErrorContext(ctx, "failed to find user on command login", slog.Any("email", email))
		return nil, err
	}

	if len(users) == 0 {

		u.logLoginEvent(ctx, "request_success", "user_not_found", start)
		return nil, customerrors.NewNotFoundError("user", fmt.Sprintf("user with email: %s not found", email))
	}

	user := users[0]
	valid := utilpassword.CheckPasswordHash(password, user.Password)
	if !valid {

		u.logLoginEvent(ctx, "request_success", "wrong_password", start)
		return nil, customerrors.NewUnauthorized("wrong password")
	}
	now := time.Now()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iss":   u.iss,
		"sub":   user.ID,
		"aud":   u.aud,
		"exp":   now.Add(u.expire).Unix(),
		"iat":   now.Unix(),
		"nbf":   now.Unix(),
		"email": user.Email,
	})

	tokenString, err := token.SignedString([]byte(u.secret))
	if err != nil {
		u.logLoginEvent(ctx, "request_failed", "token_signing_error", start)
		slog.ErrorContext(ctx, "failed to sign token on command login", slog.Any("email", email))
		return nil, err
	}

	u.logLoginEvent(ctx, "request_success", "", start)
	return &tokenString, nil
}

func (u *useCase) logLoginEvent(ctx context.Context, event, reason string, start time.Time) {
	latency := time.Since(start).Milliseconds()

	u.logger.InfoContext(ctx, "login_stat",
		slog.String("cmd", "login"),
		slog.String("event", event),
		slog.String("reason", reason),
		slog.Int64("latency_ms", latency),
		slog.Time("timestamp", time.Now()),
	)
}
