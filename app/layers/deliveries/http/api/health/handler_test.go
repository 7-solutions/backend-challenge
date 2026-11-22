package health_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	_handler "github.com/chaihanij/backend-challenge/app/layers/deliveries/http/api/health"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHandler_Health(t *testing.T) {
	t.Run("returns 200 with a status and a timestamp", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		h := _handler.NewHandler(gin.New())
		recorder := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(recorder)
		c.Request = httptest.NewRequest(http.MethodGet, "/api/health", nil)
		before := time.Now()

		h.Health(c)

		require.Equal(t, http.StatusOK, recorder.Code)
		var payload struct {
			Status    string    `json:"status"`
			Timestamp time.Time `json:"timestamp"`
		}
		require.NoError(t, json.Unmarshal(recorder.Body.Bytes(), &payload))
		assert.Equal(t, "ok", payload.Status)
		assert.False(t, payload.Timestamp.Before(before.Truncate(time.Second)))
	})

	t.Run("answers on every route through the engine", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		engine := gin.New()
		_ = _handler.NewHandler(engine)

		for _, path := range []string{"/api/health"} {
			recorder := httptest.NewRecorder()
			engine.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, path, nil))

			assert.Equal(t, http.StatusOK, recorder.Code, path)
			assert.Contains(t, recorder.Body.String(), `"status":"ok"`)
		}
	})
}
