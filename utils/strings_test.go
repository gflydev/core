package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnsafeBytes(t *testing.T) {
	tests := map[string]struct {
		input    string
		expected []byte
	}{
		"Empty string": {
			input:    "",
			expected: nil,
		},
		"Simple string": {
			input:    "hello",
			expected: []byte("hello"),
		},
		"String with special characters": {
			input:    "hello, 世界!",
			expected: []byte("hello, 世界!"),
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			result := UnsafeBytes(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestUnsafeStr(t *testing.T) {
	tests := map[string]struct {
		input    []byte
		expected string
	}{
		"Empty bytes": {
			input:    []byte{},
			expected: "",
		},
		"Simple bytes": {
			input:    []byte("hello"),
			expected: "hello",
		},
		"Bytes with special characters": {
			input:    []byte("hello, 世界!"),
			expected: "hello, 世界!",
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			result := UnsafeStr(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestIncludeStr(t *testing.T) {
	tests := map[string]struct {
		slice    []string
		search   string
		expected bool
	}{
		"Empty slice": {
			slice:    []string{},
			search:   "hello",
			expected: false,
		},
		"String found": {
			slice:    []string{"hello", "world", "test"},
			search:   "world",
			expected: true,
		},
		"String not found": {
			slice:    []string{"hello", "world", "test"},
			search:   "golang",
			expected: false,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			result := IncludeStr(tt.slice, tt.search)
			assert.Equal(t, tt.expected, result)
		})
	}
}
