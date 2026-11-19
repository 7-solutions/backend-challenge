// Package errors provides structured domain error types that carry HTTP status
// semantics without importing net/http. The library's httpx/errors package
// maps these types to RFC 7807 ProblemDetails responses.
//
// # Error types and HTTP status mapping
//
//	ValidationError        → 400 Bad Request
//	NotFoundError          → 404 Not Found
//	ConflictError          → 409 Conflict
//	PermissionDeniedError  → 403 Forbidden
//	RateLimitError         → 429 Too Many Requests
//	DependencyError        → 503 Service Unavailable
//	InternalError          → 500 Internal Server Error
//
// # Two-factory pattern
//
// Each error type offers two constructors:
//
//   - NewXxx — creates a standalone error with no wrapped cause. Use this at
//     the edge of your domain when the condition itself is the error.
//
//   - WrapXxx — wraps an underlying cause. Use this in repository or service
//     layers when you want to classify a lower-level error (e.g. a database
//     error) while preserving the original for errors.Is / errors.As chains.
//
// Example:
//
//	// service layer
//	func GetUser(id string) (*User, error) {
//	    u, err := repo.FindUser(id)
//	    if errors.Is(err, sql.ErrNoRows) {
//	        return nil, errors.WrapNotFound("User", id, err)
//	    }
//	    return u, err
//	}
//
//	// handler
//	var nfe *errors.NotFoundError
//	if errors.As(err, &nfe) {
//	    // nfe.Resource == "User", nfe.ID == "42"
//	}
//
// # Sentinel errors
//
// Four package-level sentinel values cover simple cases that do not need
// additional context:
//
//	errors.ErrNilInput     // input must not be nil
//	errors.ErrInvalidInput // input is malformed
//	errors.ErrUnauthorized // missing or invalid credentials
//	errors.ErrForbidden    // authenticated but not authorised
//
// # Wrap
//
// Wrap is a lightweight alternative to fmt.Errorf("%s: %w", msg, err) that
// returns nil when err is nil, avoiding defensive nil checks at call sites:
//
//	return errors.Wrap(repo.Save(u), "CreateUser")
package errors
