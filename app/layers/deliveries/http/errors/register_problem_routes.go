package errors

import (
	"github.com/gin-gonic/gin"
)

// RegisterProblemRoutes registers the standard /problems endpoints on the provided router.
// It returns the created RouterGroup to make the routes easier to test or further customize.
func RegisterProblemRoutes(r *gin.Engine) *gin.RouterGroup {
	return RegisterProblemRoutesWithBase(r, "/problems")
}

// RegisterProblemRoutesWithBase is like RegisterProblemRoutes but allows a custom base path.
func RegisterProblemRoutesWithBase(r *gin.Engine, basePath string) *gin.RouterGroup {
	problems := r.Group(basePath)

	// 400 Bad Request
	problems.GET("/bad-request", func(c *gin.Context) {
		NewBadRequest("bad-request").Write(c)
	})

	// 401 Unauthorized
	problems.GET("/unauthorized", func(c *gin.Context) {
		NewUnauthorized("unauthorized").Write(c)
	})

	// 403 Forbidden
	problems.GET("/forbidden", func(c *gin.Context) {
		NewForbidden("forbidden").Write(c)
	})

	// 404 Not Found
	problems.GET("/not-found", func(c *gin.Context) {
		NewNotFound("not found").Write(c)
	})

	// 500 Internal Server Error
	problems.GET("/internal-error", func(c *gin.Context) {
		NewInternal("internal error").Write(c)
	})

	return problems
}
