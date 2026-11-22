package models_test

import (
	"errors"
	"strings"
	"testing"

	apperrors "github.com/chaihanij/backend-challenge/app/errors"
	"github.com/chaihanij/backend-challenge/app/layers/deliveries/http/api/users/models"
	testhttp "github.com/chaihanij/backend-challenge/app/test_helpers/http"
)

func TestNewUpdateUserRequest(t *testing.T) {
	t.Run("returns request when payload is valid", func(t *testing.T) {
		actual, err := models.NewUpdateUserRequest(testhttp.GinContextJSON(`{"email":"user@example.com","name":"Example User"}`))
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if actual.Email != "user@example.com" || actual.Name != "Example User" {
			t.Fatalf("unexpected request: %+v", actual)
		}
	})

	tests := []struct {
		name string
		body string
	}{
		{name: "malformed JSON", body: `{"email":`},
		{name: "missing email", body: `{"name":"Example User"}`},
		{name: "invalid email", body: `{"email":"invalid","name":"Example User"}`},
		{name: "email too short", body: `{"email":"a","name":"Example User"}`},
		{name: "email too long", body: `{"email":"` + strings.Repeat("a", 250) + `@example.com","name":"Example User"}`},
		{name: "missing name", body: `{"email":"user@example.com"}`},
		{name: "name too short", body: `{"email":"user@example.com","name":"Jo"}`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual, err := models.NewUpdateUserRequest(testhttp.GinContextJSON(tt.body))
			if actual != nil {
				t.Fatalf("expected nil request, got %+v", actual)
			}
			if err == nil {
				t.Fatal("expected error")
			}
			if _, ok := errors.AsType[*apperrors.ValidationError](err); !ok {
				t.Fatalf("expected ValidationError, got %T: %v", err, err)
			}
		})
	}
}
