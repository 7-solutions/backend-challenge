package errors

import (
	"fmt"
)

// NotFoundError indicates that a requested resource does not exist
// (HTTP 404 Not Found).
type NotFoundError struct {
	// Resource is the type of entity that was not found (e.g. "User").
	Resource string
	// ID is the identifier that was looked up.
	ID    string
	Cause error
}

// Error implements the error interface.
func (e *NotFoundError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s with ID %s not found: %v", e.Resource, e.ID, e.Cause)
	}
	return fmt.Sprintf("%s with ID %s not found", e.Resource, e.ID)
}

// Unwrap returns the underlying cause, enabling errors.Is / errors.As chains.
func (e *NotFoundError) Unwrap() error {
	return e.Cause
}

// NewNotFoundError creates a NotFoundError without an underlying cause.
func NewNotFoundError(resource, id string) *NotFoundError {
	return &NotFoundError{
		Resource: resource,
		ID:       id,
	}
}

// WrapNotFound creates a NotFoundError that wraps an underlying cause.
// Use this in repository layers when classifying a lower-level error (e.g.
// sql.ErrNoRows) as a domain not-found condition.
func WrapNotFound(resource, id string, cause error) *NotFoundError {
	return &NotFoundError{
		Resource: resource,
		ID:       id,
		Cause:    cause,
	}
}
