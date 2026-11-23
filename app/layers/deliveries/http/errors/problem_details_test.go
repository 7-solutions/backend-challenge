package errors_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/chaihanij/backend-challenge/app/constants"
	httperrors "github.com/chaihanij/backend-challenge/app/layers/deliveries/http/errors"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProblemDetails_GettersAndError(t *testing.T) {
	httperrors.BaseUrl = "https://example.com"

	wrapped := errors.New("database timeout")
	pd := httperrors.NewInternal("request failed", wrapped)

	assert.Equal(t, "Internal Server Error", pd.GetTitle())
	assert.Equal(t, 500, pd.GetStatus())
	assert.Equal(t, "request failed", pd.GetDetail())
	assert.Equal(t, "", pd.GetInstance())
	assert.Nil(t, pd.GetErrors())
	assert.Equal(t, "", pd.GetRequestID())
	assert.Equal(t, "", pd.GetCorrelationID())
	assert.Equal(t, "", pd.GetTraceId())
	assert.Equal(t, wrapped, pd.Unwrap())
	assert.Equal(t, "request failed: database timeout", pd.Error())
}

func TestProblemDetails_Write_PopulatesFromContext(t *testing.T) {
	gin.SetMode(gin.TestMode)
	httperrors.BaseUrl = "https://example.com"

	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)
	ctx.Request = httptest.NewRequest(http.MethodGet, "/orders/42", nil)
	ctx.Set(constants.RequestIdKey, "req-123")
	ctx.Set(constants.CorrelationIdKey, "corr-456")
	ctx.Set(constants.TraceIdKey, "trace-789")

	pd := httperrors.NewNotFound("order not found")
	pd.Write(ctx)

	require.Equal(t, http.StatusNotFound, recorder.Code)
	assert.Equal(t, "req-123", recorder.Header().Get(constants.RequestIdHeader))
	assert.Equal(t, "corr-456", recorder.Header().Get(constants.CorrelationIdHeader))
	assert.Equal(t, "trace-789", recorder.Header().Get(constants.TraceIdHeader))

	var payload map[string]interface{}
	require.NoError(t, json.Unmarshal(recorder.Body.Bytes(), &payload))
	assert.Equal(t, "order not found", payload["detail"])
	assert.Equal(t, "/orders/42", payload["instance"])
	assert.Equal(t, "req-123", payload["requestId"])
	assert.Equal(t, "corr-456", payload["correlationId"])
	assert.Equal(t, "trace-789", payload["traceId"])
}

func TestProblemDetails_Write_DoesNothingWhenAlreadyWritten(t *testing.T) {
	gin.SetMode(gin.TestMode)

	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)
	ctx.Request = httptest.NewRequest(http.MethodGet, "/already-written", nil)
	ctx.Status(http.StatusNoContent)
	ctx.Writer.WriteHeaderNow()

	pd := httperrors.NewInternal("should not write")
	pd.Write(ctx)

	assert.Equal(t, http.StatusNoContent, recorder.Code)
	assert.Equal(t, "", recorder.Body.String())
}
