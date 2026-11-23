package slogx_test

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/trace"

	"github.com/chaihanij/backend-challenge/app/constants"
	"github.com/chaihanij/backend-challenge/app/util/slogx"
)

// stubHandler records calls for delegation assertions.
type stubHandler struct {
	enabledReturn bool
	enabledCalled bool
	handleCalled  bool
	lastLevel     slog.Level
}

func (s *stubHandler) Enabled(_ context.Context, level slog.Level) bool {
	s.enabledCalled = true
	s.lastLevel = level
	return s.enabledReturn
}

func (s *stubHandler) Handle(_ context.Context, _ slog.Record) error {
	s.handleCalled = true
	return nil
}

func (s *stubHandler) WithAttrs(_ []slog.Attr) slog.Handler { return s }
func (s *stubHandler) WithGroup(_ string) slog.Handler      { return s }

// decodeSingleJSONLine decodes one JSON log line from buf.
func decodeSingleJSONLine(t *testing.T, b *bytes.Buffer) map[string]any {
	t.Helper()

	line := bytes.TrimSpace(b.Bytes())
	require.NotEmpty(t, line, "expected log output, got empty buffer")

	var out map[string]any
	err := json.Unmarshal(line, &out)
	require.NoError(t, err, "line=%s", string(line))

	return out
}

// newJSONLogger returns a logger backed by a JSON handler writing to buf.
func newJSONLogger(buf *bytes.Buffer) *slog.Logger {
	return slog.New(slogx.NewHandler(slog.NewJSONHandler(buf, nil)))
}

func TestNewHandler_EnabledDelegatesToNextHandler(t *testing.T) {
	base := &stubHandler{enabledReturn: true}
	h := slogx.NewHandler(base)

	assert.True(t, h.Enabled(context.Background(), slog.LevelWarn))
	assert.True(t, base.enabledCalled)
	assert.Equal(t, slog.LevelWarn, base.lastLevel)
}

func TestNewHandler_HandleDelegatesToNextHandler(t *testing.T) {
	base := &stubHandler{enabledReturn: true}
	h := slogx.NewHandler(base)

	r := slog.NewRecord(time.Now(), slog.LevelInfo, "delegated", 0)
	err := h.Handle(context.Background(), r)
	require.NoError(t, err)
	assert.True(t, base.handleCalled)
}

func TestNewHandler_HandleContextIDs(t *testing.T) {
	tests := []struct {
		name     string
		ctx      context.Context
		assertFn func(t *testing.T, entry map[string]any)
	}{
		{
			name: "adds all non-empty context ids",
			ctx: func() context.Context {
				ctx := context.Background()
				ctx = slogx.WithCorrelationID(ctx, "cid-123")
				ctx = slogx.WithRequestID(ctx, "rid-123")
				ctx = slogx.WithTraceID(ctx, "tid-123")
				return ctx
			}(),
			assertFn: func(t *testing.T, entry map[string]any) {
				assert.Equal(t, "cid-123", entry[constants.CorrelationIdKey])
				assert.Equal(t, "rid-123", entry[constants.RequestIdKey])
				assert.Equal(t, "tid-123", entry[constants.TraceIdKey])
			},
		},
		{
			name: "skips empty context ids",
			ctx: func() context.Context {
				ctx := context.Background()
				ctx = slogx.WithCorrelationID(ctx, "")
				ctx = slogx.WithRequestID(ctx, "")
				ctx = slogx.WithTraceID(ctx, "")
				return ctx
			}(),
			assertFn: func(t *testing.T, entry map[string]any) {
				assert.NotContains(t, entry, constants.CorrelationIdKey)
				assert.NotContains(t, entry, constants.RequestIdKey)
				assert.NotContains(t, entry, constants.TraceIdKey)
			},
		},
		{
			name: "no ids in context produces no id fields",
			ctx:  context.Background(),
			assertFn: func(t *testing.T, entry map[string]any) {
				assert.NotContains(t, entry, constants.CorrelationIdKey)
				assert.NotContains(t, entry, constants.RequestIdKey)
				assert.NotContains(t, entry, constants.TraceIdKey)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var buf bytes.Buffer
			logger := newJSONLogger(&buf)

			logger.InfoContext(tc.ctx, "hello")
			entry := decodeSingleJSONLine(t, &buf)
			tc.assertFn(t, entry)
		})
	}
}

func TestNewHandler_HandleAddsOpenTelemetryAndDatadogIDs(t *testing.T) {
	var buf bytes.Buffer
	logger := newJSONLogger(&buf)

	traceID := trace.TraceID{0x10, 0x20, 0x30, 0x40, 0x50, 0x60, 0x70, 0x80, 0x90, 0xA0, 0xB0, 0xC0, 0xD0, 0xE0, 0xF0, 0x01}
	spanID := trace.SpanID{0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88}
	spanCtx := trace.NewSpanContext(trace.SpanContextConfig{TraceID: traceID, SpanID: spanID, TraceFlags: trace.FlagsSampled})
	ctx := trace.ContextWithSpanContext(context.Background(), spanCtx)

	logger.InfoContext(ctx, "hello")
	entry := decodeSingleJSONLine(t, &buf)

	assert.Equal(t, traceID.String(), entry[constants.OTELTraceIdKey])
	assert.Equal(t, spanID.String(), entry[constants.OTLSpanIdKey])
	assert.Equal(t, traceID.String(), entry[constants.DDTraceIdKey])
	assert.Equal(t, spanID.String(), entry[constants.DDSpanIdKey])
}

func TestNewHandler_InvalidSpanContextProducesNoOTELFields(t *testing.T) {
	var buf bytes.Buffer
	logger := newJSONLogger(&buf)

	// context.Background() has no span – SpanContextFromContext returns invalid.
	logger.InfoContext(context.Background(), "hello")
	entry := decodeSingleJSONLine(t, &buf)

	assert.NotContains(t, entry, constants.OTELTraceIdKey)
	assert.NotContains(t, entry, constants.OTLSpanIdKey)
	assert.NotContains(t, entry, constants.DDTraceIdKey)
	assert.NotContains(t, entry, constants.DDSpanIdKey)
}

func TestNewHandler_WithAttrsAndWithGroupPreservesBehavior(t *testing.T) {
	var buf bytes.Buffer

	base := slog.NewJSONHandler(&buf, nil)
	h := slogx.NewHandler(base).WithGroup("meta").WithAttrs([]slog.Attr{slog.String("service", "payments")})
	logger := slog.New(h)

	logger.InfoContext(context.Background(), "hello")
	entry := decodeSingleJSONLine(t, &buf)

	meta, ok := entry["meta"].(map[string]any)
	require.True(t, ok, "expected grouped attrs under 'meta', got %#v", entry["meta"])
	assert.Equal(t, "payments", meta["service"])
}

func TestNewHandler_AllContextIDsAndOTELTogether(t *testing.T) {
	var buf bytes.Buffer
	logger := newJSONLogger(&buf)

	traceID := trace.TraceID{0x01}
	spanID := trace.SpanID{0x02}
	spanCtx := trace.NewSpanContext(trace.SpanContextConfig{TraceID: traceID, SpanID: spanID, TraceFlags: trace.FlagsSampled})

	ctx := trace.ContextWithSpanContext(context.Background(), spanCtx)
	ctx = slogx.WithRequestID(ctx, "req-1")
	ctx = slogx.WithCorrelationID(ctx, "corr-1")
	ctx = slogx.WithTraceID(ctx, "trace-1")

	logger.InfoContext(ctx, "combined")
	entry := decodeSingleJSONLine(t, &buf)

	assert.Equal(t, "req-1", entry[constants.RequestIdKey])
	assert.Equal(t, "corr-1", entry[constants.CorrelationIdKey])
	assert.Equal(t, "trace-1", entry[constants.TraceIdKey])
	assert.Equal(t, traceID.String(), entry[constants.OTELTraceIdKey])
	assert.Equal(t, spanID.String(), entry[constants.OTLSpanIdKey])
}
