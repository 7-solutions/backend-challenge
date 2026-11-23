package models_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/chaihanij/backend-challenge/app/entities"
	"github.com/chaihanij/backend-challenge/app/layers/repositories/users/models"
	"github.com/chaihanij/backend-challenge/app/util/password"
)

func TestNewUserModel(t *testing.T) {
	t.Run("error when input is nil", func(t *testing.T) {
		model, err := models.NewUserModel(nil)
		if !errors.Is(err, models.ErrNilUser) {
			t.Fatalf("expected ErrNilUser, got: %v", err)
		}
		if model != nil {
			t.Fatalf("expected nil model, got: %#v", model)
		}
	})

	t.Run("builds model and hashes password", func(t *testing.T) {
		input := &entities.User{
			Name:     "Tester",
			Email:    "tester@example.com",
			Password: "P@ssword123",
		}

		model, err := models.NewUserModel(input)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if model == nil {
			t.Fatal("expected non-nil model")
		}

		if model.ID.IsZero() {
			t.Fatal("expected generated ID, got zero ObjectID")
		}
		if model.Name != input.Name {
			t.Fatalf("expected name %q, got %q", input.Name, model.Name)
		}
		if model.Email != input.Email {
			t.Fatalf("expected email %q, got %q", input.Email, model.Email)
		}
		if model.Password == input.Password {
			t.Fatal("expected password to be hashed, got plaintext")
		}
		if !password.CheckPasswordHash(input.Password, model.Password) {
			t.Fatal("expected hashed password to match input plaintext")
		}

		if model.CreatedAt == nil {
			t.Fatal("expected CreatedAt to be set")
		}
		if model.UpdatedAt == nil {
			t.Fatal("expected UpdatedAt to be set")
		}

		expected := model.ID.Timestamp()
		if !model.CreatedAt.Equal(expected) {
			t.Fatalf("expected CreatedAt %v, got %v", expected, *model.CreatedAt)
		}
		if !model.UpdatedAt.Equal(expected) {
			t.Fatalf("expected UpdatedAt %v, got %v", expected, *model.UpdatedAt)
		}
	})

	t.Run("returns bcrypt error when password is too long", func(t *testing.T) {
		input := &entities.User{
			Name:     "Tester",
			Email:    "tester@example.com",
			Password: strings.Repeat("a", 100),
		}

		model, err := models.NewUserModel(input)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if model != nil {
			t.Fatalf("expected nil model, got: %#v", model)
		}
	})
}
