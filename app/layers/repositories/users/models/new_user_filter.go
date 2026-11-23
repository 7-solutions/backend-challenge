package models

import (
	"errors"

	"github.com/chaihanij/backend-challenge/app/entities"
	"go.mongodb.org/mongo-driver/v2/bson"
)

var ErrInvalidFilterInput = errors.New("invalid input type for user filter")
var ErrEmptyFilter = errors.New("no valid fields provided for user filter")

// NewUserFilter builds a MongoDB filter document from one of:
//   - string: a user ObjectID hex string, filtered by _id.
//   - *entities.User: non-zero fields (ID, Email, Name) are added to the filter.
//   - *entities.UserFilter: non-nil, non-empty pointer fields are added to the filter.
func NewUserFilter(input any) (*bson.M, error) {
	if input == nil {
		return nil, ErrInvalidFilterInput
	}

	switch v := input.(type) {

	// Case 1: string → treat as ObjectID
	case string:
		return newFilterFromIDHex(v)

	// Case 2: *entities.User → build filter from entity fields
	case *entities.User:
		if v == nil {
			return nil, ErrInvalidFilterInput
		}

		return newFilterFromUser(v)

	// Case 3: *entities.UserFilter → pointer fields
	case *entities.UserFilter:
		if v == nil {
			return nil, ErrInvalidFilterInput
		}

		return newFilterFromUserFilter(v)

	default:
		return nil, ErrInvalidFilterInput
	}
}

func newFilterFromIDHex(idHex string) (*bson.M, error) {
	id, err := bson.ObjectIDFromHex(idHex)
	if err != nil {
		return nil, err
	}

	filter := bson.M{"_id": id}
	return &filter, nil
}

func newFilterFromUser(user *entities.User) (*bson.M, error) {
	filter := bson.M{}

	if err := addIDIfNotEmpty(filter, user.ID); err != nil {
		return nil, err
	}

	addFieldIfNotEmpty(filter, "email", user.Email)
	addFieldIfNotEmpty(filter, "name", user.Name)

	return finalizeFilter(filter)
}

func newFilterFromUserFilter(userFilter *entities.UserFilter) (*bson.M, error) {
	filter := bson.M{}

	if err := addIDPointerIfNotEmpty(filter, userFilter.ID); err != nil {
		return nil, err
	}

	addFieldPointerIfNotEmpty(filter, "email", userFilter.Email)
	addFieldPointerIfNotEmpty(filter, "name", userFilter.Name)

	return finalizeFilter(filter)
}

func addIDIfNotEmpty(filter bson.M, idHex string) error {
	if idHex == "" {
		return nil
	}

	id, err := bson.ObjectIDFromHex(idHex)
	if err != nil {
		return err
	}

	filter["_id"] = id
	return nil
}

func addIDPointerIfNotEmpty(filter bson.M, idHex *string) error {
	if idHex == nil {
		return nil
	}

	return addIDIfNotEmpty(filter, *idHex)
}

func addFieldIfNotEmpty(filter bson.M, key, value string) {
	if value == "" {
		return
	}

	filter[key] = value
}

func addFieldPointerIfNotEmpty(filter bson.M, key string, value *string) {
	if value == nil {
		return
	}

	addFieldIfNotEmpty(filter, key, *value)
}

func finalizeFilter(filter bson.M) (*bson.M, error) {
	if len(filter) == 0 {
		return nil, ErrEmptyFilter
	}

	return &filter, nil
}
