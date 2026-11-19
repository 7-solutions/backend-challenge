package errors_test

import (
	"errors"
	"fmt"
	"testing"

	apperrors "github.com/chaihanij/backend-challenge/app/errors"

	"github.com/stretchr/testify/assert"
)

func TestNewDependencyError_Error(t *testing.T) {
	cause := fmt.Errorf("connection refused")
	err := apperrors.NewDependencyError("database", "connect", cause)

	assert.EqualError(t, err, "dependency failure: database (connect): connection refused")
}

func TestNewDependencyError_Fields(t *testing.T) {
	cause := fmt.Errorf("timeout")
	err := apperrors.NewDependencyError("redis", "get", cause)

	assert.Equal(t, "redis", err.Dependency)
	assert.Equal(t, "get", err.Operation)
	assert.Equal(t, cause, err.Cause)
}

func TestDependencyError_Unwrap(t *testing.T) {
	cause := fmt.Errorf("network error")
	err := apperrors.NewDependencyError("external-api", "call", cause)

	assert.Equal(t, cause, errors.Unwrap(err))
}

func TestDependencyError_ErrorsIs(t *testing.T) {
	cause := fmt.Errorf("root cause")
	err := apperrors.NewDependencyError("kafka", "publish", cause)

	assert.True(t, errors.Is(err, cause))
}

func TestDependencyError_ErrorsAs(t *testing.T) {
	cause := fmt.Errorf("timeout")
	err := apperrors.NewDependencyError("s3", "upload", cause)

	var depErr *apperrors.DependencyError
	assert.True(t, errors.As(err, &depErr))
	assert.Equal(t, "s3", depErr.Dependency)
	assert.Equal(t, "upload", depErr.Operation)
}

func TestDependencyError_NilCause(t *testing.T) {
	err := apperrors.NewDependencyError("database", "query", nil)

	assert.Nil(t, errors.Unwrap(err))
	assert.Contains(t, err.Error(), "dependency failure: database (query)")
}
