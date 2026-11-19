package users

import (
	"context"
	"log/slog"

	"github.com/chaihanij/backend-challenge/app/entities"
	"github.com/chaihanij/backend-challenge/app/layers/repositories/users/models"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func (r *repo) FindUser(ctx context.Context, input *entities.UserFilter) ([]entities.User, error) {
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	filter, err := models.NewUserFilter(input)
	if err != nil {
		slog.ErrorContext(ctx, "failed to build user filter",
			slog.String("error", err.Error()))
		return nil, err
	}

	cursor, err := r.col.Find(ctx, filter)
	if err != nil {
		slog.ErrorContext(ctx, "failed to find users",
			slog.String("error", err.Error()),
			slog.Any("filter", filter))
		return nil, err
	}

	defer func(cursor *mongo.Cursor, ctx context.Context) {
		err := cursor.Close(ctx)
		if err != nil {
			slog.ErrorContext(ctx, "failed to close cursor")
		}
	}(cursor, ctx)

	var usersModel models.Users

	if err := cursor.All(ctx, &usersModel); err != nil {
		slog.ErrorContext(ctx, "failed to decode users",
			slog.String("error", err.Error()))
		return nil, err
	}

	return usersModel.ToEntity()
}
