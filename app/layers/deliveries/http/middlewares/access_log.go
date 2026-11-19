package middlewares

import (
	"bytes"
	"encoding/json"
	"io"
	"time"

	"github.com/chaihanij/backend-challenge/app/constants"
	"github.com/gin-gonic/gin"
)

// AccessLogEntry holds the request/response details captured by AccessLog middleware.
// It is passed to the caller-supplied LoggerFunc for structured logging.
type AccessLogEntry struct {
	Method        string `json:"method"`
	Path          string `json:"path"`
	Proto         string `json:"proto,omitempty"`
	Status        int    `json:"status"`
	LatencyMs     int64  `json:"latencyMs"`
	ClientIP      string `json:"clientIp,omitempty"`
	UserAgent     string `json:"userAgent,omitempty"`
	RequestId     string `json:"requestId,omitempty"`
	CorrelationId string `json:"correlationId,omitempty"`
	TraceId       string `json:"traceId,omitempty"`
	BytesIn       int64  `json:"bytesIn,omitempty"`
	BytesOut      int64  `json:"bytesOut,omitempty"`
}

// AccessLogConfig holds configuration for AccessLog middleware
type AccessLogConfig struct {
	LoggerFunc  func(entry AccessLogEntry)
	IgnorePaths []string
}

// AccessLog returns middleware that logs HTTP requests. Supports ignoring certain paths.
//
// Example:
//
//	r.Use(middleware.AccessLog(func(entry middleware.AccessLogEntry) {
//	    // log entry
//	}))
//
// With ignore paths:
//
//	r.Use(middleware.AccessLogWithConfig(&middleware.AccessLogConfig{
//	    LoggerFunc: func(entry middleware.AccessLogEntry) { /* */ },
//	    IgnorePaths: []string{"/health", "/metrics*"},
//	}))
func AccessLog(loggerFunc func(entry AccessLogEntry)) gin.HandlerFunc {
	return AccessLogWithConfig(&AccessLogConfig{
		LoggerFunc:  loggerFunc,
		IgnorePaths: []string{},
	})
}

// AccessLogWithConfig returns middleware that logs HTTP requests with custom config
func AccessLogWithConfig(config *AccessLogConfig) gin.HandlerFunc {
	if config == nil {
		config = &AccessLogConfig{}
	}

	return func(c *gin.Context) {
		start := time.Now()
		restoreRequestBody(c)

		c.Next()

		if IsPathIgnored(c.Request.URL.Path, config.IgnorePaths) {
			return
		}

		entry := buildAccessLogEntry(c, start)

		if config.LoggerFunc != nil {
			config.LoggerFunc(entry)
			return
		}

		b, _ := json.Marshal(entry)
		_, _ = gin.DefaultWriter.Write(append(b, '\n'))
	}
}

func restoreRequestBody(c *gin.Context) {
	if c.Request == nil || c.Request.Body == nil {
		return
	}

	reqBody, _ := io.ReadAll(c.Request.Body)
	c.Request.Body = io.NopCloser(bytes.NewBuffer(reqBody))
}

func buildAccessLogEntry(c *gin.Context, start time.Time) AccessLogEntry {
	entry := AccessLogEntry{
		Method:    c.Request.Method,
		Path:      c.FullPath(),
		Proto:     c.Request.Proto,
		Status:    c.Writer.Status(),
		LatencyMs: time.Since(start).Milliseconds(),
		ClientIP:  c.ClientIP(),
		UserAgent: c.Request.UserAgent(),
		BytesIn:   c.Request.ContentLength,
		BytesOut:  int64(c.Writer.Size()),
	}

	if requestID, ok := contextString(c, constants.CtxKeyRequestID); ok {
		entry.RequestId = requestID
	}

	if correlationID, ok := contextString(c, constants.CtxKeyCorrelationID); ok {
		entry.CorrelationId = correlationID
	}

	if traceID, ok := contextString(c, constants.CtxKeyTraceID); ok {
		entry.TraceId = traceID
	}

	return entry
}

func contextString(c *gin.Context, key constants.ContextKey) (string, bool) {
	value, ok := c.Get(string(key))
	if !ok {
		return "", false
	}

	result, ok := value.(string)
	return result, ok
}
