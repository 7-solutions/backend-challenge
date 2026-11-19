package slogx_test

import (
	"bytes"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/chaihanij/backend-challenge/app/util/slogx"
)

// attrsToMap converts []slog.Attr to map[string]any for easier assertions.
func attrsToMap(attrs []slog.Attr) map[string]any {
	m := make(map[string]any, len(attrs))
	for _, a := range attrs {
		m[a.Key] = a.Value.Any()
	}
	return m
}

func TestStructToAttrs_ValidStruct(t *testing.T) {
	type payload struct {
		Name  string `json:"name"`
		Score int    `json:"score"`
	}

	attrs := slogx.StructToAttrs(payload{Name: "alice", Score: 42})
	require.Len(t, attrs, 2)

	m := attrsToMap(attrs)
	assert.Equal(t, "alice", m["name"])
	assert.EqualValues(t, 42, m["score"])
}

func TestStructToAttrs_NestedStruct(t *testing.T) {
	type inner struct {
		Value string `json:"value"`
	}
	type outer struct {
		ID    int   `json:"id"`
		Inner inner `json:"inner"`
	}

	attrs := slogx.StructToAttrs(outer{ID: 1, Inner: inner{Value: "x"}})
	require.NotEmpty(t, attrs)

	m := attrsToMap(attrs)
	assert.EqualValues(t, 1, m["id"])

	nested, ok := m["inner"].(map[string]any)
	require.True(t, ok, "expected inner to be map, got %T", m["inner"])
	assert.Equal(t, "x", nested["value"])
}

func TestStructToAttrs_EmptyStruct(t *testing.T) {
	type empty struct{}
	attrs := slogx.StructToAttrs(empty{})
	assert.Empty(t, attrs)
}

func TestStructToAttrs_NilInput(t *testing.T) {
	// nil marshals to JSON "null"; json.Unmarshal sets the map to nil (no error),
	// so the result is an empty slice with no attributes.
	attrs := slogx.StructToAttrs(nil)
	assert.Empty(t, attrs)
}

func TestStructToAttrs_UnmarshalableInput(t *testing.T) {
	// Functions cannot be marshalled to JSON → marshal_error attribute expected.
	attrs := slogx.StructToAttrs(func() {})
	require.Len(t, attrs, 1)
	assert.Equal(t, "marshal_error", attrs[0].Key)
}

func TestStructToAttrs_Map(t *testing.T) {
	input := map[string]any{"foo": "bar", "num": 7}
	attrs := slogx.StructToAttrs(input)
	require.Len(t, attrs, 2)

	m := attrsToMap(attrs)
	assert.Equal(t, "bar", m["foo"])
	assert.EqualValues(t, 7, m["num"])
}

func TestAttrsToAny_ConvertsCorrectly(t *testing.T) {
	attrs := []slog.Attr{
		slog.String("a", "hello"),
		slog.Int("b", 2),
	}

	result := slogx.AttrsToAny(attrs)
	require.Len(t, result, 2)

	a0, ok := result[0].(slog.Attr)
	require.True(t, ok, "expected slog.Attr, got %T", result[0])
	assert.Equal(t, "a", a0.Key)
	assert.Equal(t, "hello", a0.Value.String())

	a1, ok := result[1].(slog.Attr)
	require.True(t, ok, "expected slog.Attr, got %T", result[1])
	assert.Equal(t, "b", a1.Key)
}

func TestAttrsToAny_EmptySlice(t *testing.T) {
	assert.Empty(t, slogx.AttrsToAny(nil))
	assert.Empty(t, slogx.AttrsToAny([]slog.Attr{}))
}

func TestStructToAny_EndToEnd(t *testing.T) {
	type item struct {
		Key string `json:"key"`
	}

	result := slogx.StructToAny(item{Key: "value"})
	require.Len(t, result, 1)

	attr, ok := result[0].(slog.Attr)
	require.True(t, ok)
	assert.Equal(t, "key", attr.Key)
	assert.Equal(t, "value", attr.Value.String())
}

func TestStructToAny_CanBePassedToLogger(t *testing.T) {
	type meta struct {
		Service string `json:"service"`
		Version string `json:"version"`
	}

	var buf bytes.Buffer
	logger := slogx.New(&slogx.Config{Level: slog.LevelInfo, Format: "json", Writer: &buf})
	logger.Info("startup", slogx.StructToAny(meta{Service: "api", Version: "1.0"})...)

	out := buf.String()
	assert.Contains(t, out, `"service":"api"`)
	assert.Contains(t, out, `"version":"1.0"`)
}
