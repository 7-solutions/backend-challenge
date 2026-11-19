package errors

import "fmt"

// PermissionDeniedError indicates that an authenticated caller lacks the
// required permission to perform an action (HTTP 403 Forbidden).
type PermissionDeniedError struct {
	// Action describes the attempted operation (e.g. "delete").
	Action string
	// Resource is the entity the action was attempted on (e.g. "User").
	Resource string
	Cause    error
}

// Error implements the error interface.
func (e *PermissionDeniedError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("permission denied: cannot %s %s (%v)", e.Action, e.Resource, e.Cause)
	}
	return fmt.Sprintf("permission denied: cannot %s %s", e.Action, e.Resource)
}

// Unwrap returns the underlying cause, enabling errors.Is / errors.As chains.
func (e *PermissionDeniedError) Unwrap() error {
	return e.Cause
}

// NewPermissionDenied creates a PermissionDeniedError without an underlying cause.
func NewPermissionDenied(action, resource string) *PermissionDeniedError {
	return &PermissionDeniedError{Action: action, Resource: resource}
}

// WrapPermissionDenied creates a PermissionDeniedError that wraps an underlying cause.
func WrapPermissionDenied(action, resource string, cause error) *PermissionDeniedError {
	return &PermissionDeniedError{Action: action, Resource: resource, Cause: cause}
}
