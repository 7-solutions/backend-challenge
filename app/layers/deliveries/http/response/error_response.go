package response

// ErrorResponse mirrors the JSON shape of an RFC 7807 ProblemDetails response.
// It is provided as a plain DTO for client-side type assertions and API
// documentation. Actual error responses are written by
// httpx/errors.ProblemDetails.Write.
type ErrorResponse struct {
	Type          string      `json:"type"`
	Title         string      `json:"title"`
	Status        int         `json:"status"`
	Detail        string      `json:"detail,omitempty"`
	Instance      string      `json:"instance,omitempty"`
	Errors        interface{} `json:"errors,omitempty"`
	RequestID     string      `json:"requestId,omitempty"`
	CorrelationID string      `json:"correlationId,omitempty"`
	TraceId       string      `json:"traceId,omitempty"`
}
