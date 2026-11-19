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

type problemDetails struct {
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

func (p *problemDetails) GetTitle() string         { return p.Title }
func (p *problemDetails) GetStatus() int           { return p.Status }
func (p *problemDetails) GetInstance() string      { return p.Instance }
func (p *problemDetails) GetDetail() string        { return p.Detail }
func (p *problemDetails) GetErrors() interface{}   { return p.Errors }
func (p *problemDetails) GetRequestID() string     { return p.RequestID }
func (p *problemDetails) GetCorrelationID() string { return p.CorrelationID }
func (p *problemDetails) GetTraceId() string       { return p.TraceId }

func (p *problemDetails) Unwrap() error { return p.Err }

func (p *problemDetails) Error() string {
	if p.Err != nil {
		return p.Detail + ": " + p.Err.Error()
	}
	return p.Detail
}

func (p *problemDetails) JSON() []byte {
	b, _ := json.Marshal(p)
	return b
}

func (p *problemDetails) Write(c *gin.Context) {
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
