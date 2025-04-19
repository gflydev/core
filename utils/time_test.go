package utils

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestFormatTime(t *testing.T) {
	// Create a fixed time for testing
	testTime := time.Date(2023, 5, 15, 14, 30, 45, 0, time.UTC)

	tests := map[string]struct {
		time     time.Time
		layout   string
		expected string
	}{
		"ISO date": {
			time:     testTime,
			layout:   "2006-01-02",
			expected: "2023-05-15",
		},
		"Full datetime": {
			time:     testTime,
			layout:   "2006-01-02 15:04:05",
			expected: "2023-05-15 14:30:45",
		},
		"Month and year": {
			time:     testTime,
			layout:   "January 2006",
			expected: "May 2023",
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			result := FormatTime(tt.time, tt.layout)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestParseTime(t *testing.T) {
	tests := map[string]struct {
		value    string
		layout   string
		expected time.Time
		hasError bool
	}{
		"Valid ISO date": {
			value:    "2023-05-15",
			layout:   "2006-01-02",
			expected: time.Date(2023, 5, 15, 0, 0, 0, 0, time.UTC),
			hasError: false,
		},
		"Valid datetime": {
			value:    "2023-05-15 14:30:45",
			layout:   "2006-01-02 15:04:05",
			expected: time.Date(2023, 5, 15, 14, 30, 45, 0, time.UTC),
			hasError: false,
		},
		"Invalid format": {
			value:    "not-a-date",
			layout:   "2006-01-02",
			expected: time.Time{},
			hasError: true,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			result, err := ParseTime(tt.value, tt.layout)
			if tt.hasError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

func TestNowAndUTC(t *testing.T) {
	// These tests are simple and just ensure the functions don't panic
	t.Run("Now returns current time", func(t *testing.T) {
		now := Now()
		assert.NotEqual(t, time.Time{}, now)
	})

	t.Run("UTC returns current UTC time", func(t *testing.T) {
		utc := UTC()
		assert.NotEqual(t, time.Time{}, utc)
		assert.Equal(t, time.UTC, utc.Location())
	})
}

func TestIsWeekendAndWeekday(t *testing.T) {
	tests := map[string]struct {
		time      time.Time
		isWeekend bool
		isWeekday bool
	}{
		"Monday": {
			time:      time.Date(2023, 5, 15, 0, 0, 0, 0, time.UTC), // Monday
			isWeekend: false,
			isWeekday: true,
		},
		"Saturday": {
			time:      time.Date(2023, 5, 20, 0, 0, 0, 0, time.UTC), // Saturday
			isWeekend: true,
			isWeekday: false,
		},
		"Sunday": {
			time:      time.Date(2023, 5, 21, 0, 0, 0, 0, time.UTC), // Sunday
			isWeekend: true,
			isWeekday: false,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, tt.isWeekend, IsWeekend(tt.time))
			assert.Equal(t, tt.isWeekday, IsWeekday(tt.time))
		})
	}
}

func TestStartAndEndOfDay(t *testing.T) {
	// Create a time in the middle of the day
	midDay := time.Date(2023, 5, 15, 14, 30, 45, 123456789, time.UTC)

	t.Run("StartOfDay", func(t *testing.T) {
		start := StartOfDay(midDay)
		expected := time.Date(2023, 5, 15, 0, 0, 0, 0, time.UTC)
		assert.Equal(t, expected, start)
	})

	t.Run("EndOfDay", func(t *testing.T) {
		end := EndOfDay(midDay)
		expected := time.Date(2023, 5, 15, 23, 59, 59, 999999999, time.UTC)
		assert.Equal(t, expected, end)
	})
}

func TestStartAndEndOfWeek(t *testing.T) {
	// Wednesday, May 17, 2023
	midWeek := time.Date(2023, 5, 17, 14, 30, 45, 0, time.UTC)

	t.Run("StartOfWeek", func(t *testing.T) {
		start := StartOfWeek(midWeek)
		// Monday, May 15, 2023
		expected := time.Date(2023, 5, 15, 0, 0, 0, 0, time.UTC)
		assert.Equal(t, expected, start)
	})

	t.Run("EndOfWeek", func(t *testing.T) {
		end := EndOfWeek(midWeek)
		// Sunday, May 21, 2023
		expected := time.Date(2023, 5, 21, 23, 59, 59, 999999999, time.UTC)
		assert.Equal(t, expected, end)
	})
}

func TestStartAndEndOfMonth(t *testing.T) {
	// May 17, 2023
	midMonth := time.Date(2023, 5, 17, 14, 30, 45, 0, time.UTC)

	t.Run("StartOfMonth", func(t *testing.T) {
		start := StartOfMonth(midMonth)
		// May 1, 2023
		expected := time.Date(2023, 5, 1, 0, 0, 0, 0, time.UTC)
		assert.Equal(t, expected, start)
	})

	t.Run("EndOfMonth", func(t *testing.T) {
		end := EndOfMonth(midMonth)
		// May 31, 2023
		expected := time.Date(2023, 5, 31, 23, 59, 59, 999999999, time.UTC)
		assert.Equal(t, expected, end)
	})
}

func TestDaysInMonth(t *testing.T) {
	tests := map[string]struct {
		time     time.Time
		expected int
	}{
		"January (31 days)": {
			time:     time.Date(2023, 1, 15, 0, 0, 0, 0, time.UTC),
			expected: 31,
		},
		"February in non-leap year (28 days)": {
			time:     time.Date(2023, 2, 15, 0, 0, 0, 0, time.UTC),
			expected: 28,
		},
		"February in leap year (29 days)": {
			time:     time.Date(2024, 2, 15, 0, 0, 0, 0, time.UTC),
			expected: 29,
		},
		"April (30 days)": {
			time:     time.Date(2023, 4, 15, 0, 0, 0, 0, time.UTC),
			expected: 30,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			result := DaysInMonth(tt.time)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestAddAndSubtractDays(t *testing.T) {
	// May 15, 2023
	baseTime := time.Date(2023, 5, 15, 14, 30, 45, 0, time.UTC)

	t.Run("AddDays", func(t *testing.T) {
		result := AddDays(baseTime, 5)
		expected := time.Date(2023, 5, 20, 14, 30, 45, 0, time.UTC)
		assert.Equal(t, expected, result)
	})

	t.Run("SubtractDays", func(t *testing.T) {
		result := SubtractDays(baseTime, 5)
		expected := time.Date(2023, 5, 10, 14, 30, 45, 0, time.UTC)
		assert.Equal(t, expected, result)
	})
}

func TestDaysBetween(t *testing.T) {
	tests := map[string]struct {
		t1       time.Time
		t2       time.Time
		expected int
	}{
		"Same day": {
			t1:       time.Date(2023, 5, 15, 10, 0, 0, 0, time.UTC),
			t2:       time.Date(2023, 5, 15, 14, 0, 0, 0, time.UTC),
			expected: 0,
		},
		"One day difference": {
			t1:       time.Date(2023, 5, 15, 10, 0, 0, 0, time.UTC),
			t2:       time.Date(2023, 5, 16, 14, 0, 0, 0, time.UTC),
			expected: 1,
		},
		"Negative difference": {
			t1:       time.Date(2023, 5, 16, 10, 0, 0, 0, time.UTC),
			t2:       time.Date(2023, 5, 15, 14, 0, 0, 0, time.UTC),
			expected: -1,
		},
		"Multiple days": {
			t1:       time.Date(2023, 5, 15, 10, 0, 0, 0, time.UTC),
			t2:       time.Date(2023, 5, 20, 14, 0, 0, 0, time.UTC),
			expected: 5,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			result := DaysBetween(tt.t1, tt.t2)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestIsSameDay(t *testing.T) {
	tests := map[string]struct {
		t1       time.Time
		t2       time.Time
		expected bool
	}{
		"Same day, different times": {
			t1:       time.Date(2023, 5, 15, 10, 0, 0, 0, time.UTC),
			t2:       time.Date(2023, 5, 15, 14, 0, 0, 0, time.UTC),
			expected: true,
		},
		"Different days": {
			t1:       time.Date(2023, 5, 15, 10, 0, 0, 0, time.UTC),
			t2:       time.Date(2023, 5, 16, 10, 0, 0, 0, time.UTC),
			expected: false,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			result := IsSameDay(tt.t1, tt.t2)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestIsSameMonth(t *testing.T) {
	tests := map[string]struct {
		t1       time.Time
		t2       time.Time
		expected bool
	}{
		"Same month, different days": {
			t1:       time.Date(2023, 5, 15, 10, 0, 0, 0, time.UTC),
			t2:       time.Date(2023, 5, 20, 14, 0, 0, 0, time.UTC),
			expected: true,
		},
		"Different months": {
			t1:       time.Date(2023, 5, 15, 10, 0, 0, 0, time.UTC),
			t2:       time.Date(2023, 6, 15, 10, 0, 0, 0, time.UTC),
			expected: false,
		},
		"Same month, different years": {
			t1:       time.Date(2023, 5, 15, 10, 0, 0, 0, time.UTC),
			t2:       time.Date(2024, 5, 15, 10, 0, 0, 0, time.UTC),
			expected: false,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			result := IsSameMonth(tt.t1, tt.t2)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestIsSameYear(t *testing.T) {
	tests := map[string]struct {
		t1       time.Time
		t2       time.Time
		expected bool
	}{
		"Same year, different months": {
			t1:       time.Date(2023, 5, 15, 10, 0, 0, 0, time.UTC),
			t2:       time.Date(2023, 6, 20, 14, 0, 0, 0, time.UTC),
			expected: true,
		},
		"Different years": {
			t1:       time.Date(2023, 5, 15, 10, 0, 0, 0, time.UTC),
			t2:       time.Date(2024, 5, 15, 10, 0, 0, 0, time.UTC),
			expected: false,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			result := IsSameYear(tt.t1, tt.t2)
			assert.Equal(t, tt.expected, result)
		})
	}
}
