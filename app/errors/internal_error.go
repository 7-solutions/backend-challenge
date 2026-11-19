package errors

import (
	"fmt"
)

// InternalError indicates an unexpected internal failure in the application
// itself (HTTP 500 Internal Server Error), as opposed to a failure in an
// external dependency.
type InternalError struct {
	// Dependency names the internal component where the failure occurred.
	Dependency string
	// Operation describes what was being attempted.
	Operation string
	Cause     error
}

// Error implements the error interface.
func (e *InternalError) Error() string {
	return fmt.Sprintf("internal failure: %s (%s): %v", e.Dependency, e.Operation, e.Cause)
}

// Unwrap returns the underlying cause, enabling errors.Is / errors.As chains.
func (e *InternalError) Unwrap() error {
	return e.Cause
}

// NewInternalError creates an InternalError wrapping the given cause.
func NewInternalError(dep, op string, cause error) *InternalError {
	return &InternalError{Dependency: dep, Operation: op, Cause: cause}
}
