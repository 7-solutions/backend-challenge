package errors_test

import (
	"errors"
	"fmt"
	"testing"

	apperrors "github.com/chaihanij/backend-challenge/app/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ── NewInternalError ──────────────────────────────────────────────────────────

func TestNewInternalError_FieldsAreSet(t *testing.T) {
	cause := errors.New("db connection refused")
	err := apperrors.NewInternalError("postgres", "query", cause)

	require.NotNil(t, err)
	assert.Equal(t, "postgres", err.Dependency)
	assert.Equal(t, "query", err.Operation)
	assert.Equal(t, cause, err.Cause)
}

func TestNewInternalError_NilCause(t *testing.T) {
	err := apperrors.NewInternalError("redis", "set", nil)

	require.NotNil(t, err)
	assert.Equal(t, "redis", err.Dependency)
	assert.Equal(t, "set", err.Operation)
	assert.Nil(t, err.Cause)
}

// ── Error() ───────────────────────────────────────────────────────────────────

func TestInternalError_ErrorMessage_Format(t *testing.T) {
	cause := errors.New("timeout")
	err := apperrors.NewInternalError("grpc", "call", cause)

	expected := fmt.Sprintf("internal failure: %s (%s): %v", "grpc", "call", cause)
	assert.EqualError(t, err, expected)
}

func TestInternalError_ErrorMessage_NilCause(t *testing.T) {
	err := apperrors.NewInternalError("s3", "upload", nil)

	assert.EqualError(t, err, "internal failure: s3 (upload): <nil>")
}

func TestInternalError_ErrorMessage_ContainsDependency(t *testing.T) {
	err := apperrors.NewInternalError("mysql", "insert", errors.New("dup key"))

	assert.Contains(t, err.Error(), "mysql")
}

func TestInternalError_ErrorMessage_ContainsOperation(t *testing.T) {
	err := apperrors.NewInternalError("mysql", "insert", errors.New("dup key"))

	assert.Contains(t, err.Error(), "insert")
}

func TestInternalError_ErrorMessage_ContainsCause(t *testing.T) {
	err := apperrors.NewInternalError("mysql", "insert", errors.New("dup key"))

	assert.Contains(t, err.Error(), "dup key")
}

// ── Unwrap() ──────────────────────────────────────────────────────────────────

func TestInternalError_Unwrap_ReturnsCause(t *testing.T) {
	cause := errors.New("original error")
	err := apperrors.NewInternalError("svc", "op", cause)

	assert.Equal(t, cause, err.Unwrap())
}

func TestInternalError_Unwrap_NilCause(t *testing.T) {
	err := apperrors.NewInternalError("svc", "op", nil)

	assert.Nil(t, err.Unwrap())
}

func TestInternalError_ErrorsIs_WrappedCause(t *testing.T) {
	sentinel := errors.New("sentinel")
	err := apperrors.NewInternalError("svc", "op", sentinel)

	assert.True(t, errors.Is(err, sentinel))
}

// ── errors.As ────────────────────────────────────────────────────────────────

func TestInternalError_ErrorsAs_DirectMatch(t *testing.T) {
	err := apperrors.NewInternalError("kafka", "publish", errors.New("broken pipe"))

	var target *apperrors.InternalError
	assert.True(t, errors.As(err, &target))
	assert.Equal(t, "kafka", target.Dependency)
	assert.Equal(t, "publish", target.Operation)
}

func TestInternalError_ErrorsAs_WrappedInAnotherError(t *testing.T) {
	inner := apperrors.NewInternalError("db", "scan", errors.New("EOF"))
	wrapped := fmt.Errorf("handler: %w", inner)

	var target *apperrors.InternalError
	assert.True(t, errors.As(wrapped, &target))
	assert.Equal(t, "db", target.Dependency)
}
