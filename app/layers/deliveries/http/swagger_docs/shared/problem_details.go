package shared

// ProblemDetailsResponse documents the RFC 7807 error response envelope.
type ProblemDetailsResponse struct {
	Type          string      `json:"type"`
	Title         string      `json:"title"`
	Status        int         `json:"status"`
	Detail        string      `json:"detail,omitempty"`
	Instance      string      `json:"instance,omitempty"`
	Errors        interface{} `json:"errors,omitempty"`
	RequestID     string      `json:"requestId,omitempty"`
	CorrelationID string      `json:"correlationId,omitempty"`
	TraceID       string      `json:"traceId,omitempty"`
}

type CreateUserRequest struct {
	Email    string `json:"email" example:"user@example.com"`
	Name     string `json:"name" example:"John Doe"`
	Password string `json:"password" example:"secret123"`
}

type UpdateUserRequest struct {
	Email string `json:"email" example:"user@example.com"`
	Name  string `json:"name" example:"John Doe"`
}

type UserResponse struct {
	ID    string `json:"id" example:"1bf2e66a05118e911a36933125bbab4c"`
	Email string `json:"email" example:"user@example.com"`
	Name  string `json:"name" example:"John Doe"`
}

type UserResponseEnvelope struct {
	Code      int         `json:"code" example:"200"`
	Message   string      `json:"message" example:"Success"`
	Data      interface{} `json:"data"`
	Meta      interface{} `json:"meta,omitempty"`
	Timestamp string      `json:"timestamp,omitempty"`
}

type UsersResponseEnvelope struct {
	Code      int         `json:"code" example:"200"`
	Message   string      `json:"message" example:"Success"`
	Data      interface{} `json:"data"`
	Meta      interface{} `json:"meta,omitempty"`
	Timestamp string      `json:"timestamp,omitempty"`
}

type LoginMetaResponse struct {
	RequestID     string `json:"requestId,omitempty"`
	CorrelationID string `json:"correlationId,omitempty"`
	TraceID       string `json:"traceId,omitempty"`
}
