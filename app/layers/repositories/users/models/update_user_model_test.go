package models_test

import (
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/chaihanij/backend-challenge/app/entities"
	"github.com/chaihanij/backend-challenge/app/layers/repositories/users/models"
	"github.com/chaihanij/backend-challenge/app/util/password"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func TestUpdateUserModel(t *testing.T) {
	t.Run("error when input is nil", func(t *testing.T) {
		update, err := models.UpdateUserModel(nil)
		if !errors.Is(err, models.ErrNilUser) {
			t.Fatalf("expected ErrNilUser, got: %v", err)
		}
		if update != nil {
			t.Fatalf("expected nil update, got: %#v", update)
		}
	})

	t.Run("always includes updated_at even when no fields are provided", func(t *testing.T) {
		before := time.Now()
		update, err := models.UpdateUserModel(&entities.User{})
		after := time.Now()

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if update == nil {
			t.Fatal("expected non-nil update")
		}
		if len(*update) != 1 {
			t.Fatalf("expected only updated_at field, got %d fields", len(*update))
		}

		entry := (*update)[0]
		if entry.Key != "updated_at" {
			t.Fatalf("expected first key updated_at, got %q", entry.Key)
		}

		timestamp, ok := entry.Value.(time.Time)
		if !ok {
			t.Fatalf("expected updated_at to be time.Time, got: %#v", entry.Value)
		}
		if timestamp.Before(before) || timestamp.After(after) {
			t.Fatalf("expected updated_at between %v and %v, got %v", before, after, timestamp)
		}
	})

	t.Run("includes non-empty fields and hashes password", func(t *testing.T) {
		input := &entities.User{
			Email:    "updated@example.com",
			Name:     "Updated Name",
			Password: "N3wPassword",
		}

		update, err := models.UpdateUserModel(input)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if update == nil {
			t.Fatal("expected non-nil update")
		}

		if (*update)[0].Key != "updated_at" {
			t.Fatalf("expected first key updated_at, got %q", (*update)[0].Key)
		}

		email, ok := findBsonDValue(*update, "email")
		if !ok || email != input.Email {
			t.Fatalf("expected email %q, got %#v", input.Email, email)
		}

		name, ok := findBsonDValue(*update, "name")
		if !ok || name != input.Name {
			t.Fatalf("expected name %q, got %#v", input.Name, name)
		}

		hashed, ok := findBsonDValue(*update, "password")
		if !ok {
			t.Fatal("expected password field to be present")
		}

		hashedPassword, ok := hashed.(string)
		if !ok {
			t.Fatalf("expected password value to be string, got %#v", hashed)
		}
		if hashedPassword == input.Password {
			t.Fatal("expected password to be hashed, got plaintext")
		}
		if !password.CheckPasswordHash(input.Password, hashedPassword) {
			t.Fatal("expected hashed password to match input plaintext")
		}
	})

	t.Run("returns bcrypt error when password is too long", func(t *testing.T) {
		input := &entities.User{Password: strings.Repeat("a", 100)}

		update, err := models.UpdateUserModel(input)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if update != nil {
			t.Fatalf("expected nil update, got: %#v", update)
		}
	})
}

func findBsonDValue(input bson.D, key string) (any, bool) {
	for _, element := range input {
		if element.Key == key {
			return element.Value, true
		}
	}

	return nil, false
}
