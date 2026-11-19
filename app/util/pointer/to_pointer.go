package pointer

import (
	"reflect"
)

// ToPointer returns a pointer to value. When value is the zero value for T,
// it returns nil so the caller can distinguish "not provided" from "set to
// zero". See the package documentation for the bool-zero-value warning.
func ToPointer[T any](value T) *T {
	if reflect.ValueOf(value).IsZero() {
		return nil
	}
	return &value
}
