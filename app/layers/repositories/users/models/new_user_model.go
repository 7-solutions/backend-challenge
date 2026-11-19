package models

import (
	"github.com/chaihanij/backend-challenge/app/entities"
	"github.com/chaihanij/backend-challenge/app/util/password"
	"github.com/chaihanij/backend-challenge/app/util/pointer"
	"go.mongodb.org/mongo-driver/v2/bson"
)

// NewUserModel builds the MongoDB document for a new user, generating a
// fresh ObjectID and hashing the plaintext password carried by u.
func NewUserModel(u *entities.User) (*User, error) {
	if u == nil {
		return nil, ErrNilUser
	}

	id := bson.NewObjectID()
	pass, err := password.HashPassword(u.Password)
	if err != nil {
		return nil, err
	}

	return &User{
		ID:        id,
		Name:      u.Name,
		Email:     u.Email,
		Password:  pass,
		CreatedAt: pointer.ToPointer(id.Timestamp()),
		UpdatedAt: pointer.ToPointer(id.Timestamp()),
	}, nil
}
