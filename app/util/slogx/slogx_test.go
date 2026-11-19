package slogx_test

import (
	"bytes"
	"context"
	"log/slog"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/chaihanij/backend-challenge/app/constants"
	"github.com/chaihanij/backend-challenge/app/util/slogx"
)

// newTestLogger creates a JSON logger that writes to buf, making tests
// independent of os.Stdout.
func newTestLogger(t *testing.T, buf *bytes.Buffer, level slog.Level) *slog.Logger {
	t.Helper()
	return slogx.New(&slogx.Config{Level: level, Format: "json", Writer: buf})
}

func TestConstants_StablePublicKeys(t *testing.T) {
	tests := []struct {
		name     string
		actual   string
		expected string
	}{
		{name: "request id", actual: constants.RequestIdKey, expected: "requestId"},
		{name: "correlation id", actual: constants.CorrelationIdKey, expected: "correlationId"},
		{name: "trace id", actual: constants.TraceIdKey, expected: "traceId"},
		{name: "otel trace id", actual: constants.OTELTraceIdKey, expected: "otelTraceId"},
		{name: "otel span id", actual: constants.OTLSpanIdKey, expected: "otelSpanId"},
		{name: "dd trace id", actual: constants.DDTraceIdKey, expected: "dd.trace_id"},
		{name: "dd span id", actual: constants.DDSpanIdKey, expected: "dd.span_id"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expected, tc.actual)
		})
	}
}

func TestNew_UsesExpectedOutputFormat(t *testing.T) {
	tests := []struct {
		name          string
		format        string
		expectJSON    bool
		expectedToken string
	}{
		{name: "json", format: "json", expectJSON: true, expectedToken: `"msg":"hello"`},
		{name: "text", format: "text", expectJSON: false, expectedToken: `msg=hello`},
		{name: "default text for unknown format", format: "unknown", expectJSON: false, expectedToken: `msg=hello`},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var buf bytes.Buffer
			logger := slogx.New(&slogx.Config{Level: slog.LevelInfo, Format: tc.format, Writer: &buf})
			logger.Info("hello")

			trimmed := strings.TrimSpace(buf.String())
			assert.Equal(t, tc.expectJSON, strings.HasPrefix(trimmed, "{"), "unexpected format: %q", trimmed)
			assert.Contains(t, trimmed, tc.expectedToken)
		})
	}
}

func TestNew_NilWriterDefaultsToStdout(t *testing.T) {
	// Just ensure no panic when Writer is nil.
	assert.NotPanics(t, func() {
		logger := slogx.New(&slogx.Config{Level: slog.LevelError, Format: "json"})
		_ = logger
	})
}

func TestNew_RespectsConfiguredLevel(t *testing.T) {
	var buf bytes.Buffer
	logger := newTestLogger(t, &buf, slog.LevelError)

	logger.Info("should be filtered")
	logger.Error("should be logged")

	out := buf.String()
	assert.NotContains(t, out, "should be filtered")
	assert.Contains(t, out, "should be logged")
}

func TestSetDefault_SetsGlobalDefaultLogger(t *testing.T) {
	old := slog.Default()
	t.Cleanup(func() { slog.SetDefault(old) })

	var buf bytes.Buffer
	got := slogx.SetDefault(&slogx.Config{Level: slog.LevelInfo, Format: "json", Writer: &buf})

	assert.Same(t, got, slog.Default())
}

// ---- Context helper tests ----

func TestWithRequestID_StoresAndRetrievesValue(t *testing.T) {
	ctx := slogx.WithRequestID(context.Background(), "req-42")

	var buf bytes.Buffer
	logger := newTestLogger(t, &buf, slog.LevelInfo)
	logger.InfoContext(ctx, "msg")

	assert.Contains(t, buf.String(), `"requestId":"req-42"`)
}

func TestWithCorrelationID_StoresAndRetrievesValue(t *testing.T) {
	ctx := slogx.WithCorrelationID(context.Background(), "corr-99")

	var buf bytes.Buffer
	logger := newTestLogger(t, &buf, slog.LevelInfo)
	logger.InfoContext(ctx, "msg")

	assert.Contains(t, buf.String(), `"correlationId":"corr-99"`)
}

func TestWithTraceID_StoresAndRetrievesValue(t *testing.T) {
	ctx := slogx.WithTraceID(context.Background(), "trace-77")

	var buf bytes.Buffer
	logger := newTestLogger(t, &buf, slog.LevelInfo)
	logger.InfoContext(ctx, "msg")

	assert.Contains(t, buf.String(), `"traceId":"trace-77"`)
}

func TestWithIDs_EmptyIDsOmittedFromOutput(t *testing.T) {
	ctx := slogx.WithRequestID(context.Background(), "")
	ctx = slogx.WithCorrelationID(ctx, "")
	ctx = slogx.WithTraceID(ctx, "")

	var buf bytes.Buffer
	logger := newTestLogger(t, &buf, slog.LevelInfo)
	logger.InfoContext(ctx, "msg")

	out := buf.String()
	assert.NotContains(t, out, "requestId")
	assert.NotContains(t, out, "correlationId")
	assert.NotContains(t, out, "traceId")
}
