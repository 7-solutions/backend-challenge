package errors_test

import (
	"errors"
	"fmt"
	"testing"

	apperrors "github.com/chaihanij/backend-challenge/app/errors"

	"github.com/stretchr/testify/assert"
)

func TestNewUnauthorized_Error(t *testing.T) {
	err := apperrors.NewUnauthorized("token expired")

	assert.EqualError(t, err, "unauthorized: token expired")
}

func TestNewUnauthorized_Fields(t *testing.T) {
	err := apperrors.NewUnauthorized("invalid token")

	assert.Equal(t, "invalid token", err.Reason)
	assert.Nil(t, err.Cause)
}

func TestWrapUnauthorized_Error(t *testing.T) {
	cause := fmt.Errorf("jwt parse failed")
	err := apperrors.WrapUnauthorized("token malformed", cause)

	assert.EqualError(t, err, "unauthorized: token malformed (jwt parse failed)")
}

func TestWrapUnauthorized_Fields(t *testing.T) {
	cause := fmt.Errorf("signature mismatch")
	err := apperrors.WrapUnauthorized("token invalid", cause)

	assert.Equal(t, "token invalid", err.Reason)
	assert.Equal(t, cause, err.Cause)
}

func TestUnauthorizedError_Unwrap(t *testing.T) {
	cause := fmt.Errorf("underlying cause")
	err := apperrors.WrapUnauthorized("token expired", cause)

	assert.Equal(t, cause, errors.Unwrap(err))
}

func TestUnauthorizedError_ErrorsIs(t *testing.T) {
	cause := fmt.Errorf("root cause")
	err := apperrors.WrapUnauthorized("token expired", cause)

	assert.True(t, errors.Is(err, cause))
}

func TestNewUnauthorized_NilCause_Unwrap(t *testing.T) {
	err := apperrors.NewUnauthorized("missing credentials")

	assert.Nil(t, errors.Unwrap(err))
}

func TestUnauthorizedError_NilCause_ErrorFormat(t *testing.T) {
	err := apperrors.NewUnauthorized("missing credentials")

	assert.NotContains(t, err.Error(), "(")
}
