package middlewares

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsPathIgnored(t *testing.T) {
	tests := []struct {
		name        string
		path        string
		ignorePaths []string
		expected    bool
	}{
		{
			name:        "empty ignore paths",
			path:        "/health",
			ignorePaths: []string{},
			expected:    false,
		},
		{
			name:        "exact match",
			path:        "/health",
			ignorePaths: []string{"/health"},
			expected:    true,
		},
		{
			name:        "no match",
			path:        "/api/users",
			ignorePaths: []string{"/health", "/metrics"},
			expected:    false,
		},
		{
			name:        "prefix match with *",
			path:        "/metrics/cpu",
			ignorePaths: []string{"/metrics*"},
			expected:    true,
		},
		{
			name:        "prefix match with * - not matching",
			path:        "/api/metrics",
			ignorePaths: []string{"/metrics*"},
			expected:    false,
		},
		{
			name:        "multiple paths - one matches",
			path:        "/ping",
			ignorePaths: []string{"/health", "/ping", "/ready"},
			expected:    true,
		},
		{
			name:        "multiple paths - none match",
			path:        "/api/users",
			ignorePaths: []string{"/health", "/ping", "/ready"},
			expected:    false,
		},
		{
			name:        "glob pattern *",
			path:        "/api/v1/internal/debug",
			ignorePaths: []string{"/api/*/internal/*"},
			expected:    true,
		},
		{
			name:        "glob pattern ? (single char)",
			path:        "/api/v1",
			ignorePaths: []string{"/api/v?"},
			expected:    true,
		},
		{
			name:        "case sensitive matching",
			path:        "/Health",
			ignorePaths: []string{"/health"},
			expected:    false,
		},
		{
			name:        "trailing slash difference",
			path:        "/health/",
			ignorePaths: []string{"/health"},
			expected:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsPathIgnored(tt.path, tt.ignorePaths)
			assert.Equal(t, tt.expected, result, "path: %s, ignorePaths: %v", tt.path, tt.ignorePaths)
		})
	}
}
