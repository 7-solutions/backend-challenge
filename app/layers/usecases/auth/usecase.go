package auth

import (
	"context"
	"log/slog"
	"time"

	"github.com/chaihanij/backend-challenge/app/layers/repositories/users"
)

type UseCase interface {
	Login(ctx context.Context, email, password string) (*string, error)
	ValidateToken(ctx context.Context, accessToken string) (*string, bool, error)
}
type useCase struct {
	usersRepo users.Repo
	logger    *slog.Logger
	secret    string
	iss       string
	aud       string
	expire    time.Duration
}

func New(usersRepo users.Repo, logger *slog.Logger, secret string, iss string, aud string, expire time.Duration) UseCase {
	return &useCase{
		usersRepo: usersRepo,
		logger:    logger,
		secret:    secret,
		iss:       iss,
		aud:       aud,
		expire:    expire,
	}
}
