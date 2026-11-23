package middlewares_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/chaihanij/backend-challenge/app/constants"
	middleware "github.com/chaihanij/backend-challenge/app/layers/deliveries/http/middlewares"
	"github.com/gin-gonic/gin"
)

func TestAccessLog_LogsRequestDetailsAndRestoresRequestBody(t *testing.T) {
	gin.SetMode(gin.TestMode)

	r := gin.New()
	var entries []middleware.AccessLogEntry

	r.Use(func(c *gin.Context) {
		c.Set(string(constants.CtxKeyRequestID), "req-123")
		c.Set(string(constants.CtxKeyCorrelationID), "corr-456")
		c.Next()
	})
	r.Use(middleware.AccessLog(func(entry middleware.AccessLogEntry) {
		entries = append(entries, entry)
	}))
	r.POST("/users/:id", func(c *gin.Context) {
		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			t.Fatalf("unexpected read error: %v", err)
		}
		c.String(http.StatusCreated, string(body))
	})

	payload := "create-user-payload"
	req := httptest.NewRequest(http.MethodPost, "/users/123", strings.NewReader(payload))
	req.Header.Set("User-Agent", "copilot-test")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if got := w.Code; got != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, got)
	}
	if got := w.Body.String(); got != payload {
		t.Fatalf("expected response body %q, got %q", payload, got)
	}
	if len(entries) != 1 {
		t.Fatalf("expected 1 access log entry, got %d", len(entries))
	}

	entry := entries[0]
	if entry.Method != http.MethodPost {
		t.Errorf("expected method %s, got %s", http.MethodPost, entry.Method)
	}
	if entry.Path != "/users/:id" {
		t.Errorf("expected path %s, got %s", "/users/:id", entry.Path)
	}
	if entry.Proto != "HTTP/1.1" {
		t.Errorf("expected proto %s, got %s", "HTTP/1.1", entry.Proto)
	}
	if entry.Status != http.StatusCreated {
		t.Errorf("expected status %d, got %d", http.StatusCreated, entry.Status)
	}
	if entry.UserAgent != "copilot-test" {
		t.Errorf("expected user agent %q, got %q", "copilot-test", entry.UserAgent)
	}
	if entry.BytesIn != int64(len(payload)) {
		t.Errorf("expected bytes in %d, got %d", len(payload), entry.BytesIn)
	}
	if entry.BytesOut != int64(len(payload)) {
		t.Errorf("expected bytes out %d, got %d", len(payload), entry.BytesOut)
	}
	if entry.RequestId != "req-123" {
		t.Errorf("expected request id %q, got %q", "req-123", entry.RequestId)
	}
	if entry.CorrelationId != "corr-456" {
		t.Errorf("expected correlation id %q, got %q", "corr-456", entry.CorrelationId)
	}
	if entry.LatencyMs < 0 {
		t.Errorf("expected non-negative latency, got %d", entry.LatencyMs)
	}
}

func TestAccessLog_DoesNotLogIgnoredPath(t *testing.T) {
	gin.SetMode(gin.TestMode)

	r := gin.New()
	logged := false

	r.Use(middleware.AccessLogWithConfig(&middleware.AccessLogConfig{
		LoggerFunc: func(entry middleware.AccessLogEntry) {
			logged = true
		},
		IgnorePaths: []string{"/health"},
	}))
	r.GET("/health", func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if got := w.Code; got != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, got)
	}
	if logged {
		t.Fatal("expected ignored path not to be logged")
	}
}

func TestAccessLog_WritesJSONToDefaultWriterWhenLoggerFuncIsNil(t *testing.T) {
	gin.SetMode(gin.TestMode)

	prevWriter := gin.DefaultWriter
	buf := &bytes.Buffer{}
	gin.DefaultWriter = buf
	t.Cleanup(func() {
		gin.DefaultWriter = prevWriter
	})

	r := gin.New()
	r.Use(middleware.AccessLogWithConfig(&middleware.AccessLogConfig{}))
	r.GET("/ping", func(c *gin.Context) {
		c.Status(http.StatusNoContent)
	})

	req := httptest.NewRequest(http.MethodGet, "/ping", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if got := w.Code; got != http.StatusNoContent {
		t.Fatalf("expected status %d, got %d", http.StatusNoContent, got)
	}

	line := bytes.TrimSpace(buf.Bytes())
	if len(line) == 0 {
		t.Fatal("expected access log output in default writer")
	}

	var entry middleware.AccessLogEntry
	if err := json.Unmarshal(line, &entry); err != nil {
		t.Fatalf("expected valid json access log entry, got %v", err)
	}
	if entry.Method != http.MethodGet {
		t.Errorf("expected method %s, got %s", http.MethodGet, entry.Method)
	}
	if entry.Path != "/ping" {
		t.Errorf("expected path %s, got %s", "/ping", entry.Path)
	}
	if entry.Status != http.StatusNoContent {
		t.Errorf("expected status %d, got %d", http.StatusNoContent, entry.Status)
	}
}
