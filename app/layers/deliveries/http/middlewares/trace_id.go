package middlewares

import (
	"github.com/chaihanij/backend-challenge/app/constants"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// TraceId returns a middleware that propagates X-Trace-ID across request and
// response. If the incoming header is empty, it generates a UUID.
func TraceId() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceId := c.GetHeader(constants.TraceIdHeader)
		if traceId == "" {
			traceId = uuid.New().String()
		}
		c.Set(constants.TraceIdKey, traceId)
		c.Writer.Header().Set(constants.TraceIdHeader, traceId)
		c.Next()
	}
}
