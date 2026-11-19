package middlewares_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/chaihanij/backend-challenge/app/constants"
	middleware "github.com/chaihanij/backend-challenge/app/layers/deliveries/http/middlewares"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func TestRequestId_UsesIncomingHeader(t *testing.T) {
	t.Parallel()
	gin.SetMode(gin.TestMode)

	r := gin.New()
	r.Use(middleware.RequestId())
	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, c.GetString(constants.RequestIdKey))
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set(constants.RequestIdHeader, "req-123")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if got := w.Code; got != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, got)
	}
	if got := w.Body.String(); got != "req-123" {
		t.Fatalf("unexpected context request id: %q", got)
	}
	if got := w.Header().Get(constants.RequestIdHeader); got != "req-123" {
		t.Fatalf("unexpected response request id header: %q", got)
	}
}

func TestRequestId_GeneratesUuidWhenMissing(t *testing.T) {
	t.Parallel()
	gin.SetMode(gin.TestMode)

	r := gin.New()
	r.Use(middleware.RequestId())
	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, c.GetString(constants.RequestIdKey))
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if got := w.Code; got != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, got)
	}

	headerVal := w.Header().Get(constants.RequestIdHeader)
	if headerVal == "" {
		t.Fatal("expected generated request id header")
	}
	if _, err := uuid.Parse(headerVal); err != nil {
		t.Fatalf("request id header is not a valid UUID: %v", err)
	}
	if got := w.Body.String(); got != headerVal {
		t.Fatalf("context request id does not match response header: %q != %q", got, headerVal)
	}
}
