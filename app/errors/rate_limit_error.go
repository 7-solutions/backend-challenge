package errors

import "fmt"

// RateLimitError indicates that a caller has exceeded an allowed request rate
// (HTTP 429 Too Many Requests).
type RateLimitError struct {
	// Limit is the maximum number of allowed requests within WindowSec.
	Limit int
	// WindowSec is the sliding window duration in seconds.
	WindowSec int
	Cause     error
}

// Error implements the error interface.
func (e *RateLimitError) Error() string {
	return fmt.Sprintf("rate limit exceeded: %d requests per %d seconds", e.Limit, e.WindowSec)
}

// Unwrap returns the underlying cause, enabling errors.Is / errors.As chains.
func (e *RateLimitError) Unwrap() error {
	return e.Cause
}

// NewRateLimitError creates a RateLimitError without an underlying cause.
func NewRateLimitError(limit, windowSec int) *RateLimitError {
	return &RateLimitError{Limit: limit, WindowSec: windowSec}
}

// WrapRateLimitError creates a RateLimitError that wraps an underlying cause.
func WrapRateLimitError(limit, windowSec int, cause error) *RateLimitError {
	return &RateLimitError{Limit: limit, WindowSec: windowSec, Cause: cause}
}
