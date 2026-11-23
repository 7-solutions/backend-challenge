package http_test

import (
	"net/http"
	"testing"

	testHelpersHttp "github.com/chaihanij/backend-challenge/app/test_helpers/http"
	"github.com/stretchr/testify/require"
)

func TestHelpersHttp_GinContextJSON_ReturnsContextWithCorrectRequestBody(t *testing.T) {
	json := `{"foo":"bar"}`
	c := testHelpersHttp.GinContextJSON(json)
	require.Equal(t, http.MethodPost, c.Request.Method)
	require.Equal(t, "application/json", c.Request.Header.Get("Content-Type"))
	buf := make([]byte, c.Request.ContentLength)
	n, _ := c.Request.Body.Read(buf)
	require.Equal(t, len(json), n)
	require.Equal(t, json, string(buf))
}

func TestHelpersHttp_GinContextJSON_ReturnsContextWithRecorder(t *testing.T) {
	c := testHelpersHttp.GinContextJSON(`{}`)
	require.NotNil(t, c.Writer)
}

func TestHelpersHttp_GinContextJSON_HandlesEmptyJSONString(t *testing.T) {
	c := testHelpersHttp.GinContextJSON("")
	require.Equal(t, int64(0), c.Request.ContentLength)
}
