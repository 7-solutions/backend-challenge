package models_test

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	apperrors "github.com/chaihanij/backend-challenge/app/errors"
	"github.com/chaihanij/backend-challenge/app/layers/deliveries/http/api/auth/models"
	testhelpers "github.com/chaihanij/backend-challenge/app/test_helpers"
	"github.com/gin-gonic/gin"
)

func TestNewLoginRequest(t *testing.T) {
	testhelpers.Init()
	gin.SetMode(gin.TestMode)

	t.Run("returns request when payload is valid", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		email := "email@test.com"
		password := "P@ssw0rd"
		body := []byte(`{"email":"email@test.com","password":"P@ssw0rd"}`)

		c.Request = httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewBuffer(body))

		actual, err := models.NewLoginRequest(c)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if actual.Email != email {
			t.Fatalf("expected email %s, got %s", email, actual.Email)
		}

		if actual.Password != password {
			t.Fatalf("expected password %s, got %s", password, actual.Password)
		}
	})

	t.Run("returns validation error when JSON body is invalid", func(t *testing.T) {
		w := httptest.NewRecorder()

		c, _ := gin.CreateTestContext(w)
		body := []byte(``)

		c.Request = httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewBuffer(body))

		_, err := models.NewLoginRequest(c)

		if err == nil {
			t.Fatal("expected error")
		}

		if _, ok := errors.AsType[*apperrors.ValidationError](err); !ok {
			t.Fatalf("expected ValidationError, got %T: %v", err, err)
		}
	})

	t.Run("returns validation error when email is missing", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		body := []byte(`{"password":"P@ssw0rd"}`)

		c.Request = httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewBuffer(body))

		_, err := models.NewLoginRequest(c)

		if err == nil {
			t.Fatal("expected error")
		}

		if _, ok := errors.AsType[*apperrors.ValidationError](err); !ok {
			t.Fatalf("expected ValidationError, got %T: %v", err, err)
		}

	})

	t.Run("returns validation error when password is missing", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		body := []byte(`{"email":"email@test.com"}`)

		c.Request = httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewBuffer(body))

		_, err := models.NewLoginRequest(c)

		if err == nil {
			t.Fatal("expected error")
		}

		if _, ok := errors.AsType[*apperrors.ValidationError](err); !ok {
			t.Fatalf("expected ValidationError, got %T: %v", err, err)
		}
	})

	t.Run("returns validation error when email format is invalid", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		body := []byte(`{"email":"invalid-email","password":"P@ssw0rd"}`)

		c.Request = httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewBuffer(body))

		_, err := models.NewLoginRequest(c)

		if err == nil {
			t.Fatal("expected error")
		}

		if _, ok := errors.AsType[*apperrors.ValidationError](err); !ok {
			t.Fatalf("expected ValidationError, got %T: %v", err, err)
		}
	})

	t.Run("returns validation error when password is too short", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		body := []byte(`{"email":"email@test.com","password":"12345"}`)

		c.Request = httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewBuffer(body))

		_, err := models.NewLoginRequest(c)

		if err == nil {
			t.Fatal("expected error")
		}

		if _, ok := errors.AsType[*apperrors.ValidationError](err); !ok {
			t.Fatalf("expected ValidationError, got %T: %v", err, err)
		}

	})
}
