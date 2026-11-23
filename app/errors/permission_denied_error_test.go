package errors_test

import (
	"errors"
	"fmt"
	"testing"

	apperrors "github.com/chaihanij/backend-challenge/app/errors"

	"github.com/stretchr/testify/assert"
)

func TestNewPermissionDenied_Error(t *testing.T) {
	err := apperrors.NewPermissionDenied("delete", "Invoice")

	assert.EqualError(t, err, "permission denied: cannot delete Invoice")
}

func TestNewPermissionDenied_Fields(t *testing.T) {
	err := apperrors.NewPermissionDenied("update", "User")

	assert.Equal(t, "update", err.Action)
	assert.Equal(t, "User", err.Resource)
	assert.Nil(t, err.Cause)
}

func TestWrapPermissionDenied_Error(t *testing.T) {
	cause := fmt.Errorf("role check failed")
	err := apperrors.WrapPermissionDenied("read", "Report", cause)

	assert.EqualError(t, err, "permission denied: cannot read Report (role check failed)")
}

func TestWrapPermissionDenied_Fields(t *testing.T) {
	cause := fmt.Errorf("auth error")
	err := apperrors.WrapPermissionDenied("export", "Document", cause)

	assert.Equal(t, "export", err.Action)
	assert.Equal(t, "Document", err.Resource)
	assert.Equal(t, cause, err.Cause)
}

func TestPermissionDeniedError_Unwrap(t *testing.T) {
	cause := fmt.Errorf("underlying cause")
	err := apperrors.WrapPermissionDenied("write", "Config", cause)

	assert.Equal(t, cause, errors.Unwrap(err))
}

func TestPermissionDeniedError_ErrorsIs(t *testing.T) {
	cause := fmt.Errorf("root cause")
	err := apperrors.WrapPermissionDenied("delete", "Record", cause)

	assert.True(t, errors.Is(err, cause))
}

func TestNewPermissionDenied_NilCause_Unwrap(t *testing.T) {
	err := apperrors.NewPermissionDenied("create", "Project")

	assert.Nil(t, errors.Unwrap(err))
}

func TestPermissionDeniedError_NilCause_ErrorFormat(t *testing.T) {
	err := apperrors.NewPermissionDenied("access", "Dashboard")

	// ต้องไม่มี cause ต่อท้าย
	assert.NotContains(t, err.Error(), "(")
}
