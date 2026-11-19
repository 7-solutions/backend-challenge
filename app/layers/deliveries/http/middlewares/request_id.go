package middlewares

import (
	"github.com/chaihanij/backend-challenge/app/constants"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// RequestId returns a middleware that propagates X-Request-ID across request and
// response. If the incoming header is empty, it generates a UUID.
func RequestId() gin.HandlerFunc {
	return func(c *gin.Context) {
		reqId := c.GetHeader(constants.RequestIdHeader)
		if reqId == "" {
			reqId = uuid.New().String()
		}
		c.Set(constants.RequestIdKey, reqId)
		c.Writer.Header().Set(constants.RequestIdHeader, reqId)
		c.Next()
	}
}
