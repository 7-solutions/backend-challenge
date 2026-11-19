package models_test

import (
	"errors"
	"testing"

	"github.com/chaihanij/backend-challenge/app/entities"
	"github.com/chaihanij/backend-challenge/app/layers/repositories/users/models"
	"github.com/chaihanij/backend-challenge/app/util/pointer"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func TestNewUserFilter(t *testing.T) {
	t.Run("error when input is nil", func(t *testing.T) {
		filter, err := models.NewUserFilter(nil)
		if !errors.Is(err, models.ErrInvalidFilterInput) {
			t.Fatalf("expected ErrInvalidFilterInput, got: %v", err)
		}
		if filter != nil {
			t.Fatalf("expected nil filter, got: %#v", filter)
		}
	})

	t.Run("error when input type is unsupported", func(t *testing.T) {
		filter, err := models.NewUserFilter(123)
		if !errors.Is(err, models.ErrInvalidFilterInput) {
			t.Fatalf("expected ErrInvalidFilterInput, got: %v", err)
		}
		if filter != nil {
			t.Fatalf("expected nil filter, got: %#v", filter)
		}
	})

	t.Run("builds _id filter from valid id hex string", func(t *testing.T) {
		id := bson.NewObjectID()

		filter, err := models.NewUserFilter(id.Hex())
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if filter == nil {
			t.Fatal("expected non-nil filter")
		}

		gotID, ok := (*filter)["_id"].(bson.ObjectID)
		if !ok {
			t.Fatalf("expected _id to be bson.ObjectID, got: %#v", (*filter)["_id"])
		}
		if gotID != id {
			t.Fatalf("expected _id %s, got %s", id.Hex(), gotID.Hex())
		}
	})

	t.Run("error for invalid id hex string", func(t *testing.T) {
		filter, err := models.NewUserFilter("not-an-object-id")
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if filter != nil {
			t.Fatalf("expected nil filter, got: %#v", filter)
		}
	})

	t.Run("builds filter from user with id email and name", func(t *testing.T) {
		id := bson.NewObjectID()
		input := &entities.User{
			ID:    id.Hex(),
			Email: "tester@example.com",
			Name:  "Tester",
		}

		filter, err := models.NewUserFilter(input)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if filter == nil {
			t.Fatal("expected non-nil filter")
		}
		if len(*filter) != 3 {
			t.Fatalf("expected 3 fields, got %d", len(*filter))
		}

		gotID, ok := (*filter)["_id"].(bson.ObjectID)
		if !ok {
			t.Fatalf("expected _id to be bson.ObjectID, got: %#v", (*filter)["_id"])
		}
		if gotID != id {
			t.Fatalf("expected _id %s, got %s", id.Hex(), gotID.Hex())
		}
		if (*filter)["email"] != input.Email {
			t.Fatalf("expected email %q, got %#v", input.Email, (*filter)["email"])
		}
		if (*filter)["name"] != input.Name {
			t.Fatalf("expected name %q, got %#v", input.Name, (*filter)["name"])
		}
	})

	t.Run("error from user with invalid id", func(t *testing.T) {
		input := &entities.User{ID: "invalid-id", Email: "x@example.com"}

		filter, err := models.NewUserFilter(input)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if filter != nil {
			t.Fatalf("expected nil filter, got: %#v", filter)
		}
	})

	t.Run("error when user has no filterable fields", func(t *testing.T) {
		filter, err := models.NewUserFilter(&entities.User{})
		if !errors.Is(err, models.ErrEmptyFilter) {
			t.Fatalf("expected ErrEmptyFilter, got: %v", err)
		}
		if filter != nil {
			t.Fatalf("expected nil filter, got: %#v", filter)
		}
	})

	t.Run("builds filter from user filter pointers", func(t *testing.T) {
		id := bson.NewObjectID().Hex()
		email := "pointer@example.com"
		name := "Pointer"

		input := &entities.UserFilter{
			ID:    pointer.ToPointer(id),
			Email: pointer.ToPointer(email),
			Name:  pointer.ToPointer(name),
		}

		filter, err := models.NewUserFilter(input)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if filter == nil {
			t.Fatal("expected non-nil filter")
		}
		if len(*filter) != 3 {
			t.Fatalf("expected 3 fields, got %d", len(*filter))
		}

		gotID, ok := (*filter)["_id"].(bson.ObjectID)
		if !ok {
			t.Fatalf("expected _id to be bson.ObjectID, got: %#v", (*filter)["_id"])
		}
		if gotID.Hex() != id {
			t.Fatalf("expected _id %s, got %s", id, gotID.Hex())
		}
		if (*filter)["email"] != email {
			t.Fatalf("expected email %q, got %#v", email, (*filter)["email"])
		}
		if (*filter)["name"] != name {
			t.Fatalf("expected name %q, got %#v", name, (*filter)["name"])
		}
	})

	t.Run("error when user filter has no filterable fields", func(t *testing.T) {
		filter, err := models.NewUserFilter(&entities.UserFilter{})
		if !errors.Is(err, models.ErrEmptyFilter) {
			t.Fatalf("expected ErrEmptyFilter, got: %v", err)
		}
		if filter != nil {
			t.Fatalf("expected nil filter, got: %#v", filter)
		}
	})

	t.Run("error when typed nil user pointer is passed", func(t *testing.T) {
		var user *entities.User

		filter, err := models.NewUserFilter(user)
		if !errors.Is(err, models.ErrInvalidFilterInput) {
			t.Fatalf("expected ErrInvalidFilterInput, got: %v", err)
		}
		if filter != nil {
			t.Fatalf("expected nil filter, got: %#v", filter)
		}
	})
}
