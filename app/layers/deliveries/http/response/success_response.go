package response

import (
	"net/http"
	"time"

	"github.com/chaihanij/backend-challenge/app/constants"
	"github.com/gin-gonic/gin"
)

// SuccessResponse is the standard JSON envelope for successful responses.
type SuccessResponse[T any] struct {
	Code      int    `json:"code"`
	Message   string `json:"message"`
	Data      T      `json:"data,omitempty"`
	Meta      *Meta  `json:"meta,omitempty"`
	Links     *Links `json:"links,omitempty"`
	Timestamp string `json:"timestamp,omitempty"`
}

// Meta carries request tracing identifiers and optional pagination details.
type Meta struct {
	RequestID     string      `json:"requestId,omitempty"`
	CorrelationID string      `json:"correlationId,omitempty"`
	TraceID       string      `json:"traceId,omitempty"`
	Pagination    *Pagination `json:"pagination,omitempty"`
}

// Pagination describes a paginated result set.
type Pagination struct {
	// Total is the total number of matching records.
	Total int `json:"total"`
	// Count is the number of records in the current page.
	Count       int    `json:"count"`
	PerPage     int    `json:"perPage"`
	CurrentPage int    `json:"currentPage"`
	TotalPages  int    `json:"totalPages"`
	NextCursor  string `json:"nextCursor,omitempty"`
}

// Links provides navigation URLs for paginated responses.
type Links struct {
	Self string `json:"self,omitempty"`
	Next string `json:"next,omitempty"`
	Prev string `json:"prev,omitempty"`
}

// Success writes a standard HTTP 200 response with data payload, request
// metadata, and RFC3339 timestamp.
func Success(c *gin.Context, data interface{}) {

	requestId := c.GetString(constants.RequestIdKey)
	if requestId != "" {
		c.Writer.Header().Set(constants.RequestIdHeader, requestId)
	}

	correlationId := c.GetString(constants.CorrelationIdKey)
	if correlationId != "" {
		c.Writer.Header().Set(constants.CorrelationIdHeader, correlationId)
	}

	traceId := c.GetString(constants.TraceIdKey)
	if traceId != "" {
		c.Writer.Header().Set(constants.TraceIdHeader, traceId)
	}

	c.JSON(http.StatusOK, SuccessResponse[interface{}]{
		Code:    http.StatusOK,
		Message: "OK",
		Data:    data,
		Meta: &Meta{
			RequestID:     requestId,
			CorrelationID: correlationId,
			TraceID:       traceId,
		},
		Timestamp: time.Now().Format(time.RFC3339),
	})
}

// SuccessWithPagination writes a standard HTTP 200 response with data payload,
// pagination metadata, navigation links, and RFC3339 timestamp.
func SuccessWithPagination(c *gin.Context, data interface{}, pagination *Pagination, links *Links) {

	requestId := c.GetString(constants.RequestIdKey)
	if requestId != "" {
		c.Writer.Header().Set(constants.RequestIdHeader, requestId)
	}

	correlationId := c.GetString(constants.CorrelationIdKey)
	if correlationId != "" {
		c.Writer.Header().Set(constants.CorrelationIdHeader, correlationId)
	}

	traceId := c.GetString(constants.TraceIdKey)
	if traceId != "" {
		c.Writer.Header().Set(constants.TraceIdHeader, traceId)
	}

	c.JSON(http.StatusOK, SuccessResponse[interface{}]{
		Code:    http.StatusOK,
		Message: "OK",
		Data:    data,
		Meta: &Meta{
			RequestID:     requestId,
			CorrelationID: correlationId,
			TraceID:       traceId,
			Pagination:    pagination,
		},
		Links:     links,
		Timestamp: time.Now().Format(time.RFC3339),
	})
}

// SuccessCreated writes a standard HTTP 201 response with data payload, request
// metadata, and RFC3339 timestamp.
func SuccessCreated(c *gin.Context, data interface{}) {

	requestId := c.GetString(constants.RequestIdKey)
	if requestId != "" {
		c.Writer.Header().Set(constants.RequestIdHeader, requestId)
	}

	correlationId := c.GetString(constants.CorrelationIdKey)
	if correlationId != "" {
		c.Writer.Header().Set(constants.CorrelationIdHeader, correlationId)
	}

	traceId := c.GetString(constants.TraceIdKey)
	if traceId != "" {
		c.Writer.Header().Set(constants.TraceIdHeader, traceId)
	}

	c.JSON(http.StatusCreated, SuccessResponse[interface{}]{
		Code:    http.StatusCreated,
		Message: "Created",
		Data:    data,
		Meta: &Meta{
			RequestID:     requestId,
			CorrelationID: correlationId,
			TraceID:       traceId,
		},
		Timestamp: time.Now().Format(time.RFC3339),
	})
}
