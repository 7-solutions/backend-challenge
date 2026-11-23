package middlewares_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	apperrors "github.com/chaihanij/backend-challenge/app/errors"
	httperrors "github.com/chaihanij/backend-challenge/app/layers/deliveries/http/errors"
	"github.com/chaihanij/backend-challenge/app/layers/deliveries/http/middlewares"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// helper: build a router with ErrorHandler and a single GET handler that injects err
func newRouter(mw gin.HandlerFunc, path string, injectedErr error) *gin.Engine {
	r := gin.New()
	r.Use(mw)
	r.GET(path, func(c *gin.Context) {
		_ = c.Error(injectedErr)
	})
	return r
}

func doGet(r *gin.Engine, path string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(http.MethodGet, path, nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

// ── ErrorHandler (default) ────────────────────────────────────────────────────

func TestErrorHandler_Default_NoError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(middlewares.ErrorHandler())
	r.GET("/ok", func(c *gin.Context) { c.Status(http.StatusOK) })

	w := doGet(r, "/ok")
	require.Equal(t, http.StatusOK, w.Code)
}

// ── ErrorHandlerWithConfig – nil config ──────────────────────────────────────

func TestErrorHandlerWithConfig_NilConfig(t *testing.T) {
	gin.SetMode(gin.TestMode)
	//nolint:staticcheck
	r := newRouter(middlewares.ErrorHandlerWithConfig(nil), "/test", errors.New("boom"))
	w := doGet(r, "/test")
	require.Equal(t, http.StatusInternalServerError, w.Code)
}

// ── IgnorePaths ───────────────────────────────────────────────────────────────

func TestErrorHandler_IgnorePaths(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name        string
		path        string
		ignorePaths []string
		shouldError bool
	}{
		{
			name:        "error on path not in ignore list",
			path:        "/api/users",
			ignorePaths: []string{"/health", "/metrics*"},
			shouldError: true,
		},
		{
			name:        "error ignored on health endpoint",
			path:        "/health",
			ignorePaths: []string{"/health", "/metrics*"},
			shouldError: false,
		},
		{
			name:        "error ignored on metrics with pattern",
			path:        "/metrics/cpu",
			ignorePaths: []string{"/health", "/metrics*"},
			shouldError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := gin.New()
			r.Use(middlewares.ErrorHandlerWithConfig(&middlewares.ErrorHandlerConfig{
				IgnorePaths: tt.ignorePaths,
			}))

			r.GET(tt.path, func(c *gin.Context) {
				_ = c.Error(httperrors.NewForbidden("access denied"))
			})

			w := doGet(r, tt.path)

			if tt.shouldError {
				require.Equal(t, http.StatusForbidden, w.Code)
				require.NotEmpty(t, w.Body.String())
			} else {
				// If path is ignored, error should not be written back
				require.Equal(t, http.StatusOK, w.Code)
			}
		})
	}
}

func TestErrorHandlerWithConfig_NoIgnorePaths(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(middlewares.ErrorHandlerWithConfig(&middlewares.ErrorHandlerConfig{
		IgnorePaths: []string{},
	}))

	r.GET("/test", func(c *gin.Context) {
		_ = c.Error(httperrors.NewNotFound("resource not found"))
	})

	w := doGet(r, "/test")

	require.Equal(t, http.StatusNotFound, w.Code)
	require.NotEmpty(t, w.Body.String())
}

// ── domainErrorToProblemDetails – each error branch ──────────────────────────

func TestErrorHandler_ProblemDetailsPassthrough(t *testing.T) {
	gin.SetMode(gin.TestMode)
	// error is already a ProblemDetails (e.g. NewUnauthorized)
	r := newRouter(middlewares.ErrorHandler(), "/test", httperrors.NewUnauthorized("token expired"))
	w := doGet(r, "/test")
	require.Equal(t, http.StatusUnauthorized, w.Code)

	var body map[string]any
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &body))
	assert.Equal(t, float64(http.StatusUnauthorized), body["status"])
}

func TestErrorHandler_NotFoundError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	err := apperrors.NewNotFoundError("User", "42")
	r := newRouter(middlewares.ErrorHandler(), "/test", err)
	w := doGet(r, "/test")
	require.Equal(t, http.StatusNotFound, w.Code)

	var body map[string]any
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &body))
	assert.Equal(t, float64(http.StatusNotFound), body["status"])
}

func TestErrorHandler_ValidationError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	err := apperrors.NewValidationError(map[string]string{"name": "required"})
	r := newRouter(middlewares.ErrorHandler(), "/test", err)
	w := doGet(r, "/test")
	require.Equal(t, http.StatusBadRequest, w.Code)

	var body map[string]any
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &body))
	assert.Equal(t, float64(http.StatusBadRequest), body["status"])
}

func TestErrorHandler_ConflictError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	err := apperrors.NewConflictError("User", "email already exists")
	r := newRouter(middlewares.ErrorHandler(), "/test", err)
	w := doGet(r, "/test")
	require.Equal(t, http.StatusConflict, w.Code)

	var body map[string]any
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &body))
	assert.Equal(t, float64(http.StatusConflict), body["status"])
}

func TestErrorHandler_PermissionDeniedError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	err := apperrors.NewPermissionDenied("delete", "post")
	r := newRouter(middlewares.ErrorHandler(), "/test", err)
	w := doGet(r, "/test")
	require.Equal(t, http.StatusForbidden, w.Code)

	var body map[string]any
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &body))
	assert.Equal(t, float64(http.StatusForbidden), body["status"])
}

func TestErrorHandler_RateLimitError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	err := apperrors.NewRateLimitError(100, 60)
	r := newRouter(middlewares.ErrorHandler(), "/test", err)
	w := doGet(r, "/test")
	require.Equal(t, http.StatusInternalServerError, w.Code)

	var body map[string]any
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &body))
	assert.Equal(t, float64(http.StatusInternalServerError), body["status"])
}

func TestErrorHandler_DependencyError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	err := apperrors.NewDependencyError("postgres", "query", errors.New("connection refused"))
	r := newRouter(middlewares.ErrorHandler(), "/test", err)
	w := doGet(r, "/test")
	require.Equal(t, http.StatusInternalServerError, w.Code)

	var body map[string]any
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &body))
	assert.Equal(t, float64(http.StatusInternalServerError), body["status"])
}

func TestErrorHandler_InternalError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	err := apperrors.NewInternalError("service", "process", errors.New("unexpected"))
	r := newRouter(middlewares.ErrorHandler(), "/test", err)
	w := doGet(r, "/test")
	require.Equal(t, http.StatusInternalServerError, w.Code)

	var body map[string]any
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &body))
	assert.Equal(t, float64(http.StatusInternalServerError), body["status"])
}

func TestErrorHandler_GenericError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := newRouter(middlewares.ErrorHandler(), "/test", errors.New("something went wrong"))
	w := doGet(r, "/test")
	require.Equal(t, http.StatusInternalServerError, w.Code)

	var body map[string]any
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &body))
	assert.Equal(t, float64(http.StatusInternalServerError), body["status"])
}

func TestErrorHandler_WrappedNotFoundError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	cause := errors.New("db row not found")
	err := apperrors.WrapNotFound("Order", "99", cause)
	r := newRouter(middlewares.ErrorHandler(), "/test", err)
	w := doGet(r, "/test")
	require.Equal(t, http.StatusNotFound, w.Code)
}

func TestErrorHandler_WrappedConflictError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	cause := errors.New("unique constraint")
	err := apperrors.WrapConflict("Product", "SKU duplicate", cause)
	r := newRouter(middlewares.ErrorHandler(), "/test", err)
	w := doGet(r, "/test")
	require.Equal(t, http.StatusConflict, w.Code)
}

func TestErrorHandler_ResponseBodyIsJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := newRouter(middlewares.ErrorHandler(), "/test", errors.New("oops"))
	w := doGet(r, "/test")
	require.Equal(t, "application/json; charset=utf-8", w.Header().Get("Content-Type"))
	require.True(t, json.Valid(w.Body.Bytes()), "response body should be valid JSON")
}
