package health

import (
	"github.com/gin-gonic/gin"
)

// Health
//
//	@Summary		Health check
//	@Description	Returns the gateway's liveness status and current server time.
//	@ID				health
//	@Produce		json
//	@Tags			health
//	@Success		200	{object}	map[string]interface{}	"Service is healthy"
//	@Router			/api/health [get]
func Health(c *gin.Context) {
}
