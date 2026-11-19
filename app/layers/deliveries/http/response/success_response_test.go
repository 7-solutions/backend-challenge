package response_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/chaihanij/backend-challenge/app/constants"
	"github.com/chaihanij/backend-challenge/app/layers/deliveries/http/response"
	"github.com/gin-gonic/gin"
)

func newResponseTestContext(t *testing.T) (*gin.Context, *httptest.ResponseRecorder) {
	t.Helper()
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	return c, w
}

func TestSuccess(t *testing.T) {
	t.Parallel()

	c, w := newResponseTestContext(t)
	c.Set(constants.RequestIdKey, "req-123")
	c.Set(constants.CorrelationIdKey, "corr-123")

	response.Success(c, gin.H{"name": "alice"})

	if got := w.Code; got != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, got)
	}

	var body response.SuccessResponse[interface{}]
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("failed to decode body: %v", err)
	}

	if body.Code != http.StatusOK {
		t.Fatalf("unexpected code: %d", body.Code)
	}
	if body.Message != "OK" {
		t.Fatalf("unexpected message: %s", body.Message)
	}
	if body.Meta == nil {
		t.Fatal("expected meta")
	}
	if body.Meta.RequestID != "req-123" {
		t.Fatalf("unexpected request id: %s", body.Meta.RequestID)
	}
	if body.Meta.CorrelationID != "corr-123" {
		t.Fatalf("unexpected correlation id: %s", body.Meta.CorrelationID)
	}
	if _, err := time.Parse(time.RFC3339, body.Timestamp); err != nil {
		t.Fatalf("timestamp is not RFC3339: %v", err)
	}

	data, ok := body.Data.(map[string]interface{})
	if !ok {
		t.Fatalf("expected map data, got %T", body.Data)
	}
	if got := data["name"]; got != "alice" {
		t.Fatalf("unexpected data name: %v", got)
	}
}

func TestSuccessWithPagination(t *testing.T) {
	t.Parallel()

	c, w := newResponseTestContext(t)
	c.Set(constants.RequestIdKey, "req-456")
	c.Set(constants.CorrelationIdKey, "corr-456")

	pagination := &response.Pagination{
		Total:       100,
		Count:       10,
		PerPage:     10,
		CurrentPage: 1,
		TotalPages:  10,
		NextCursor:  "cursor-2",
	}
	links := &response.Links{
		Self: "/v1/users?page=1",
		Next: "/v1/users?page=2",
		Prev: "",
	}

	response.SuccessWithPagination(c, []string{"alice", "bob"}, pagination, links)

	if got := w.Code; got != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, got)
	}

	var body response.SuccessResponse[interface{}]
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("failed to decode body: %v", err)
	}

	if body.Meta == nil {
		t.Fatal("expected meta")
	}
	if body.Meta.RequestID != "req-456" {
		t.Fatalf("unexpected request id: %s", body.Meta.RequestID)
	}
	if body.Meta.CorrelationID != "corr-456" {
		t.Fatalf("unexpected correlation id: %s", body.Meta.CorrelationID)
	}
	if body.Meta.Pagination == nil {
		t.Fatal("expected pagination")
	}
	if body.Meta.Pagination.Total != 100 || body.Meta.Pagination.NextCursor != "cursor-2" {
		t.Fatalf("unexpected pagination: %+v", body.Meta.Pagination)
	}
	if body.Links == nil {
		t.Fatal("expected links")
	}
	if body.Links.Self != "/v1/users?page=1" || body.Links.Next != "/v1/users?page=2" {
		t.Fatalf("unexpected links: %+v", body.Links)
	}
	if _, err := time.Parse(time.RFC3339, body.Timestamp); err != nil {
		t.Fatalf("timestamp is not RFC3339: %v", err)
	}

	data, ok := body.Data.([]interface{})
	if !ok {
		t.Fatalf("expected array data, got %T", body.Data)
	}
	if len(data) != 2 {
		t.Fatalf("unexpected data length: %d", len(data))
	}
}
