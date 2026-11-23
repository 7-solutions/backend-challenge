package errors_test

import (
	"errors"
	"fmt"
	"testing"

	apperrors "github.com/chaihanij/backend-challenge/app/errors"
	"github.com/stretchr/testify/assert"
)

func TestNewConflictError_Error(t *testing.T) {
	err := apperrors.NewConflictError("User", "email already exists")

	assert.EqualError(t, err, "User conflict: email already exists")
}

func TestNewConflictError_Fields(t *testing.T) {
	err := apperrors.NewConflictError("Order", "duplicate order")

	assert.Equal(t, "Order", err.Resource)
	assert.Equal(t, "duplicate order", err.Reason)
	assert.Nil(t, err.Cause)
}

func TestWrapConflict_Error(t *testing.T) {
	cause := fmt.Errorf("unique constraint violation")
	err := apperrors.WrapConflict("Product", "SKU already exists", cause)

	assert.EqualError(t, err, "Product conflict: SKU already exists (unique constraint violation)")
}

func TestWrapConflict_Fields(t *testing.T) {
	cause := fmt.Errorf("db error")
	err := apperrors.WrapConflict("Item", "duplicate key", cause)

	assert.Equal(t, "Item", err.Resource)
	assert.Equal(t, "duplicate key", err.Reason)
	assert.Equal(t, cause, err.Cause)
}

func TestConflictError_Unwrap(t *testing.T) {
	cause := fmt.Errorf("original cause")
	err := apperrors.WrapConflict("Resource", "conflict reason", cause)

	assert.Equal(t, cause, errors.Unwrap(err))
}

func TestConflictError_ErrorsIs(t *testing.T) {
	cause := fmt.Errorf("root cause")
	err := apperrors.WrapConflict("Resource", "conflict reason", cause)

	assert.True(t, errors.Is(err, cause))
}

func TestNewConflictError_NilCause_Unwrap(t *testing.T) {
	err := apperrors.NewConflictError("File", "already exists")

	assert.Nil(t, errors.Unwrap(err))
}

func TestConflictError_NilCause_ErrorFormat(t *testing.T) {
	err := apperrors.NewConflictError("User", "username taken")

	// ต้องไม่มี cause ต่อท้าย
	assert.NotContains(t, err.Error(), "(")
}
