package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCleanPath(t *testing.T) {
	tests := map[string]struct {
		input    string
		expected string
	}{
		"No trailing dot": {
			input:    "/api/users",
			expected: "/api/users",
		},
		"With trailing dot": {
			input:    "/api/users.",
			expected: "/api/users",
		},
		"Multiple trailing dots": {
			input:    "/api/users...",
			expected: "/api/users..",
		},
		"Empty string": {
			input:    "",
			expected: "",
		},
		"Only dot": {
			input:    ".",
			expected: "",
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			result := cleanPath(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestValidatePath(t *testing.T) {
	tests := map[string]struct {
		input       string
		shouldPanic bool
	}{
		"Valid path": {
			input:       "/api/users",
			shouldPanic: false,
		},
		"Empty path": {
			input:       "",
			shouldPanic: true,
		},
		"Path without leading slash": {
			input:       "api/users",
			shouldPanic: true,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			if tt.shouldPanic {
				assert.Panics(t, func() {
					validatePath(tt.input)
				})
			} else {
				assert.NotPanics(t, func() {
					validatePath(tt.input)
				})
			}
		})
	}
}

func TestGetOptionalPaths(t *testing.T) {
	tests := map[string]struct {
		input    string
		expected []string
	}{
		"No optional parameters": {
			input:    "/api/users",
			expected: []string{},
		},
		"Simple optional parameter": {
			input:    "/api/users/{id?}",
			expected: []string{"/api/users", "/api/users/{id}"},
		},
		"Multiple optional parameters": {
			input:    "/api/{version?}/users/{id?}",
			expected: []string{"/api", "/api/{version}", "/api/{version}/users", "/api/{version}/users/{id}"},
		},
		"Optional parameter with regex": {
			input:    "/api/users/{id?:[0-9]+}",
			expected: []string{"/api/users", "/api/users/{id:[0-9]+}"},
		},
		"Root path with optional parameter": {
			input:    "/{param?}",
			expected: []string{"/", "/{param}"},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			result := getOptionalPaths(tt.input)
			assert.ElementsMatch(t, tt.expected, result, "Paths should match regardless of order")
		})
	}
}
