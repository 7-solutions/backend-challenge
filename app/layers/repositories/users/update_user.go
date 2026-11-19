package users

import (
	"context"
	"log/slog"

	"github.com/chaihanij/backend-challenge/app/entities"
	customerrors "github.com/chaihanij/backend-challenge/app/errors"
	"github.com/chaihanij/backend-challenge/app/layers/repositories/users/models"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func (r *repo) UpdateUser(ctx context.Context, input *entities.User) error {

	if input == nil {
		slog.ErrorContext(ctx, "input is nil")
		return customerrors.ErrNilInput
	}

	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	filter, err := models.NewUserFilter(input.ID)
	if err != nil {
		slog.ErrorContext(ctx, "failed to build user filter",
			slog.String("error", err.Error()))
		return err
	}

	update, err := models.UpdateUserModel(input)
	if err != nil {
		slog.ErrorContext(ctx, "failed to build update user model",
			slog.String("error", err.Error()),
			slog.Any("input", input))
		return err
	}

	result, err := r.col.UpdateOne(ctx, filter, bson.M{"$set": update})
	if err != nil {
		slog.ErrorContext(ctx, "failed to update user",
			slog.String("error", err.Error()),
			slog.Any("filter", filter),
			slog.Any("update", update))
		return err
	}

	if result == nil || !result.Acknowledged || result.ModifiedCount == 0 {
		slog.WarnContext(ctx, "user not updated",
			slog.String("id", input.ID))
		return customerrors.NewNotFoundError("User", input.ID)
	}

	return nil

}
