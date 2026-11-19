package errors

import "fmt"

// UnauthorizedError indicates that a request lacks valid authentication
// credentials (HTTP 401 Unauthorized).
type UnauthorizedError struct {
	// Reason describes why authentication failed (e.g. "token expired").
	Reason string
	Cause  error
}

// Error implements the error interface.
func (e *UnauthorizedError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("unauthorized: %s (%v)", e.Reason, e.Cause)
	}
	return fmt.Sprintf("unauthorized: %s", e.Reason)
}

// Unwrap returns the underlying cause, enabling errors.Is / errors.As chains.
func (e *UnauthorizedError) Unwrap() error {
	return e.Cause
}

// NewUnauthorized creates an UnauthorizedError without an underlying cause.
func NewUnauthorized(reason string) *UnauthorizedError {
	return &UnauthorizedError{Reason: reason}
}

// WrapUnauthorized creates an UnauthorizedError that wraps an underlying cause.
func WrapUnauthorized(reason string, cause error) *UnauthorizedError {
	return &UnauthorizedError{Reason: reason, Cause: cause}
}
