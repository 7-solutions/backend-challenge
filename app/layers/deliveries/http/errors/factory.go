package errors

import (
	"errors"

	apperrors "github.com/chaihanij/backend-challenge/app/errors"
	"github.com/chaihanij/backend-challenge/app/layers/deliveries/http/validators"
)

// NewBadRequest returns a 400 Bad Request ProblemDetails. When the optional
// err is a *apperrors.ValidationError or validator.ValidationErrors, its field
// messages are extracted and set in the response "errors" field automatically.
func NewBadRequest(detail string, err ...error) ProblemDetails {
	rootErr := firstErr(err)
	var fields map[string]string

	if rootErr != nil {
		items := validators.MapValidationErrorsSlice(rootErr)
		fields = make(map[string]string, len(items))
		for _, item := range items {
			fields[item.Field] = item.Message
		}

		if validation, ok := errors.AsType[*apperrors.ValidationError](rootErr); ok {
			fields = validation.Fields
		}
	}

	return &problemDetails{
		Type:   problemType("bad-request"),
		Title:  "Bad Request",
		Status: 400,
		Detail: detail,
		Errors: fields,
		Err:    rootErr,
	}
}

// NewUnauthorized returns a 401 Unauthorized ProblemDetails.
func NewUnauthorized(detail string, err ...error) ProblemDetails {
	return &problemDetails{
		Type:   problemType("unauthorized"),
		Title:  "Unauthorized",
		Status: 401,
		Detail: detail,
		Err:    firstErr(err),
	}
}

// NewForbidden returns a 403 Forbidden ProblemDetails.
func NewForbidden(detail string, err ...error) ProblemDetails {
	return &problemDetails{
		Type:   problemType("forbidden"),
		Title:  "Forbidden",
		Status: 403,
		Detail: detail,
		Err:    firstErr(err),
	}
}

// NewNotFound returns a 404 Not Found ProblemDetails.
func NewNotFound(detail string, err ...error) ProblemDetails {
	return &problemDetails{
		Type:   problemType("not-found"),
		Title:  "Not Found",
		Status: 404,
		Detail: detail,
		Err:    firstErr(err),
	}
}

// NewMethodNotAllowed returns a 405 Method Not Allowed ProblemDetails.
func NewMethodNotAllowed(detail string, err ...error) ProblemDetails {
	return &problemDetails{
		Type:   problemType("method-not-allowed"),
		Title:  "Method Not Allowed",
		Status: 405,
		Detail: detail,
		Err:    firstErr(err),
	}
}

// NewConflict returns a 409 Conflict ProblemDetails.
func NewConflict(detail string, err ...error) ProblemDetails {
	return &problemDetails{
		Type:   problemType("conflict"),
		Title:  "Conflict",
		Status: 409,
		Detail: detail,
		Err:    firstErr(err),
	}
}

// NewInternal returns a 500 Internal Server Error ProblemDetails.
func NewInternal(detail string, err ...error) ProblemDetails {
	return &problemDetails{
		Type:   problemType("internal-server-error"),
		Title:  "Internal Server Error",
		Status: 500,
		Detail: detail,
		Err:    firstErr(err),
	}
}

// NewNotImplemented returns a 501 Not Implemented ProblemDetails.
func NewNotImplemented(detail string, err ...error) ProblemDetails {
	return &problemDetails{
		Type:   problemType("not-implemented"),
		Title:  "Not Implemented",
		Status: 501,
		Detail: detail,
		Err:    firstErr(err),
	}
}

// NewBadGateway returns a 502 Bad Gateway ProblemDetails.
func NewBadGateway(detail string, err ...error) ProblemDetails {
	return &problemDetails{
		Type:   problemType("bad-gateway"),
		Title:  "Bad Gateway",
		Status: 502,
		Detail: detail,
		Err:    firstErr(err),
	}
}
