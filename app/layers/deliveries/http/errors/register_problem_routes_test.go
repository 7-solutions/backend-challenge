package errors_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	httperrors "github.com/chaihanij/backend-challenge/app/layers/deliveries/http/errors"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRegisterProblemRoutes_DefaultBasePath(t *testing.T) {
	gin.SetMode(gin.TestMode)
	httperrors.BaseUrl = "https://example.com"

	router := gin.New()
	group := httperrors.RegisterProblemRoutes(router)
	require.NotNil(t, group)
	assert.Equal(t, "/problems", group.BasePath())

	tests := []struct {
		path   string
		status int
		title  string
		detail string
	}{
		{path: "/problems/bad-request", status: 400, title: "Bad Request", detail: "bad-request"},
		{path: "/problems/unauthorized", status: 401, title: "Unauthorized", detail: "unauthorized"},
		{path: "/problems/forbidden", status: 403, title: "Forbidden", detail: "forbidden"},
		{path: "/problems/not-found", status: 404, title: "Not Found", detail: "not found"},
		{path: "/problems/internal-error", status: 500, title: "Internal Server Error", detail: "internal error"},
	}

	for _, tc := range tests {
		t.Run(tc.path, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tc.path, nil)
			rec := httptest.NewRecorder()

			router.ServeHTTP(rec, req)

			require.Equal(t, tc.status, rec.Code)
			assert.Contains(t, rec.Header().Get("Content-Type"), "application/json")

			var payload map[string]interface{}
			require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &payload))
			assert.Equal(t, tc.title, payload["title"])
			assert.Equal(t, tc.detail, payload["detail"])
			assert.Equal(t, tc.path, payload["instance"])
		})
	}
}

func TestRegisterProblemRoutesWithBase_CustomBasePath(t *testing.T) {
	gin.SetMode(gin.TestMode)
	httperrors.BaseUrl = "https://example.com"

	router := gin.New()
	group := httperrors.RegisterProblemRoutesWithBase(router, "/errors")
	require.NotNil(t, group)
	assert.Equal(t, "/errors", group.BasePath())

	req := httptest.NewRequest(http.MethodGet, "/errors/forbidden", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	require.Equal(t, http.StatusForbidden, rec.Code)

	var payload map[string]interface{}
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &payload))
	assert.Equal(t, "Forbidden", payload["title"])
	assert.Equal(t, "forbidden", payload["detail"])
	assert.Equal(t, "/errors/forbidden", payload["instance"])
}
