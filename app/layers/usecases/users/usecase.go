package users

import (
	"context"
	"log/slog"

	"github.com/chaihanij/backend-challenge/app/entities"
	_usersrepo "github.com/chaihanij/backend-challenge/app/layers/repositories/users"
)

type UseCase interface {
	CountUser(ctx context.Context) (int64, error)
	CreateUser(ctx context.Context, input *entities.User) (*entities.User, error)
	DeleteUser(ctx context.Context, id string) error
	FindOneUser(ctx context.Context, id string) (*entities.User, error)
	FindUser(ctx context.Context, input *entities.UserFilter) ([]entities.User, error)
	UpdateUser(ctx context.Context, input *entities.User) (*entities.User, error)
}

type useCase struct {
	usersRepo _usersrepo.Repo
	logger    *slog.Logger
}

func New(usersRepo _usersrepo.Repo, logger *slog.Logger) UseCase {
	return &useCase{
		usersRepo: usersRepo,
		logger:    logger,
	}
}
