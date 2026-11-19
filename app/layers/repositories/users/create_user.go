package users

import (
	"context"
	"errors"
	"log/slog"

	"github.com/chaihanij/backend-challenge/app/entities"
	"github.com/chaihanij/backend-challenge/app/layers/repositories/users/models"
)

func (r *repo) CreateUser(ctx context.Context, input *entities.User) (*entities.User, error) {
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	// map entity → model
	newUser, err := models.NewUserModel(input)
	if err != nil {
		slog.ErrorContext(ctx, "failed to build user model",
			slog.String("error", err.Error()))
		return nil, err
	}

	result, err := r.col.InsertOne(ctx, newUser)
	if err != nil {
		slog.ErrorContext(ctx, "failed to insert user",
			slog.String("error", err.Error()),
			slog.Any("user", newUser))
		return nil, err
	}

	if result.InsertedID == nil {
		slog.WarnContext(ctx, "mongo returned nil InsertedID")
		return nil, errors.New("mongo returned nil InsertedID")
	}
	return newUser.ToEntity()
}
