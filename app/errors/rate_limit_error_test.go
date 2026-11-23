package errors_test

import (
	"errors"
	"fmt"
	"testing"

	apperrors "github.com/chaihanij/backend-challenge/app/errors"

	"github.com/stretchr/testify/assert"
)

func TestNewRateLimitError_Error(t *testing.T) {
	err := apperrors.NewRateLimitError(100, 60)

	assert.EqualError(t, err, "rate limit exceeded: 100 requests per 60 seconds")
}

func TestNewRateLimitError_Fields(t *testing.T) {
	err := apperrors.NewRateLimitError(50, 30)

	assert.Equal(t, 50, err.Limit)
	assert.Equal(t, 30, err.WindowSec)
	assert.Nil(t, err.Cause)
}

func TestWrapRateLimitError_Error(t *testing.T) {
	cause := fmt.Errorf("too many requests from IP")
	err := apperrors.WrapRateLimitError(10, 1, cause)

	assert.EqualError(t, err, "rate limit exceeded: 10 requests per 1 seconds")
}

func TestWrapRateLimitError_Fields(t *testing.T) {
	cause := fmt.Errorf("upstream throttled")
	err := apperrors.WrapRateLimitError(200, 60, cause)

	assert.Equal(t, 200, err.Limit)
	assert.Equal(t, 60, err.WindowSec)
	assert.Equal(t, cause, err.Cause)
}

func TestRateLimitError_Unwrap(t *testing.T) {
	cause := fmt.Errorf("underlying cause")
	err := apperrors.WrapRateLimitError(5, 10, cause)

	assert.Equal(t, cause, errors.Unwrap(err))
}

func TestRateLimitError_ErrorsIs(t *testing.T) {
	cause := fmt.Errorf("root cause")
	err := apperrors.WrapRateLimitError(5, 10, cause)

	assert.True(t, errors.Is(err, cause))
}

func TestRateLimitError_ErrorsAs(t *testing.T) {
	cause := fmt.Errorf("throttle")
	err := apperrors.WrapRateLimitError(20, 5, cause)

	var rlErr *apperrors.RateLimitError
	assert.True(t, errors.As(err, &rlErr))
	assert.Equal(t, 20, rlErr.Limit)
	assert.Equal(t, 5, rlErr.WindowSec)
}

func TestNewRateLimitError_NilCause_Unwrap(t *testing.T) {
	err := apperrors.NewRateLimitError(100, 60)

	assert.Nil(t, errors.Unwrap(err))
}

func TestRateLimitError_ZeroValues(t *testing.T) {
	err := apperrors.NewRateLimitError(0, 0)

	assert.EqualError(t, err, "rate limit exceeded: 0 requests per 0 seconds")
}
