// Package errors implements RFC 7807 Problem Details for HTTP APIs as a
// Gin-aware interface that is both an error value and a response writer.
//
// # ProblemDetails interface
//
// ProblemDetails extends the standard error interface with HTTP-level methods.
// The dual role lets you pass a ProblemDetails through the error chain (e.g.
// via c.Error) while also being able to write it directly to a Gin response:
//
//	pd := errors.NewNotFound("user 42 not found")
//	pd.Write(c)          // sets Content-Type, status, and JSON body
//	_ = c.Error(pd)      // or delegate to ErrorHandler middleware
//
// # Factory functions
//
//	NewBadRequest(detail, err...)        → 400
//	NewUnauthorized(detail, err...)      → 401
//	NewForbidden(detail, err...)         → 403
//	NewNotFound(detail, err...)          → 404
//	NewMethodNotAllowed(detail, err...)  → 405
//	NewConflict(detail, err...)          → 409
//	NewInternal(detail, err...)          → 500
//	NewNotImplemented(detail, err...)    → 501
//	NewBadGateway(detail, err...)        → 502
//
// The optional err variadic argument wraps an underlying cause accessible via
// errors.Unwrap / errors.Is / errors.As.
//
// # NewBadRequest and validation errors
//
// NewBadRequest has special handling for validation causes. When the wrapped
// error is a *cmerrors.ValidationError or validator.ValidationErrors, the
// fields map is extracted and set in the response "errors" field automatically.
//
// # BaseUrl and type URIs
//
// Every ProblemDetails response includes a "type" URI of the form
// <BaseUrl>/problems/<slug>. Call SetBaseUrl once at application startup
// (typically from an environment variable) before serving any requests:
//
//	errors.SetBaseUrl(os.Getenv("BASE_URL"))
//
// Note: the SetBaseUrl implementation contains a known issue where it
// temporarily assigns os.Getenv("BASE_URL") before overwriting with the
// argument — the net result is that the argument always wins. The env-var
// branch is effectively dead code.
//
// # Lazy ID population
//
// When Write(c) is called, requestId, correlationId, and traceId are read from
// the Gin context if they have not already been set on the ProblemDetails
// value. This means the ID middleware must have run before Write is called.
//
// # Problem routes
//
// RegisterProblemRoutes registers GET endpoints under /problems/* that respond
// with example ProblemDetails JSON for each status code. These endpoints are
// useful for client SDK generation and API documentation:
//
//	errors.RegisterProblemRoutes(router)
//	// GET /problems/bad-request, /problems/not-found, etc.
package errors
