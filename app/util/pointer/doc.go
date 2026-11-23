// Package pointer provides generic helpers for converting between values and
// pointers using Go generics.
//
// # ToPointer — value to pointer
//
// ToPointer returns a pointer to the given value. When the value is the zero
// value for its type, it returns nil instead of a pointer to the zero value.
// This is useful for optional struct fields where nil explicitly signals "not set":
//
//	pointer.ToPointer("hello")  // → *string pointing to "hello"
//	pointer.ToPointer("")       // → nil  (zero value for string)
//	pointer.ToPointer(42)       // → *int pointing to 42
//	pointer.ToPointer(0)        // → nil  (zero value for int)
//
// Warning — bool: false is the zero value for bool, so ToPointer(false) returns
// nil. When the distinction between nil and false matters (e.g. a nullable
// boolean column), do not use ToPointer; assign the pointer directly instead.
//
// # ToValue — pointer to value
//
// ToValue dereferences a pointer. When the pointer is nil it returns the zero
// value for T rather than panicking:
//
//	s := pointer.ToValue(ptr)           // "hello"
//	z := pointer.ToValue[string](nil)   // ""
package pointer
