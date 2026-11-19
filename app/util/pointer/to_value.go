package pointer

// ToValue dereferences pointer. When pointer is nil it returns the zero value
// for T rather than panicking.
func ToValue[T any](pointer *T) T {
	if pointer == nil {

		var zeroValue T

		return zeroValue
	}
	return *pointer
}
