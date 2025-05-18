package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDataGet(t *testing.T) {
	tests := map[string]struct {
		data     Data
		key      string
		expected any
	}{
		"Empty data": {
			data:     Data{},
			key:      "key",
			expected: nil,
		},
		"Key exists": {
			data:     Data{"key": "value"},
			key:      "key",
			expected: "value",
		},
		"Key doesn't exist": {
			data:     Data{"other": "value"},
			key:      "key",
			expected: nil,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			result := tt.data.Get(tt.key)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestDataGetInt(t *testing.T) {
	tests := map[string]struct {
		data     Data
		key      string
		expected int
	}{
		"Empty data": {
			data:     Data{},
			key:      "key",
			expected: 0,
		},
		"Key doesn't exist": {
			data:     Data{"other": 42},
			key:      "key",
			expected: 0,
		},
		"Value is int": {
			data:     Data{"key": 42},
			key:      "key",
			expected: 42,
		},
		"Value is int64": {
			data:     Data{"key": int64(42)},
			key:      "key",
			expected: 42,
		},
		"Value is int32": {
			data:     Data{"key": int32(42)},
			key:      "key",
			expected: 42,
		},
		"Value is int16": {
			data:     Data{"key": int16(42)},
			key:      "key",
			expected: 42,
		},
		"Value is int8": {
			data:     Data{"key": int8(42)},
			key:      "key",
			expected: 42,
		},
		"Value is uint": {
			data:     Data{"key": uint(42)},
			key:      "key",
			expected: 42,
		},
		"Value is uint64": {
			data:     Data{"key": uint64(42)},
			key:      "key",
			expected: 42,
		},
		"Value is uint32": {
			data:     Data{"key": uint32(42)},
			key:      "key",
			expected: 42,
		},
		"Value is uint16": {
			data:     Data{"key": uint16(42)},
			key:      "key",
			expected: 42,
		},
		"Value is uint8": {
			data:     Data{"key": uint8(42)},
			key:      "key",
			expected: 42,
		},
		"Value is not numeric": {
			data:     Data{"key": "not a number"},
			key:      "key",
			expected: 0,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			result := tt.data.GetInt(tt.key)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestDataGetFloat(t *testing.T) {
	tests := map[string]struct {
		data     Data
		key      string
		expected float64
	}{
		"Empty data": {
			data:     Data{},
			key:      "key",
			expected: 0,
		},
		"Key doesn't exist": {
			data:     Data{"other": 42.5},
			key:      "key",
			expected: 0,
		},
		"Value is float64": {
			data:     Data{"key": 42.5},
			key:      "key",
			expected: 42.5,
		},
		"Value is float32": {
			data:     Data{"key": float32(42.5)},
			key:      "key",
			expected: 42.5,
		},
		"Value is int": {
			data:     Data{"key": 42},
			key:      "key",
			expected: 42.0,
		},
		"Value is int64": {
			data:     Data{"key": int64(42)},
			key:      "key",
			expected: 42.0,
		},
		"Value is int32": {
			data:     Data{"key": int32(42)},
			key:      "key",
			expected: 42.0,
		},
		"Value is int16": {
			data:     Data{"key": int16(42)},
			key:      "key",
			expected: 42.0,
		},
		"Value is int8": {
			data:     Data{"key": int8(42)},
			key:      "key",
			expected: 42.0,
		},
		"Value is uint": {
			data:     Data{"key": uint(42)},
			key:      "key",
			expected: 42.0,
		},
		"Value is uint64": {
			data:     Data{"key": uint64(42)},
			key:      "key",
			expected: 42.0,
		},
		"Value is uint32": {
			data:     Data{"key": uint32(42)},
			key:      "key",
			expected: 42.0,
		},
		"Value is uint16": {
			data:     Data{"key": uint16(42)},
			key:      "key",
			expected: 42.0,
		},
		"Value is uint8": {
			data:     Data{"key": uint8(42)},
			key:      "key",
			expected: 42.0,
		},
		"Value is not numeric": {
			data:     Data{"key": "not a number"},
			key:      "key",
			expected: 0,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			result := tt.data.GetFloat(tt.key)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestDataGetString(t *testing.T) {
	tests := map[string]struct {
		data     Data
		key      string
		expected string
	}{
		"Empty data": {
			data:     Data{},
			key:      "key",
			expected: "",
		},
		"Key doesn't exist": {
			data:     Data{"other": "value"},
			key:      "key",
			expected: "",
		},
		"Value is string": {
			data:     Data{"key": "value"},
			key:      "key",
			expected: "value",
		},
		"Value is not string": {
			data:     Data{"key": 42},
			key:      "key",
			expected: "",
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			result := tt.data.GetString(tt.key)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestDataGetBool(t *testing.T) {
	tests := map[string]struct {
		data     Data
		key      string
		expected bool
	}{
		"Empty data": {
			data:     Data{},
			key:      "key",
			expected: false,
		},
		"Key doesn't exist": {
			data:     Data{"other": true},
			key:      "key",
			expected: false,
		},
		"Value is true": {
			data:     Data{"key": true},
			key:      "key",
			expected: true,
		},
		"Value is false": {
			data:     Data{"key": false},
			key:      "key",
			expected: false,
		},
		"Value is not bool": {
			data:     Data{"key": "not a bool"},
			key:      "key",
			expected: false,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			result := tt.data.GetBool(tt.key)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestDataGetData(t *testing.T) {
	tests := map[string]struct {
		data     Data
		key      string
		expected Data
	}{
		"Empty data": {
			data:     Data{},
			key:      "key",
			expected: nil,
		},
		"Key doesn't exist": {
			data:     Data{"other": Data{"nested": "value"}},
			key:      "key",
			expected: nil,
		},
		"Value is Map": {
			data:     Data{"key": Data{"nested": "value"}},
			key:      "key",
			expected: Data{"nested": "value"},
		},
		"Value is not Map": {
			data:     Data{"key": "not a map"},
			key:      "key",
			expected: nil,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			result := tt.data.GetData(tt.key)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestDataSet(t *testing.T) {
	tests := map[string]struct {
		data     Data
		key      string
		value    any
		expected Data
	}{
		"Empty data": {
			data:     Data{},
			key:      "key",
			value:    "value",
			expected: Data{"key": "value"},
		},
		"Key doesn't exist": {
			data:     Data{"existing": "value"},
			key:      "key",
			value:    "new value",
			expected: Data{"existing": "value", "key": "new value"},
		},
		"Key exists": {
			data:     Data{"key": "old value"},
			key:      "key",
			value:    "new value",
			expected: Data{"key": "new value"},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			result := tt.data.Set(tt.key, tt.value)
			assert.Equal(t, tt.expected, result)
			// Also verify that the original data was modified
			assert.Equal(t, tt.expected, tt.data)
		})
	}
}

func TestDataHas(t *testing.T) {
	tests := map[string]struct {
		data     Data
		key      string
		expected bool
	}{
		"Empty data": {
			data:     Data{},
			key:      "key",
			expected: false,
		},
		"Key exists": {
			data:     Data{"key": "value"},
			key:      "key",
			expected: true,
		},
		"Key doesn't exist": {
			data:     Data{"other": "value"},
			key:      "key",
			expected: false,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			result := tt.data.Has(tt.key)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestDataDelete(t *testing.T) {
	tests := map[string]struct {
		data     Data
		key      string
		expected Data
	}{
		"Empty data": {
			data:     Data{},
			key:      "key",
			expected: Data{},
		},
		"Key exists": {
			data:     Data{"key": "value", "other": "value2"},
			key:      "key",
			expected: Data{"other": "value2"},
		},
		"Key doesn't exist": {
			data:     Data{"other": "value"},
			key:      "key",
			expected: Data{"other": "value"},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			result := tt.data.Delete(tt.key)
			assert.Equal(t, tt.expected, result)
			// Also verify that the original data was modified
			assert.Equal(t, tt.expected, tt.data)
		})
	}
}
