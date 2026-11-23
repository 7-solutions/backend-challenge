package constants

const (
	AuthorizationHeaderKey = "Authorization"
	RequestIdHeader        = "X-Request-ID"
	CorrelationIdHeader    = "X-Correlation-ID"
	TraceIdHeader          = "X-Trace-ID"
	RequestIdKey           = "requestId"
	CorrelationIdKey       = "correlationId"
	TraceIdKey             = "traceId"
	UserIdKey              = "userId"
	OTELTraceIdKey         = "otelTraceId"
	OTLSpanIdKey           = "otelSpanId"
	DDTraceIdKey           = "dd.trace_id"
	DDSpanIdKey            = "dd.span_id"
)

// ContextKey is an unexported type for context keys to prevent collisions
// with keys defined in other packages.
type ContextKey string

const (
	CtxKeyRequestID     ContextKey = RequestIdKey
	CtxKeyCorrelationID ContextKey = CorrelationIdKey
	CtxKeyTraceID       ContextKey = TraceIdKey
)
