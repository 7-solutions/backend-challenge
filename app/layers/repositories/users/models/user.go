package models

import (
	"errors"
	"time"

	"github.com/chaihanij/backend-challenge/app/entities"
	"github.com/chaihanij/backend-challenge/app/util/pointer"
	"go.mongodb.org/mongo-driver/v2/bson"
)

// ErrNilUser is returned by model constructors that require a non-nil
// *entities.User input.
var ErrNilUser = errors.New("user entity is nil")

// User is the MongoDB document representation of a user.
type User struct {
	ID        bson.ObjectID `bson:"_id"`
	Name      string        `bson:"name"`
	Email     string        `bson:"email"`
	Password  string        `bson:"password"`
	CreatedAt *time.Time    `bson:"created_at"`
	UpdatedAt *time.Time    `bson:"updated_at"`
}

// ToEntity converts the document into its domain entity representation.
func (u User) ToEntity() (*entities.User, error) {
	return &entities.User{
		ID:        u.ID.Hex(),
		Name:      u.Name,
		Email:     u.Email,
		Password:  u.Password,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}, nil
}

// Users is a collection of User documents.
type Users []User

// ToEntity converts each document in the collection into its domain entity representation.
func (us Users) ToEntity() ([]entities.User, error) {
	results := make([]entities.User, 0, len(us))
	for _, u := range us {
		result, err := u.ToEntity()
		if err != nil {

			return nil, err
		}
		results = append(results, pointer.ToValue(result))
	}
	return results, nil
}
