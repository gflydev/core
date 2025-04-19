package utils

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestKeys(t *testing.T) {
	tests := map[string]struct {
		input    map[string]int
		expected []string
	}{
		"Empty map": {
			input:    map[string]int{},
			expected: []string{},
		},
		"Map with keys": {
			input: map[string]int{
				"a": 1,
				"b": 2,
				"c": 3,
			},
			expected: []string{"a", "b", "c"},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			result := Keys(tt.input)
			// Since map iteration order is not guaranteed, we need to check that all expected keys are present
			assert.ElementsMatch(t, tt.expected, result)
		})
	}
}

func TestValues(t *testing.T) {
	tests := map[string]struct {
		input    map[string]int
		expected []int
	}{
		"Empty map": {
			input:    map[string]int{},
			expected: []int{},
		},
		"Map with values": {
			input: map[string]int{
				"a": 1,
				"b": 2,
				"c": 3,
			},
			expected: []int{1, 2, 3},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			result := Values(tt.input)
			// Since map iteration order is not guaranteed, we need to check that all expected values are present
			assert.ElementsMatch(t, tt.expected, result)
		})
	}
}

func TestEntries(t *testing.T) {
	tests := map[string]struct {
		input map[string]int
	}{
		"Empty map": {
			input: map[string]int{},
		},
		"Map with entries": {
			input: map[string]int{
				"a": 1,
				"b": 2,
				"c": 3,
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			result := Entries(tt.input)
			assert.Equal(t, len(tt.input), len(result))

			// Create a map from the entries and verify it matches the original
			reconstructed := make(map[string]int)
			for _, entry := range result {
				reconstructed[entry.Key] = entry.Value
			}
			assert.Equal(t, tt.input, reconstructed)
		})
	}
}

func TestHasKey(t *testing.T) {
	testMap := map[string]int{
		"a": 1,
		"b": 2,
		"c": 3,
	}

	tests := map[string]struct {
		key      string
		expected bool
	}{
		"Key exists": {
			key:      "b",
			expected: true,
		},
		"Key does not exist": {
			key:      "d",
			expected: false,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			result := HasKey(testMap, tt.key)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestMerge(t *testing.T) {
	tests := map[string]struct {
		map1     map[string]int
		map2     map[string]int
		expected map[string]int
	}{
		"Both maps empty": {
			map1:     map[string]int{},
			map2:     map[string]int{},
			expected: map[string]int{},
		},
		"First map empty": {
			map1: map[string]int{},
			map2: map[string]int{
				"a": 1,
				"b": 2,
			},
			expected: map[string]int{
				"a": 1,
				"b": 2,
			},
		},
		"Second map empty": {
			map1: map[string]int{
				"a": 1,
				"b": 2,
			},
			map2: map[string]int{},
			expected: map[string]int{
				"a": 1,
				"b": 2,
			},
		},
		"No overlapping keys": {
			map1: map[string]int{
				"a": 1,
				"b": 2,
			},
			map2: map[string]int{
				"c": 3,
				"d": 4,
			},
			expected: map[string]int{
				"a": 1,
				"b": 2,
				"c": 3,
				"d": 4,
			},
		},
		"Overlapping keys": {
			map1: map[string]int{
				"a": 1,
				"b": 2,
				"c": 3,
			},
			map2: map[string]int{
				"b": 20,
				"c": 30,
				"d": 4,
			},
			expected: map[string]int{
				"a": 1,
				"b": 20,
				"c": 30,
				"d": 4,
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			result := Merge(tt.map1, tt.map2)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestPick(t *testing.T) {
	testMap := map[string]int{
		"a": 1,
		"b": 2,
		"c": 3,
		"d": 4,
	}

	tests := map[string]struct {
		keys     []string
		expected map[string]int
	}{
		"Empty keys": {
			keys:     []string{},
			expected: map[string]int{},
		},
		"All existing keys": {
			keys: []string{"a", "b", "c", "d"},
			expected: map[string]int{
				"a": 1,
				"b": 2,
				"c": 3,
				"d": 4,
			},
		},
		"Some existing keys": {
			keys: []string{"a", "c"},
			expected: map[string]int{
				"a": 1,
				"c": 3,
			},
		},
		"Non-existing keys": {
			keys:     []string{"e", "f"},
			expected: map[string]int{},
		},
		"Mix of existing and non-existing keys": {
			keys: []string{"a", "e", "c"},
			expected: map[string]int{
				"a": 1,
				"c": 3,
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			result := Pick(testMap, tt.keys)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestOmit(t *testing.T) {
	testMap := map[string]int{
		"a": 1,
		"b": 2,
		"c": 3,
		"d": 4,
	}

	tests := map[string]struct {
		keys     []string
		expected map[string]int
	}{
		"Empty keys": {
			keys: []string{},
			expected: map[string]int{
				"a": 1,
				"b": 2,
				"c": 3,
				"d": 4,
			},
		},
		"All existing keys": {
			keys:     []string{"a", "b", "c", "d"},
			expected: map[string]int{},
		},
		"Some existing keys": {
			keys: []string{"a", "c"},
			expected: map[string]int{
				"b": 2,
				"d": 4,
			},
		},
		"Non-existing keys": {
			keys: []string{"e", "f"},
			expected: map[string]int{
				"a": 1,
				"b": 2,
				"c": 3,
				"d": 4,
			},
		},
		"Mix of existing and non-existing keys": {
			keys: []string{"a", "e", "c"},
			expected: map[string]int{
				"b": 2,
				"d": 4,
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			result := Omit(testMap, tt.keys)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestMapKeys(t *testing.T) {
	testMap := map[string]int{
		"a": 1,
		"b": 2,
		"c": 3,
	}

	tests := map[string]struct {
		mapper   func(string) string
		expected map[string]int
	}{
		"Identity mapper": {
			mapper: func(s string) string { return s },
			expected: map[string]int{
				"a": 1,
				"b": 2,
				"c": 3,
			},
		},
		"Uppercase mapper": {
			mapper: strings.ToUpper,
			expected: map[string]int{
				"A": 1,
				"B": 2,
				"C": 3,
			},
		},
		"Prefix mapper": {
			mapper: func(s string) string { return "key_" + s },
			expected: map[string]int{
				"key_a": 1,
				"key_b": 2,
				"key_c": 3,
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			result := MapKeys(testMap, tt.mapper)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestMapValues(t *testing.T) {
	testMap := map[string]int{
		"a": 1,
		"b": 2,
		"c": 3,
	}

	tests := map[string]struct {
		mapper   func(int) int
		expected map[string]int
	}{
		"Identity mapper": {
			mapper: func(n int) int { return n },
			expected: map[string]int{
				"a": 1,
				"b": 2,
				"c": 3,
			},
		},
		"Double mapper": {
			mapper: func(n int) int { return n * 2 },
			expected: map[string]int{
				"a": 2,
				"b": 4,
				"c": 6,
			},
		},
		"Square mapper": {
			mapper: func(n int) int { return n * n },
			expected: map[string]int{
				"a": 1,
				"b": 4,
				"c": 9,
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			result := MapValues(testMap, tt.mapper)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestFilterKeys(t *testing.T) {
	testMap := map[string]int{
		"a": 1,
		"b": 2,
		"c": 3,
		"d": 4,
	}

	tests := map[string]struct {
		predicate func(string) bool
		expected  map[string]int
	}{
		"All keys": {
			predicate: func(s string) bool { return true },
			expected: map[string]int{
				"a": 1,
				"b": 2,
				"c": 3,
				"d": 4,
			},
		},
		"No keys": {
			predicate: func(s string) bool { return false },
			expected:  map[string]int{},
		},
		"Keys after 'b'": {
			predicate: func(s string) bool { return s > "b" },
			expected: map[string]int{
				"c": 3,
				"d": 4,
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			result := FilterKeys(testMap, tt.predicate)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestFilterValues(t *testing.T) {
	testMap := map[string]int{
		"a": 1,
		"b": 2,
		"c": 3,
		"d": 4,
	}

	tests := map[string]struct {
		predicate func(int) bool
		expected  map[string]int
	}{
		"All values": {
			predicate: func(n int) bool { return true },
			expected: map[string]int{
				"a": 1,
				"b": 2,
				"c": 3,
				"d": 4,
			},
		},
		"No values": {
			predicate: func(n int) bool { return false },
			expected:  map[string]int{},
		},
		"Even values": {
			predicate: func(n int) bool { return n%2 == 0 },
			expected: map[string]int{
				"b": 2,
				"d": 4,
			},
		},
		"Values greater than 2": {
			predicate: func(n int) bool { return n > 2 },
			expected: map[string]int{
				"c": 3,
				"d": 4,
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			result := FilterValues(testMap, tt.predicate)
			assert.Equal(t, tt.expected, result)
		})
	}
}
