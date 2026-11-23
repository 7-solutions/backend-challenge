package auth

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	customerrors "github.com/chaihanij/backend-challenge/app/errors"
	"github.com/golang-jwt/jwt/v5"
)

func (u *useCase) ValidateToken(ctx context.Context, accessToken string) (*string, bool, error) {
	start := time.Now()
	u.logValidateTokenEvent(ctx, "request_received", "", start)

	token, err := jwt.Parse(accessToken, func(t *jwt.Token) (interface{}, error) {

		if t.Method != jwt.SigningMethodHS256 {

			u.logValidateTokenEvent(ctx, "request_failed", "invalid_signing_method", start)
			return false, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}

		return []byte(u.secret), nil
	})

	if err != nil {

		u.logValidateTokenEvent(ctx, "request_failed", "parse_error", start)
		slog.ErrorContext(ctx, "failed to validate token on command validate_token", slog.Any("error", err.Error()))
		return nil, false, err
	}

	if !token.Valid {
		u.logValidateTokenEvent(ctx, "request_failed", "invalid_token", start)
		return nil, false, customerrors.NewUnauthorized("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		u.logValidateTokenEvent(ctx, "request_failed", "claims_cast_error", start)
		return nil, false, customerrors.NewUnauthorized("invalid token claims")
	}

	iss, ok := claims["iss"].(string)
	if !ok || iss != u.iss {
		u.logValidateTokenEvent(ctx, "request_failed", "iss_cast_error", start)
		return nil, false, customerrors.NewUnauthorized("invalid token claims")
	}

	aud, ok := claims["aud"].(string)
	if !ok || aud != u.aud {
		u.logValidateTokenEvent(ctx, "request_failed", "aud_cast_error", start)
		return nil, false, customerrors.NewUnauthorized("invalid token claims")
	}

	sub, ok := claims["sub"].(string)
	if !ok {
		u.logValidateTokenEvent(ctx, "request_failed", "sub_cast_error", start)
		return nil, false, customerrors.NewUnauthorized("invalid token claims")
	}

	u.logValidateTokenEvent(ctx, "request_success", "", start)
	return &sub, true, nil
}

func (u *useCase) logValidateTokenEvent(ctx context.Context, event, reason string, start time.Time) {
	latency := time.Since(start).Milliseconds()

	u.logger.InfoContext(ctx, "validate_token_stat",
		slog.String("cmd", "validate_token"),
		slog.String("event", event),
		slog.String("reason", reason),
		slog.Int64("latency_ms", latency),
		slog.Time("timestamp", time.Now()),
	)
}
