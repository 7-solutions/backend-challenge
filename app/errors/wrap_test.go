package errors_test

import (
	"errors"
	"fmt"
	"testing"

	apperrors "github.com/chaihanij/backend-challenge/app/errors"

	"github.com/stretchr/testify/assert"
)

func TestWrap_NilError_ReturnsNil(t *testing.T) {
	result := apperrors.Wrap(nil, "something went wrong")

	assert.Nil(t, result)
}

func TestWrap_WithError_ReturnsWrappedError(t *testing.T) {
	original := fmt.Errorf("original error")
	result := apperrors.Wrap(original, "context message")

	assert.NotNil(t, result)
	assert.EqualError(t, result, "context message: original error")
}

func TestWrap_ErrorMessage_Format(t *testing.T) {
	original := fmt.Errorf("db timeout")
	result := apperrors.Wrap(original, "query failed")

	assert.EqualError(t, result, "query failed: db timeout")
}

func TestWrap_ErrorsIs_CanUnwrapToCause(t *testing.T) {
	original := fmt.Errorf("root cause")
	result := apperrors.Wrap(original, "wrapper")

	assert.True(t, errors.Is(result, original))
}

func TestWrap_MultipleWraps_CanUnwrapChain(t *testing.T) {
	root := fmt.Errorf("root")
	first := apperrors.Wrap(root, "first")
	second := apperrors.Wrap(first, "second")

	assert.EqualError(t, second, "second: first: root")
	assert.True(t, errors.Is(second, root))
}

func TestWrap_NilError_WithEmptyMessage_ReturnsNil(t *testing.T) {
	result := apperrors.Wrap(nil, "")

	assert.Nil(t, result)
}
