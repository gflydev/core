package utils

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestContains(t *testing.T) {
	tests := map[string]struct {
		slice    []int
		element  int
		expected bool
	}{
		"Empty slice": {
			slice:    []int{},
			element:  5,
			expected: false,
		},
		"Element present": {
			slice:    []int{1, 2, 3, 4, 5},
			element:  3,
			expected: true,
		},
		"Element not present": {
			slice:    []int{1, 2, 3, 4, 5},
			element:  6,
			expected: false,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			result := Contains(tt.slice, tt.element)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestIndexOf(t *testing.T) {
	tests := map[string]struct {
		slice    []string
		element  string
		expected int
	}{
		"Empty slice": {
			slice:    []string{},
			element:  "test",
			expected: -1,
		},
		"Element present": {
			slice:    []string{"apple", "banana", "cherry"},
			element:  "banana",
			expected: 1,
		},
		"Element not present": {
			slice:    []string{"apple", "banana", "cherry"},
			element:  "grape",
			expected: -1,
		},
		"First element": {
			slice:    []string{"apple", "banana", "cherry"},
			element:  "apple",
			expected: 0,
		},
		"Last element": {
			slice:    []string{"apple", "banana", "cherry"},
			element:  "cherry",
			expected: 2,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			result := IndexOf(tt.slice, tt.element)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestFilter(t *testing.T) {
	tests := map[string]struct {
		slice     []int
		predicate func(int) bool
		expected  []int
	}{
		"Empty slice": {
			slice:     []int{},
			predicate: func(n int) bool { return n > 3 },
			expected:  []int{},
		},
		"No matches": {
			slice:     []int{1, 2, 3},
			predicate: func(n int) bool { return n > 3 },
			expected:  []int{},
		},
		"Some matches": {
			slice:     []int{1, 2, 3, 4, 5},
			predicate: func(n int) bool { return n > 3 },
			expected:  []int{4, 5},
		},
		"All matches": {
			slice:     []int{4, 5, 6},
			predicate: func(n int) bool { return n > 3 },
			expected:  []int{4, 5, 6},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			result := Filter(tt.slice, tt.predicate)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestMap(t *testing.T) {
	tests := map[string]struct {
		slice    []int
		mapper   func(int) string
		expected []string
	}{
		"Empty slice": {
			slice:    []int{},
			mapper:   func(n int) string { return strconv.Itoa(n) },
			expected: []string{},
		},
		"Transform integers to strings": {
			slice:    []int{1, 2, 3},
			mapper:   func(n int) string { return strconv.Itoa(n) },
			expected: []string{"1", "2", "3"},
		},
		"Double integers": {
			slice:    []int{1, 2, 3},
			mapper:   func(n int) string { return strconv.Itoa(n * 2) },
			expected: []string{"2", "4", "6"},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			result := Map(tt.slice, tt.mapper)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestReduce(t *testing.T) {
	tests := map[string]struct {
		slice        []int
		initialValue int
		reducer      func(int, int) int
		expected     int
	}{
		"Empty slice": {
			slice:        []int{},
			initialValue: 0,
			reducer:      func(acc, n int) int { return acc + n },
			expected:     0,
		},
		"Sum": {
			slice:        []int{1, 2, 3, 4, 5},
			initialValue: 0,
			reducer:      func(acc, n int) int { return acc + n },
			expected:     15,
		},
		"Product": {
			slice:        []int{1, 2, 3, 4, 5},
			initialValue: 1,
			reducer:      func(acc, n int) int { return acc * n },
			expected:     120,
		},
		"Max": {
			slice:        []int{3, 1, 4, 1, 5, 9, 2, 6},
			initialValue: -1,
			reducer: func(acc, n int) int {
				if n > acc {
					return n
				} else {
					return acc
				}
			},
			expected: 9,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			result := Reduce(tt.slice, tt.initialValue, tt.reducer)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestAny(t *testing.T) {
	tests := map[string]struct {
		slice     []int
		predicate func(int) bool
		expected  bool
	}{
		"Empty slice": {
			slice:     []int{},
			predicate: func(n int) bool { return n > 3 },
			expected:  false,
		},
		"No matches": {
			slice:     []int{1, 2, 3},
			predicate: func(n int) bool { return n > 3 },
			expected:  false,
		},
		"Some matches": {
			slice:     []int{1, 2, 3, 4, 5},
			predicate: func(n int) bool { return n > 3 },
			expected:  true,
		},
		"All matches": {
			slice:     []int{4, 5, 6},
			predicate: func(n int) bool { return n > 3 },
			expected:  true,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			result := Any(tt.slice, tt.predicate)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestAll(t *testing.T) {
	tests := map[string]struct {
		slice     []int
		predicate func(int) bool
		expected  bool
	}{
		"Empty slice": {
			slice:     []int{},
			predicate: func(n int) bool { return n > 3 },
			expected:  true,
		},
		"No matches": {
			slice:     []int{1, 2, 3},
			predicate: func(n int) bool { return n > 3 },
			expected:  false,
		},
		"Some matches": {
			slice:     []int{1, 2, 3, 4, 5},
			predicate: func(n int) bool { return n > 3 },
			expected:  false,
		},
		"All matches": {
			slice:     []int{4, 5, 6},
			predicate: func(n int) bool { return n > 3 },
			expected:  true,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			result := All(tt.slice, tt.predicate)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestUnique(t *testing.T) {
	tests := map[string]struct {
		slice    []int
		expected []int
	}{
		"Empty slice": {
			slice:    []int{},
			expected: []int{},
		},
		"No duplicates": {
			slice:    []int{1, 2, 3, 4, 5},
			expected: []int{1, 2, 3, 4, 5},
		},
		"With duplicates": {
			slice:    []int{1, 2, 2, 3, 3, 3, 4, 5, 5},
			expected: []int{1, 2, 3, 4, 5},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			result := Unique(tt.slice)
			assert.ElementsMatch(t, tt.expected, result)
		})
	}
}

func TestChunk(t *testing.T) {
	tests := map[string]struct {
		slice    []int
		size     int
		expected [][]int
	}{
		"Empty slice": {
			slice:    []int{},
			size:     2,
			expected: [][]int{},
		},
		"Size zero or negative": {
			slice:    []int{1, 2, 3, 4, 5},
			size:     0,
			expected: [][]int{},
		},
		"Even chunks": {
			slice:    []int{1, 2, 3, 4, 5, 6},
			size:     2,
			expected: [][]int{{1, 2}, {3, 4}, {5, 6}},
		},
		"Uneven chunks": {
			slice:    []int{1, 2, 3, 4, 5},
			size:     2,
			expected: [][]int{{1, 2}, {3, 4}, {5}},
		},
		"Size larger than slice": {
			slice:    []int{1, 2, 3},
			size:     5,
			expected: [][]int{{1, 2, 3}},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			result := Chunk(tt.slice, tt.size)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestReverse(t *testing.T) {
	tests := map[string]struct {
		slice    []int
		expected []int
	}{
		"Empty slice": {
			slice:    []int{},
			expected: []int{},
		},
		"Single element": {
			slice:    []int{1},
			expected: []int{1},
		},
		"Multiple elements": {
			slice:    []int{1, 2, 3, 4, 5},
			expected: []int{5, 4, 3, 2, 1},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			result := Reverse(tt.slice)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestShuffle(t *testing.T) {
	// Shuffle is non-deterministic, so we can only test that:
	// 1. The result has the same length as the input
	// 2. The result contains the same elements as the input
	// 3. The result is not the same as the input (for sufficiently large inputs)

	t.Run("Empty slice", func(t *testing.T) {
		slice := []int{}
		result := Shuffle(slice)
		assert.Equal(t, len(slice), len(result))
		assert.ElementsMatch(t, slice, result)
	})

	t.Run("Single element", func(t *testing.T) {
		slice := []int{1}
		result := Shuffle(slice)
		assert.Equal(t, len(slice), len(result))
		assert.ElementsMatch(t, slice, result)
	})

	t.Run("Multiple elements", func(t *testing.T) {
		slice := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
		result := Shuffle(slice)
		assert.Equal(t, len(slice), len(result))
		assert.ElementsMatch(t, slice, result)

		// This test could occasionally fail if the shuffle happens to return the original order,
		// but the probability is very low for a sufficiently large slice
		different := false
		for i := 0; i < len(slice); i++ {
			if slice[i] != result[i] {
				different = true
				break
			}
		}
		assert.True(t, different, "Shuffled slice should be different from original")
	})
}
