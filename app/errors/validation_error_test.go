package errors_test

import (
	"errors"
	"fmt"
	"testing"

	apperrors "github.com/chaihanij/backend-challenge/app/errors"

	"github.com/stretchr/testify/assert"
)

func TestNewValidationError_Error(t *testing.T) {
	fields := map[string]string{"email": "invalid format"}
	err := apperrors.NewValidationError(fields)

	assert.EqualError(t, err, "validation failed")
}

func TestNewValidationError_Fields(t *testing.T) {
	fields := map[string]string{
		"email":    "invalid format",
		"username": "must not be empty",
	}
	err := apperrors.NewValidationError(fields)

	assert.Equal(t, fields, err.Fields)
	assert.Nil(t, err.Cause)
}

func TestWrapValidationError_Error(t *testing.T) {
	fields := map[string]string{"age": "must be positive"}
	cause := fmt.Errorf("parse error")
	err := apperrors.WrapValidationError(fields, cause)

	assert.EqualError(t, err, "validation failed")
}

func TestWrapValidationError_Fields(t *testing.T) {
	fields := map[string]string{"name": "too long"}
	cause := fmt.Errorf("binding error")
	err := apperrors.WrapValidationError(fields, cause)

	assert.Equal(t, fields, err.Fields)
	assert.Equal(t, cause, err.Cause)
}

func TestValidationError_Unwrap(t *testing.T) {
	cause := fmt.Errorf("underlying parse error")
	err := apperrors.WrapValidationError(map[string]string{}, cause)

	assert.Equal(t, cause, errors.Unwrap(err))
}

func TestValidationError_ErrorsIs(t *testing.T) {
	cause := fmt.Errorf("root cause")
	err := apperrors.WrapValidationError(map[string]string{"field": "invalid"}, cause)

	assert.True(t, errors.Is(err, cause))
}

func TestNewValidationError_NilCause_Unwrap(t *testing.T) {
	err := apperrors.NewValidationError(map[string]string{"x": "required"})

	assert.Nil(t, errors.Unwrap(err))
}

func TestValidationError_ErrorsAs(t *testing.T) {
	fields := map[string]string{"phone": "invalid"}
	err := apperrors.NewValidationError(fields)

	var valErr *apperrors.ValidationError
	assert.True(t, errors.As(err, &valErr))
	assert.Equal(t, fields, valErr.Fields)
}

func TestNewValidationError_EmptyFields(t *testing.T) {
	err := apperrors.NewValidationError(map[string]string{})

	assert.EqualError(t, err, "validation failed")
	assert.Empty(t, err.Fields)
}
