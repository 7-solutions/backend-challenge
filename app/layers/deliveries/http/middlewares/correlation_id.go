package middlewares

import (
	"github.com/chaihanij/backend-challenge/app/constants"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// CorrelationId returns a middleware that propagates X-Correlation-ID across
// request and response. If the incoming header is empty, it generates a UUID.
func CorrelationId() gin.HandlerFunc {
	return func(c *gin.Context) {

		correlationId := c.GetHeader(constants.CorrelationIdHeader)
		if correlationId == "" {
			correlationId = uuid.New().String()
		}
		c.Set(constants.CorrelationIdKey, correlationId)
		c.Writer.Header().Set(constants.CorrelationIdHeader, correlationId)
		c.Next()
	}
}
