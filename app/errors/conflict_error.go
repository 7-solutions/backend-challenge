package errors

import (
	"fmt"
)

// ConflictError indicates a state conflict for a resource (HTTP 409 Conflict),
// for example a duplicate key or an optimistic-lock violation.
type ConflictError struct {
	// Resource is the type of entity involved in the conflict.
	Resource string
	// Reason describes why the conflict occurred.
	Reason string
	Cause  error
}

// Error implements the error interface.
func (e *ConflictError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s conflict: %s (%v)", e.Resource, e.Reason, e.Cause)
	}
	return fmt.Sprintf("%s conflict: %s", e.Resource, e.Reason)
}

// Unwrap returns the underlying cause, enabling errors.Is / errors.As chains.
func (e *ConflictError) Unwrap() error {
	return e.Cause
}

// NewConflictError creates a ConflictError without an underlying cause.
func NewConflictError(resource, reason string) *ConflictError {
	return &ConflictError{Resource: resource, Reason: reason}
}

// WrapConflict creates a ConflictError that wraps an underlying cause.
func WrapConflict(resource, reason string, cause error) *ConflictError {
	return &ConflictError{Resource: resource, Reason: reason, Cause: cause}
}
