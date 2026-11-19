package errors

// ValidationError represents one or more field-level validation failures
// (HTTP 400 Bad Request). Fields maps each failing field name to a
// human-readable message.
type ValidationError struct {
	// Fields maps field names to human-readable validation messages.
	Fields map[string]string
	Cause  error
}

// Error implements the error interface.
func (e *ValidationError) Error() string {
	return "validation failed"
}

// Unwrap returns the underlying cause, enabling errors.Is / errors.As chains.
func (e *ValidationError) Unwrap() error {
	return e.Cause
}

// NewValidationError creates a ValidationError with the given field messages.
func NewValidationError(fields map[string]string) *ValidationError {
	return &ValidationError{Fields: fields}
}

// WrapValidationError creates a ValidationError that wraps an underlying cause.
func WrapValidationError(fields map[string]string, cause error) *ValidationError {
	return &ValidationError{Fields: fields, Cause: cause}
}
