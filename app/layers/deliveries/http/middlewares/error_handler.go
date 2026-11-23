package middlewares

import (
	"errors"

	apperrors "github.com/chaihanij/backend-challenge/app/errors"
	httperrors "github.com/chaihanij/backend-challenge/app/layers/deliveries/http/errors"

	"github.com/gin-gonic/gin"
)

// ErrorHandlerConfig holds configuration for ErrorHandler middleware
type ErrorHandlerConfig struct {
	IgnorePaths []string
}

// ErrorHandler returns a middleware that converts the last Gin error into a
// normalized RFC7807 ProblemDetails response. If the normalized problem includes
// an internal error, it is also written to the provided logger.
//
// Example with ignore paths:
//
//	r.Use(middleware.ErrorHandlerWithConfig(&middleware.ErrorHandlerConfig{
//	    IgnorePaths: []string{"/health", "/metrics*"},
//	}))
func ErrorHandler() gin.HandlerFunc {
	return ErrorHandlerWithConfig(&ErrorHandlerConfig{
		IgnorePaths: []string{},
	})
}

// ErrorHandlerWithConfig returns a middleware that converts errors to RFC7807 responses
// with custom configuration for ignoring certain paths
func ErrorHandlerWithConfig(config *ErrorHandlerConfig) gin.HandlerFunc {
	if config == nil {
		config = &ErrorHandlerConfig{}
	}

	return func(c *gin.Context) {
		c.Next()

		if IsPathIgnored(c.Request.URL.Path, config.IgnorePaths) {
			return
		}

		last := c.Errors.Last()
		if last == nil {
			return
		}

		domainErrorToProblemDetails(last.Err).Write(c)
	}
}

// domainErrorToProblemDetails maps domain errors to RFC7807 ProblemDetails
func domainErrorToProblemDetails(err error) httperrors.ProblemDetails {
	if pd, ok := errors.AsType[httperrors.ProblemDetails](err); ok {
		return pd
	}

	if unauthorized, ok := errors.AsType[*apperrors.UnauthorizedError](err); ok {
		return httperrors.NewUnauthorized(unauthorized.Error(), err)
	}

	if notFound, ok := errors.AsType[*apperrors.NotFoundError](err); ok {
		return httperrors.NewNotFound(notFound.Error())
	}

	if validation, ok := errors.AsType[*apperrors.ValidationError](err); ok {
		return httperrors.NewBadRequest(err.Error(), validation)
	}

	if conflict, ok := errors.AsType[*apperrors.ConflictError](err); ok {
		return httperrors.NewConflict(conflict.Error(), err)
	}

	if permDenied, ok := errors.AsType[*apperrors.PermissionDeniedError](err); ok {
		return httperrors.NewForbidden(permDenied.Error())
	}

	if rateLimit, ok := errors.AsType[*apperrors.RateLimitError](err); ok {
		return httperrors.NewInternal(rateLimit.Error())
	}

	if dependency, ok := errors.AsType[*apperrors.DependencyError](err); ok {
		return httperrors.NewInternal(dependency.Error())
	}

	if internal, ok := errors.AsType[*apperrors.InternalError](err); ok {
		return httperrors.NewInternal(internal.Error())
	}

	return httperrors.NewInternal(err.Error())
}
