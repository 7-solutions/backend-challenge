package slogx

import (
	"context"
	"log/slog"

	"github.com/chaihanij/backend-challenge/app/constants"
	"go.opentelemetry.io/otel/trace"
)

type handler struct {
	next slog.Handler
}

// NewHandler wraps next with a handler that enriches every log record with:
//   - correlation ID, request ID, and trace ID extracted from the context
//   - OpenTelemetry trace/span IDs
//   - Datadog-compatible dd.trace_id / dd.span_id fields
func NewHandler(next slog.Handler) slog.Handler {
	return &handler{next: next}
}

func (h *handler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.next.Enabled(ctx, level)
}

func (h *handler) Handle(ctx context.Context, r slog.Record) error {
	// 1. Custom IDs from context.
	addStringAttr(&r, ctx, constants.CtxKeyCorrelationID, constants.CorrelationIdKey)
	addStringAttr(&r, ctx, constants.CtxKeyRequestID, constants.RequestIdKey)
	addStringAttr(&r, ctx, constants.CtxKeyTraceID, constants.TraceIdKey)

	// 2. OpenTelemetry span context.
	if spanCtx := trace.SpanContextFromContext(ctx); spanCtx.IsValid() {
		traceID := spanCtx.TraceID().String()
		spanID := spanCtx.SpanID().String()

		r.AddAttrs(
			slog.String(constants.OTELTraceIdKey, traceID),
			slog.String(constants.OTLSpanIdKey, spanID),
			// Datadog trace correlation fields.
			slog.String(constants.DDTraceIdKey, traceID),
			slog.String(constants.DDSpanIdKey, spanID),
		)
	}

	return h.next.Handle(ctx, r)
}

func (h *handler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &handler{next: h.next.WithAttrs(attrs)}
}

func (h *handler) WithGroup(name string) slog.Handler {
	return &handler{next: h.next.WithGroup(name)}
}

// addStringAttr reads a string value from ctx using ctxKey and, when non-empty,
// appends a slog.String attribute with fieldName to r.
func addStringAttr(r *slog.Record, ctx context.Context, key constants.ContextKey, fieldName string) {
	if v, _ := ctx.Value(key).(string); v != "" {
		r.AddAttrs(slog.String(fieldName, v))
	}
}
