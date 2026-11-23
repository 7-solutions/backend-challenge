package models

import (
	"time"

	"github.com/chaihanij/backend-challenge/app/entities"
	"github.com/chaihanij/backend-challenge/app/util/password"
	"go.mongodb.org/mongo-driver/v2/bson"
)

// UpdateUserModel builds a MongoDB update document containing an updated_at
// timestamp plus any non-empty field on input. It returns ErrNilUser when
// input is nil.
func UpdateUserModel(input *entities.User) (*bson.D, error) {
	if input == nil {
		return nil, ErrNilUser
	}

	update := bson.D{
		{Key: "updated_at", Value: time.Now()},
	}

	if input.Email != "" {
		update = append(update, bson.E{Key: "email", Value: input.Email})
	}

	if input.Name != "" {
		update = append(update, bson.E{Key: "name", Value: input.Name})
	}

	if input.Password != "" {
		pass, err := password.HashPassword(input.Password)
		if err != nil {
			return nil, err
		}
		update = append(update, bson.E{Key: "password", Value: pass})
	}

	return &update, nil
}
