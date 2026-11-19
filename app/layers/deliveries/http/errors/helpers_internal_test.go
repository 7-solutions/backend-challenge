package errors

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFirstErr(t *testing.T) {
	t.Run("returns nil for empty slice", func(t *testing.T) {
		assert.Nil(t, firstErr(nil))
		assert.Nil(t, firstErr([]error{}))
	})

	t.Run("returns first error when multiple provided", func(t *testing.T) {
		err1 := errors.New("first")
		err2 := errors.New("second")

		assert.Equal(t, err1, firstErr([]error{err1, err2}))
	})
}

func TestProblemType(t *testing.T) {
	BaseUrl = "https://api.example.com"

	assert.Equal(t, "https://api.example.com/problems/bad-request", problemType("bad-request"))
	assert.Equal(t, "https://api.example.com/problems/not-found", problemType("not-found"))
}
