package middlewares_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/chaihanij/backend-challenge/app/constants"
	"github.com/chaihanij/backend-challenge/app/layers/deliveries/http/middlewares"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func TestCorrelationId_UsesIncomingHeader(t *testing.T) {
	t.Parallel()
	gin.SetMode(gin.TestMode)

	r := gin.New()
	r.Use(middlewares.CorrelationId())
	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, c.GetString(constants.CorrelationIdKey))
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set(constants.CorrelationIdHeader, "corr-123")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if got := w.Code; got != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, got)
	}
	if got := w.Body.String(); got != "corr-123" {
		t.Fatalf("unexpected context correlation id: %q", got)
	}
	if got := w.Header().Get(constants.CorrelationIdHeader); got != "corr-123" {
		t.Fatalf("unexpected response correlation id header: %q", got)
	}
}

func TestCorrelationId_GeneratesUuidWhenMissing(t *testing.T) {
	t.Parallel()
	gin.SetMode(gin.TestMode)

	r := gin.New()
	r.Use(middlewares.CorrelationId())
	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, c.GetString(constants.CorrelationIdKey))
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if got := w.Code; got != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, got)
	}

	headerVal := w.Header().Get(constants.CorrelationIdHeader)
	if headerVal == "" {
		t.Fatal("expected generated correlation id header")
	}
	if _, err := uuid.Parse(headerVal); err != nil {
		t.Fatalf("correlation id header is not a valid UUID: %v", err)
	}
	if got := w.Body.String(); got != headerVal {
		t.Fatalf("context correlation id does not match response header: %q != %q", got, headerVal)
	}
}
