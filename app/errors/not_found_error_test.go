package errors_test

import (
	"errors"
	"fmt"
	"testing"

	apperrors "github.com/chaihanij/backend-challenge/app/errors"

	"github.com/stretchr/testify/assert"
)

func TestNewNotFoundError_Error(t *testing.T) {
	err := apperrors.NewNotFoundError("User", "42")

	assert.Equal(t, "User with ID 42 not found", err.Error())
}

func TestNewNotFoundError_Fields(t *testing.T) {
	err := apperrors.NewNotFoundError("Order", "99")

	assert.Equal(t, "Order", err.Resource)
	assert.Equal(t, "99", err.ID)
	assert.Nil(t, err.Cause)
}

func TestWrapNotFound_Error(t *testing.T) {
	cause := fmt.Errorf("db connection failed")
	err := apperrors.WrapNotFound("Product", "7", cause)

	assert.Equal(t, "Product with ID 7 not found: db connection failed", err.Error())
}

func TestWrapNotFound_Fields(t *testing.T) {
	cause := fmt.Errorf("some cause")
	err := apperrors.WrapNotFound("Item", "1", cause)

	assert.Equal(t, "Item", err.Resource)
	assert.Equal(t, "1", err.ID)
	assert.Equal(t, cause, err.Cause)
}

func TestNotFoundError_Unwrap(t *testing.T) {
	cause := fmt.Errorf("original error")
	err := apperrors.WrapNotFound("Resource", "5", cause)

	assert.Equal(t, cause, errors.Unwrap(err))
}

func TestNotFoundError_ErrorsIs(t *testing.T) {
	cause := fmt.Errorf("root cause")
	err := apperrors.WrapNotFound("Resource", "5", cause)

	// errors.Is ควรสามารถ unwrap ถึง cause ได้
	assert.True(t, errors.Is(err, cause))
}

func TestNewNotFoundError_NilCause(t *testing.T) {
	err := apperrors.NewNotFoundError("File", "abc")

	// Unwrap ควร return nil
	assert.Nil(t, errors.Unwrap(err))
}
