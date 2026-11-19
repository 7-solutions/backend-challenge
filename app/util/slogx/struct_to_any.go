package slogx

import (
	"encoding/json"
	"log/slog"
)

// StructToAttrs marshals v to JSON and returns its top-level fields as a
// slice of [slog.Attr]. Nested objects and arrays are preserved as-is.
// On marshal/unmarshal failure a single error attribute is returned instead.
func StructToAttrs(v any) []slog.Attr {
	b, err := json.Marshal(v)
	if err != nil {
		return []slog.Attr{slog.String("marshal_error", err.Error())}
	}

	var m map[string]any
	if err := json.Unmarshal(b, &m); err != nil {
		return []slog.Attr{slog.String("unmarshal_error", err.Error())}
	}

	attrs := make([]slog.Attr, 0, len(m))
	for k, v := range m {
		attrs = append(attrs, slog.Any(k, v))
	}
	return attrs
}

// AttrsToAny converts a []slog.Attr slice to []any so it can be passed
// directly to slog variadic logging methods.
func AttrsToAny(attrs []slog.Attr) []any {
	out := make([]any, len(attrs))
	for i, a := range attrs {
		out[i] = a
	}
	return out
}

// StructToAny is a convenience wrapper: StructToAttrs followed by AttrsToAny.
func StructToAny(v any) []any {
	return AttrsToAny(StructToAttrs(v))
}
