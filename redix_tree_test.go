package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// MockHandler is a simple implementation of IHandler for testing
type MockHandler struct {
	Endpoint
	ID string // Used to identify the handler in tests
}

// TestTreeAdd tests the Add method of the Tree struct
func TestTreeAdd(t *testing.T) {
	tree := NewTree()

	// Test adding a simple path
	handler1 := &MockHandler{ID: "/api/users"}
	assert.NotPanics(t, func() {
		tree.Add("/api/users", handler1)
	}, "Adding a valid path should not panic")

	// Test adding a path with parameters
	handler2 := &MockHandler{ID: "/api/users/{id}"}
	assert.NotPanics(t, func() {
		tree.Add("/api/users/{id}", handler2)
	}, "Adding a path with parameters should not panic")

	// Test adding a path with optional parameters
	// Note: We need a new tree here because adding an optional parameter path
	// after adding a similar non-optional path would cause a conflict
	treeForOptional := NewTree()
	handler3 := &MockHandler{ID: "/api/users/{id?}"}
	assert.NotPanics(t, func() {
		treeForOptional.Add("/api/users/{id?}", handler3)
	}, "Adding a path with optional parameters should not panic")

	// Test that adding conflicting paths causes a panic
	assert.Panics(t, func() {
		tree.Add("/api/users/{id?}", handler3)
	}, "Adding a conflicting path should panic")

	// Test adding an invalid path (should panic)
	assert.Panics(t, func() {
		tree.Add("api/invalid", &MockHandler{})
	}, "Adding a path without leading slash should panic")

	// Test adding a nil handler (should panic)
	assert.Panics(t, func() {
		tree.Add("/api/nil", nil)
	}, "Adding a nil handler should panic")
}

// TestTreeMutable tests the Mutable property of the Tree struct
func TestTreeMutable(t *testing.T) {
	tree := NewTree()
	tree.Mutable = true

	// Add a handler
	handler1 := &MockHandler{ID: "original"}
	tree.Add("/api/users", handler1)

	// Replace the handler
	handler2 := &MockHandler{ID: "replacement"}
	assert.NotPanics(t, func() {
		tree.Add("/api/users", handler2)
	}, "Replacing a handler when Mutable is true should not panic")
}

// TestTreeAddMultiplePaths tests adding multiple paths to the tree
func TestTreeAddMultiplePaths(t *testing.T) {
	tree := NewTree()

	paths := []string{
		"/api/users",
		"/api/users/new",
		"/api/posts",
		"/api/comments",
		"/api/users/{id}",
		"/api/posts/{id}",
		"/api/users/{id}/comments",
		"/api/users/{id}/posts/{postId}",
		"/api/optional/{param?}",
	}

	// Add all paths to the tree
	for _, path := range paths {
		handler := &MockHandler{ID: path}
		assert.NotPanics(t, func() {
			tree.Add(path, handler)
		}, "Adding path %s should not panic", path)
	}
}

// TestLongestCommonPrefix verifies byte-based prefix matching, including the
// previously-broken multi-byte UTF-8 case where the rune-count bound made the
// loop exit early and under-report the prefix length.
func TestLongestCommonPrefix(t *testing.T) {
	tests := []struct {
		name string
		a, b string
		want int
	}{
		{"identical ascii", "/users", "/users", 6},
		{"partial ascii", "/users", "/us_rs", 3},
		{"no common", "abc", "xyz", 0},
		{"empty a", "", "abc", 0},
		{"empty b", "abc", "", 0},
		{"prefix shorter", "/user", "/users", 5},
		{"identical multibyte", "ééé", "ééé", len("ééé")},
		{"partial multibyte", "é-a", "é-b", len("é-")},
		{"multibyte then differ", "café", "cafx", 3},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, longestCommonPrefix(tt.a, tt.b))
		})
	}
}
