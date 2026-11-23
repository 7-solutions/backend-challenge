package errors

import (
	"fmt"
)

// Wrap annotates err with a contextual message and returns the result.
// It preserves the error chain so that errors.Is and errors.As still work.
// When err is nil, Wrap returns nil, eliminating the need for a nil guard at
// the call site.
func Wrap(err error, msg string) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("%s: %w", msg, err)
}
