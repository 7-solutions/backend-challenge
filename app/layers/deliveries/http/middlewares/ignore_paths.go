package middlewares

import (
	"path/filepath"
	"strings"
)

// IsPathIgnored checks if a path should be ignored by comparing it against ignorePaths array.
// ignorePaths supports:
// - Exact matches: "/health"
// - Prefix matches: "/metrics*"
// - Wildcard patterns: "/api/*/debug"
// - Returns false if ignorePaths is empty
func IsPathIgnored(path string, ignorePaths []string) bool {
	if len(ignorePaths) == 0 {
		return false
	}

	for _, ignorePath := range ignorePaths {
		if matchPath(path, ignorePath) {
			return true
		}
	}

	return false
}

// matchPath checks if a path matches a pattern
func matchPath(path, pattern string) bool {
	// Exact match
	if path == pattern {
		return true
	}

	// Prefix match (pattern ends with *)
	if strings.HasSuffix(pattern, "*") {
		prefix := strings.TrimSuffix(pattern, "*")
		if strings.HasPrefix(path, prefix) {
			return true
		}
	}

	// Glob pattern match
	matched, err := filepath.Match(pattern, path)
	if err == nil && matched {
		return true
	}

	return false
}
