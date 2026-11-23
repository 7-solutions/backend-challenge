package users

import (
	"context"
	"log/slog"

	customerrors "github.com/chaihanij/backend-challenge/app/errors"
	"github.com/chaihanij/backend-challenge/app/layers/repositories/users/models"
)

func (r *repo) DeleteUser(ctx context.Context, id string) error {

	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	filter, err := models.NewUserFilter(id)
	if err != nil {
		slog.ErrorContext(ctx, "failed to build user filter",
			slog.String("error", err.Error()))
		return err
	}

	result, err := r.col.DeleteOne(ctx, filter)
	if err != nil {
		slog.ErrorContext(ctx, "failed to delete user",
			slog.String("error", err.Error()),
			slog.Any("filter", filter))
		return err
	}

	if result == nil || !result.Acknowledged || result.DeletedCount == 0 {
		slog.WarnContext(ctx, "user not deleted",
			slog.String("id", id))
		return customerrors.NewNotFoundError("user", id)
	}

	return nil
}
