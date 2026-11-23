package errors_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"testing"

	apperrors "github.com/chaihanij/backend-challenge/app/errors"
	httperrors "github.com/chaihanij/backend-challenge/app/layers/deliveries/http/errors"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFactoryFunctions_CreateExpectedProblemDetails(t *testing.T) {
	httperrors.BaseUrl = "https://example.com"

	firstErr := errors.New("first error")
	secondErr := errors.New("second error")

	tests := []struct {
		name       string
		factory    func(string, ...error) httperrors.ProblemDetails
		detail     string
		title      string
		status     int
		typeSuffix string
	}{
		{name: "BadRequest", factory: httperrors.NewBadRequest, detail: "bad payload", title: "Bad Request", status: 400, typeSuffix: "bad-request"},
		{name: "Unauthorized", factory: httperrors.NewUnauthorized, detail: "missing token", title: "Unauthorized", status: 401, typeSuffix: "unauthorized"},
		{name: "Forbidden", factory: httperrors.NewForbidden, detail: "insufficient permission", title: "Forbidden", status: 403, typeSuffix: "forbidden"},
		{name: "NotFound", factory: httperrors.NewNotFound, detail: "resource not found", title: "Not Found", status: 404, typeSuffix: "not-found"},
		{name: "MethodNotAllowed", factory: httperrors.NewMethodNotAllowed, detail: "method not allowed", title: "Method Not Allowed", status: 405, typeSuffix: "method-not-allowed"},
		{name: "Conflict", factory: httperrors.NewConflict, detail: "version conflict", title: "Conflict", status: 409, typeSuffix: "conflict"},
		{name: "Internal", factory: httperrors.NewInternal, detail: "internal error", title: "Internal Server Error", status: 500, typeSuffix: "internal-server-error"},
		{name: "NotImplemented", factory: httperrors.NewNotImplemented, detail: "feature unavailable", title: "Not Implemented", status: 501, typeSuffix: "not-implemented"},
		{name: "BadGateway", factory: httperrors.NewBadGateway, detail: "upstream failed", title: "Bad Gateway", status: 502, typeSuffix: "bad-gateway"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			pd := tc.factory(tc.detail, firstErr, secondErr)
			require.NotNil(t, pd)

			assert.Equal(t, tc.title, pd.GetTitle())
			assert.Equal(t, tc.status, pd.GetStatus())
			assert.Equal(t, tc.detail, pd.GetDetail())
			assert.Equal(t, firstErr, pd.Unwrap())
			assert.Equal(t, fmt.Sprintf("%s: %s", tc.detail, firstErr.Error()), pd.Error())

			var payload map[string]interface{}
			require.NoError(t, json.Unmarshal(pd.JSON(), &payload))
			assert.Equal(t, tc.title, payload["title"])
			assert.Equal(t, float64(tc.status), payload["status"])
			assert.Equal(t, fmt.Sprintf("https://example.com/problems/%s", tc.typeSuffix), payload["type"])
		})
	}
}

func TestFactoryFunctions_WithoutWrappedError(t *testing.T) {
	pd := httperrors.NewBadRequest("detail only")

	assert.Nil(t, pd.Unwrap())
	assert.Equal(t, "detail only", pd.Error())
}

// ---------- NewBadRequest-specific tests ----------

func TestNewBadRequest_WithPlainError_NoFieldErrors(t *testing.T) {
	httperrors.BaseUrl = "https://example.com"

	plainErr := errors.New("something went wrong")
	pd := httperrors.NewBadRequest("bad payload", plainErr)

	require.NotNil(t, pd)
	assert.Equal(t, "Bad Request", pd.GetTitle())
	assert.Equal(t, 400, pd.GetStatus())
	assert.Equal(t, "bad payload", pd.GetDetail())
	assert.Equal(t, plainErr, pd.Unwrap())

	// plain error is not a validator.ValidationErrors → errors map should be empty
	errs, ok := pd.GetErrors().(map[string]string)
	require.True(t, ok, "GetErrors() should return map[string]string")
	assert.Empty(t, errs)

	var payload map[string]interface{}
	require.NoError(t, json.Unmarshal(pd.JSON(), &payload))
	assert.Equal(t, "https://example.com/problems/bad-request", payload["type"])
	// plain error → errors map is initialised but empty, serialises as {}
	errorsInPayload, ok := payload["errors"].(map[string]interface{})
	require.True(t, ok, "errors in JSON should be an empty object")
	assert.Empty(t, errorsInPayload)
}

func TestNewBadRequest_WithValidationErrors_PopulatesFieldErrors(t *testing.T) {
	httperrors.BaseUrl = "https://example.com"

	type LoginRequest struct {
		Email    string `validate:"required,email"    json:"email"`
		Password string `validate:"required,min=8"    json:"password"`
	}

	v := validator.New()
	// Register json tag names so field paths match json keys
	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := fld.Tag.Get("json")
		if name == "" || name == "-" {
			return fld.Name
		}
		return strings.SplitN(name, ",", 2)[0]
	})

	validationErr := v.Struct(LoginRequest{}) // both fields empty → two errors
	require.Error(t, validationErr)

	pd := httperrors.NewBadRequest("validation failed", validationErr)

	require.NotNil(t, pd)
	assert.Equal(t, "Bad Request", pd.GetTitle())
	assert.Equal(t, 400, pd.GetStatus())
	assert.Equal(t, "validation failed", pd.GetDetail())
	assert.Equal(t, validationErr, pd.Unwrap())

	// errors map must contain entries for both fields
	errs, ok := pd.GetErrors().(map[string]string)
	require.True(t, ok, "GetErrors() should return map[string]string")
	assert.Len(t, errs, 2)
	assert.Contains(t, errs, "email")
	assert.Contains(t, errs, "password")

	// JSON must serialise the errors object
	var payload map[string]interface{}
	require.NoError(t, json.Unmarshal(pd.JSON(), &payload))
	assert.Equal(t, "https://example.com/problems/bad-request", payload["type"])
	errorsInPayload, ok := payload["errors"].(map[string]interface{})
	require.True(t, ok, "errors in JSON should be an object")
	assert.Contains(t, errorsInPayload, "email")
	assert.Contains(t, errorsInPayload, "password")
}

func TestNewBadRequest_WithSingleValidationError(t *testing.T) {
	httperrors.BaseUrl = "https://example.com"

	type CreateUserRequest struct {
		Username string `validate:"required" json:"username"`
	}

	v := validator.New()
	validationErr := v.Struct(CreateUserRequest{}) // username empty → one error
	require.Error(t, validationErr)

	pd := httperrors.NewBadRequest("invalid request", validationErr)

	errs, ok := pd.GetErrors().(map[string]string)
	require.True(t, ok)
	assert.Len(t, errs, 1)
	assert.Contains(t, errs, "Username") // no json tag registered → Go field name
	assert.Equal(t, "is required", errs["Username"])
}

func TestNewBadRequest_WithoutError_EmptyFieldErrors(t *testing.T) {
	pd := httperrors.NewBadRequest("missing body")

	assert.Nil(t, pd.Unwrap())
	assert.Equal(t, "missing body", pd.Error())

	errs, ok := pd.GetErrors().(map[string]string)
	require.True(t, ok)
	assert.Empty(t, errs)
}

// TestNewBadRequest_WithCmValidationError verifies that when a
// apperrors.ValidationError is passed the fields map from that error
// takes precedence over any go-playground/validator derived fields.
func TestNewBadRequest_WithCmValidationError(t *testing.T) {
	httperrors.BaseUrl = "https://example.com"

	cmErr := apperrors.NewValidationError(map[string]string{
		"email": "must be a valid email address",
		"age":   "must be at least 18",
	})

	pd := httperrors.NewBadRequest("validation failed", cmErr)

	require.NotNil(t, pd)
	assert.Equal(t, "Bad Request", pd.GetTitle())
	assert.Equal(t, 400, pd.GetStatus())
	assert.Equal(t, "validation failed", pd.GetDetail())
	assert.Equal(t, cmErr, pd.Unwrap())

	errs, ok := pd.GetErrors().(map[string]string)
	require.True(t, ok, "GetErrors() should return map[string]string")
	assert.Len(t, errs, 2)
	assert.Equal(t, "must be a valid email address", errs["email"])
	assert.Equal(t, "must be at least 18", errs["age"])

	var payload map[string]interface{}
	require.NoError(t, json.Unmarshal(pd.JSON(), &payload))
	assert.Equal(t, "https://example.com/problems/bad-request", payload["type"])
	errorsInPayload, ok := payload["errors"].(map[string]interface{})
	require.True(t, ok)
	assert.Equal(t, "must be a valid email address", errorsInPayload["email"])
	assert.Equal(t, "must be at least 18", errorsInPayload["age"])
}

// TestNewBadRequest_WithWrappedCmValidationError verifies that a
// apperrors.ValidationError wrapped inside another error is still detected
// via errors.As and its fields are correctly populated.
func TestNewBadRequest_WithWrappedCmValidationError(t *testing.T) {
	httperrors.BaseUrl = "https://example.com"

	inner := apperrors.NewValidationError(map[string]string{
		"username": "must not be empty",
	})
	wrapped := fmt.Errorf("handler: %w", inner)

	pd := httperrors.NewBadRequest("bad request", wrapped)

	errs, ok := pd.GetErrors().(map[string]string)
	require.True(t, ok)
	assert.Len(t, errs, 1)
	assert.Equal(t, "must not be empty", errs["username"])
}

// TestNewBadRequest_MultipleErrors_OnlyFirstIsWrapped verifies that when
// multiple errors are provided only the first one is stored as the wrapped
// error (Unwrap), while the others are ignored.
func TestNewBadRequest_MultipleErrors_OnlyFirstIsWrapped(t *testing.T) {
	firstErr := errors.New("first error")
	secondErr := errors.New("second error")
	thirdErr := errors.New("third error")

	pd := httperrors.NewBadRequest("multiple errors", firstErr, secondErr, thirdErr)

	assert.Equal(t, firstErr, pd.Unwrap())
	assert.NotEqual(t, secondErr, pd.Unwrap())
	assert.Equal(t, "multiple errors: first error", pd.Error())
}

// TestNewBadRequest_WithCmValidationError_EmptyFields verifies that a
// apperrors.ValidationError with an empty fields map results in an empty
// errors map (not nil) in the response.
func TestNewBadRequest_WithCmValidationError_EmptyFields(t *testing.T) {
	cmErr := apperrors.NewValidationError(map[string]string{})

	pd := httperrors.NewBadRequest("no field errors", cmErr)

	errs, ok := pd.GetErrors().(map[string]string)
	require.True(t, ok)
	assert.Empty(t, errs)
}

// TestNewBadRequest_ErrorString verifies the Error() string format with and
// without a wrapped error.
func TestNewBadRequest_ErrorString(t *testing.T) {
	cause := errors.New("decode failed")
	pd := httperrors.NewBadRequest("invalid json", cause)
	assert.Equal(t, "invalid json: decode failed", pd.Error())

	pdNoErr := httperrors.NewBadRequest("invalid json")
	assert.Equal(t, "invalid json", pdNoErr.Error())
}
