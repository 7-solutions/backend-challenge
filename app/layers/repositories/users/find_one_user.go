package users

import (
	"context"
	"log/slog"

	"github.com/chaihanij/backend-challenge/app/entities"
	"github.com/chaihanij/backend-challenge/app/layers/repositories/users/models"
)

func (r *repo) FindOneUser(ctx context.Context, id string) (*entities.User, error) {
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	filter, err := models.NewUserFilter(id)
	if err != nil {
		slog.ErrorContext(ctx, "failed to build user filter",
			slog.String("id", id),
			slog.String("error", err.Error()))

		return nil, err
	}

	res := r.col.FindOne(ctx, filter)
	var result models.User
	if err := res.Decode(&result); err != nil {

		slog.ErrorContext(ctx, "failed to find user",
			slog.String("error", err.Error()),
			slog.Any("filter", filter))
		return nil, err
	}

	return result.ToEntity()
}
