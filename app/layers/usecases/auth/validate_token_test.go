package auth_test

import (
	"context"
	"errors"
	"testing"
	"time"

	customerrors "github.com/chaihanij/backend-challenge/app/errors"
	"github.com/chaihanij/backend-challenge/app/layers/usecases/auth"
	_mockuserrepo "github.com/chaihanij/backend-challenge/app/mocks/users"
	testhelpers "github.com/chaihanij/backend-challenge/app/test_helpers"
	"github.com/chaihanij/backend-challenge/app/util/pointer"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/mock/gomock"
)

func TestUseCase_ValidateToken(t *testing.T) {
	testhelpers.Init()
	ctx := context.Background()

	newValidator := func(t *testing.T) interface {
		ValidateToken(ctx context.Context, accessToken string) (*string, bool, error)
	} {
		t.Helper()
		ctrl := gomock.NewController(t)
		t.Cleanup(ctrl.Finish)

		mockUserRepo := _mockuserrepo.NewMockRepo(ctrl)
		u := auth.New(mockUserRepo, testhelpers.Logger, "secret", "iss", "aud", time.Hour)
		validator, ok := u.(interface {
			ValidateToken(ctx context.Context, accessToken string) (*string, bool, error)
		})
		if !ok {
			t.Fatalf("expected usecase to implement ValidateToken")
		}
		return validator
	}

	t.Run("Returns Subject When Token Is Valid", func(t *testing.T) {
		now := time.Now()
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"iss": "iss",
			"aud": "aud",
			"sub": "user-123",
			"exp": now.Add(time.Hour).Unix(),
			"iat": now.Unix(),
			"nbf": now.Unix(),
		})
		tokenString, err := token.SignedString([]byte("secret"))
		if err != nil {
			t.Fatalf("failed to sign token: %v", err)
		}

		validator := newValidator(t)
		sub, ok, err := validator.ValidateToken(ctx, tokenString)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if !ok {
			t.Errorf("expected ok true, got false")
		}
		if sub == nil {
			t.Errorf("expected non-nil subject, got nil")
		}
		if sub != nil && pointer.ToValue(sub) != "user-123" {
			t.Errorf("expected subject user-123, got %s", pointer.ToValue(sub))
		}
	})

	t.Run("Returns Error When Token Format Is Invalid", func(t *testing.T) {
		validator := newValidator(t)
		sub, ok, err := validator.ValidateToken(ctx, "not-a-jwt")
		if err == nil {
			t.Errorf("expected error, got nil")
		}
		if ok {
			t.Errorf("expected ok false, got true")
		}
		if sub != nil {
			t.Errorf("expected nil subject, got %v", pointer.ToValue(sub))
		}
	})

	t.Run("Returns Error When Signing Method Is Not HMAC", func(t *testing.T) {
		now := time.Now()
		token := jwt.NewWithClaims(jwt.SigningMethodNone, jwt.MapClaims{
			"iss": "iss",
			"aud": "aud",
			"sub": "user-123",
			"exp": now.Add(time.Hour).Unix(),
			"iat": now.Unix(),
			"nbf": now.Unix(),
		})
		tokenString, err := token.SignedString(jwt.UnsafeAllowNoneSignatureType)
		if err != nil {
			t.Fatalf("failed to sign token: %v", err)
		}

		validator := newValidator(t)
		sub, ok, err := validator.ValidateToken(ctx, tokenString)
		if err == nil {
			t.Errorf("expected error, got nil")
		}
		if ok {
			t.Errorf("expected ok false, got true")
		}
		if sub != nil {
			t.Errorf("expected nil subject, got %v", pointer.ToValue(sub))
		}
	})

	t.Run("Returns Unauthorized Error When Issuer Does Not Match", func(t *testing.T) {
		now := time.Now()
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"iss": "other-iss",
			"aud": "aud",
			"sub": "user-123",
			"exp": now.Add(time.Hour).Unix(),
			"iat": now.Unix(),
			"nbf": now.Unix(),
		})
		tokenString, err := token.SignedString([]byte("secret"))
		if err != nil {
			t.Fatalf("failed to sign token: %v", err)
		}

		validator := newValidator(t)
		sub, ok, err := validator.ValidateToken(ctx, tokenString)
		if err == nil {
			t.Errorf("expected error, got nil")
		}
		if ok {
			t.Errorf("expected ok false, got true")
		}
		if sub != nil {
			t.Errorf("expected nil subject, got %v", pointer.ToValue(sub))
		}
		if _, castOk := errors.AsType[*customerrors.UnauthorizedError](err); !castOk {
			t.Errorf("expected unauthorized error, got %v", err)
		}
	})

	t.Run("Returns Unauthorized Error When Audience Does Not Match", func(t *testing.T) {
		now := time.Now()
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"iss": "iss",
			"aud": "other-aud",
			"sub": "user-123",
			"exp": now.Add(time.Hour).Unix(),
			"iat": now.Unix(),
			"nbf": now.Unix(),
		})
		tokenString, err := token.SignedString([]byte("secret"))
		if err != nil {
			t.Fatalf("failed to sign token: %v", err)
		}

		validator := newValidator(t)
		sub, ok, err := validator.ValidateToken(ctx, tokenString)
		if err == nil {
			t.Errorf("expected error, got nil")
		}
		if ok {
			t.Errorf("expected ok false, got true")
		}
		if sub != nil {
			t.Errorf("expected nil subject, got %v", pointer.ToValue(sub))
		}
		if _, castOk := errors.AsType[*customerrors.UnauthorizedError](err); !castOk {
			t.Errorf("expected unauthorized error, got %v", err)
		}
	})

	t.Run("Returns Unauthorized Error When Subject Claim Is Not String", func(t *testing.T) {
		now := time.Now()
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"iss": "iss",
			"aud": "aud",
			"sub": 123,
			"exp": now.Add(time.Hour).Unix(),
			"iat": now.Unix(),
			"nbf": now.Unix(),
		})
		tokenString, err := token.SignedString([]byte("secret"))
		if err != nil {
			t.Fatalf("failed to sign token: %v", err)
		}

		validator := newValidator(t)
		sub, ok, err := validator.ValidateToken(ctx, tokenString)
		if err == nil {
			t.Errorf("expected error, got nil")
		}
		if ok {
			t.Errorf("expected ok false, got true")
		}
		if sub != nil {
			t.Errorf("expected nil subject, got %v", pointer.ToValue(sub))
		}
		if _, castOk := errors.AsType[*customerrors.UnauthorizedError](err); !castOk {
			t.Errorf("expected unauthorized error, got %v", err)
		}
	})
}
