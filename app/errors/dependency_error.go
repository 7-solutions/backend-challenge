package errors

import "fmt"

// DependencyError indicates a failure in an external dependency such as a
// database, cache, or third-party API (HTTP 503 Service Unavailable).
type DependencyError struct {
	// Dependency names the external system (e.g. "postgres", "redis").
	Dependency string
	// Operation describes what was being attempted (e.g. "query", "set").
	Operation string
	Cause     error
}

// Error implements the error interface.
func (e *DependencyError) Error() string {
	return fmt.Sprintf("dependency failure: %s (%s): %v", e.Dependency, e.Operation, e.Cause)
}

// Unwrap returns the underlying cause, enabling errors.Is / errors.As chains.
func (e *DependencyError) Unwrap() error {
	return e.Cause
}

// NewDependencyError creates a DependencyError wrapping the given cause.
func NewDependencyError(dep, op string, cause error) *DependencyError {
	return &DependencyError{Dependency: dep, Operation: op, Cause: cause}
}
