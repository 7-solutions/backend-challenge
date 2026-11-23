// Package slogx extends the standard log/slog package with a handler that
// automatically enriches every log record with request-tracing IDs and
// distributed-tracing fields from OpenTelemetry and Datadog.
//
// # Handler decorator
//
// NewHandler wraps any existing slog.Handler (JSON or text) and adds the
// following attributes to every record before forwarding it:
//   - requestId, correlationId, traceId — extracted from context via the typed
//     ContextKey constants defined in the constants package
//   - otelTraceId, otelSpanId — from the active OpenTelemetry span context;
//     fields are omitted when no valid span is attached to the context
//   - dd.trace_id, dd.span_id — same values as the OTEL fields, formatted for
//     Datadog log correlation
//
// The convenience constructors New and SetDefault wrap NewHandler around a
// standard JSON or text handler according to Config.Format.
//
// # Context propagation
//
// Use the context helpers to carry IDs when code runs outside of a Gin handler
// (e.g. in background goroutines, service layers, or outgoing HTTP requests
// made with httpclient):
//
//	ctx = slogx.WithRequestID(ctx, reqID)
//	ctx = slogx.WithCorrelationID(ctx, corrID)
//	ctx = slogx.WithTraceID(ctx, traceID)
//
// # Struct logging helpers
//
// StructToAttrs marshals a struct to JSON then converts its top-level fields to
// []slog.Attr. AttrsToAny converts that slice to []any for use with the slog
// variadic logging methods. StructToAny combines both steps:
//
//	slog.InfoContext(ctx, "claims", slogx.StructToAny(claims)...)
package slogx
