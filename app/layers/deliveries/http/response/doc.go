// Package response provides standard HTTP response helpers for Gin handlers,
// producing a consistent success envelope across all services.
//
// # Success response shape
//
// Success and SuccessWithPagination write HTTP 200 with the following JSON:
//
//	{
//	  "code": 200,
//	  "message": "OK",
//	  "data": <any>,
//	  "meta": {
//	    "requestId": "...",
//	    "correlationId": "...",
//	    "traceId": "...",
//	    "pagination": { ... }   // only with SuccessWithPagination
//	  },
//	  "links": { "self": "...", "next": "...", "prev": "..." },  // optional
//	  "timestamp": "2006-01-02T15:04:05Z07:00"
//	}
//
// The requestId, correlationId, and traceId values are read from the Gin
// context. They are also set as response headers (X-Request-ID etc.) so
// clients can correlate responses without parsing the body. This requires the
// RequestId, CorrelationId, and TraceId middleware to have run before the
// handler writes its response.
//
// # ErrorResponse
//
// ErrorResponse is a plain DTO that mirrors the ProblemDetails JSON shape. It
// is provided for client-side type assertions and documentation purposes.
// Actual error responses are written by httpx/errors.ProblemDetails.Write.
package response
