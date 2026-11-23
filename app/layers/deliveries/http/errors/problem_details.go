package errors

import (
	"encoding/json"

	"github.com/chaihanij/backend-challenge/app/constants"
	"github.com/gin-gonic/gin"
)

// ProblemDetails is an RFC 7807-compliant error that also knows how to write
// itself as a Gin JSON response. Factory functions (NewBadRequest, NewNotFound,
// etc.) return this interface. Pass the value to c.Error to delegate response
// writing to ErrorHandler middleware, or call Write(c) directly.
type ProblemDetails interface {
	error
	// GetTitle returns the short human-readable summary of the problem type.
	GetTitle() string
	// GetStatus returns the HTTP status code.
	GetStatus() int
	// GetInstance returns the URI reference identifying the specific occurrence.
	GetInstance() string
	// GetDetail returns the human-readable explanation of this specific occurrence.
	GetDetail() string
	// GetErrors returns optional structured field errors (e.g. for validation failures).
	GetErrors() interface{}
	// GetRequestID returns the X-Request-ID associated with this problem.
	GetRequestID() string
	// GetCorrelationID returns the X-Correlation-ID associated with this problem.
	GetCorrelationID() string
	// GetTraceId returns the X-Trace-ID associated with this problem.
	GetTraceId() string
	// Unwrap returns the underlying cause, enabling errors.Is / errors.As chains.
	Unwrap() error
	// JSON serialises the problem as a JSON byte slice.
	JSON() []byte
	// Write sends the problem as an RFC 7807 JSON response to c. It populates
	// requestId, correlationId, and traceId from the Gin context when they have
	// not already been set, and echoes them as response headers.
	Write(c *gin.Context)
}

type ProblemDetailsResponse struct {
	Type          string      `json:"type"`
	Title         string      `json:"title"`
	Status        int         `json:"status"`
	Detail        string      `json:"detail,omitempty"`
	Instance      string      `json:"instance,omitempty"`
	Errors        interface{} `json:"errors,omitempty"`
	RequestID     string      `json:"requestId,omitempty"`
	CorrelationID string      `json:"correlationId,omitempty"`
	TraceId       string      `json:"traceId,omitempty"`
	Err           error       `json:"-"`
}

func (p *ProblemDetailsResponse) GetTitle() string         { return p.Title }
func (p *ProblemDetailsResponse) GetStatus() int           { return p.Status }
func (p *ProblemDetailsResponse) GetInstance() string      { return p.Instance }
func (p *ProblemDetailsResponse) GetDetail() string        { return p.Detail }
func (p *ProblemDetailsResponse) GetErrors() interface{}   { return p.Errors }
func (p *ProblemDetailsResponse) GetRequestID() string     { return p.RequestID }
func (p *ProblemDetailsResponse) GetCorrelationID() string { return p.CorrelationID }
func (p *ProblemDetailsResponse) GetTraceId() string       { return p.TraceId }

func (p *ProblemDetailsResponse) Unwrap() error { return p.Err }

func (p *ProblemDetailsResponse) Error() string {
	if p.Err != nil {
		return p.Detail + ": " + p.Err.Error()
	}
	return p.Detail
}

func (p *ProblemDetailsResponse) JSON() []byte {
	b, _ := json.Marshal(p)
	return b
}

func (p *ProblemDetailsResponse) Write(c *gin.Context) {
	if c.Writer.Written() {
		return
	}

	if p.RequestID == "" {
		p.RequestID = c.GetString(constants.RequestIdKey)
	}
	if p.CorrelationID == "" {
		p.CorrelationID = c.GetString(constants.CorrelationIdKey)
	}
	if p.TraceId == "" {
		p.TraceId = c.GetString(constants.TraceIdKey)
	}
	if p.Instance == "" {
		p.Instance = c.Request.URL.Path
	}

	if p.RequestID != "" {
		c.Writer.Header().Set(constants.RequestIdHeader, p.RequestID)
	}
	if p.CorrelationID != "" {
		c.Writer.Header().Set(constants.CorrelationIdHeader, p.CorrelationID)
	}
	if p.TraceId != "" {
		c.Writer.Header().Set(constants.TraceIdHeader, p.TraceId)
	}

	c.JSON(p.Status, p)
}
