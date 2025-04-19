package utils

import (
	"crypto/rand"
	"errors"
	"math"
	"math/big"
	"strconv"
)

// RandInt64 generates a cryptographically secure random integer within the range [0, m).
//
// Parameters:
//   - m int64: The upper limit (exclusive) for the random integer.
//
// Returns:
//   - int64: A random integer within the specified range.
//
// NOTE:
//
//	Get error `G404 (CWE-338): Use of weak random number generator (math/rand instead of crypto/rand)
//	(Confidence: MEDIUM, Severity: HIGH)` when use `rand.Intn(max)` from "math/rand".
//	Fixed Refer https://github.com/securego/gosec/issues/294#issuecomment-487452731
func RandInt64(m int64) int64 {
	// Get random Int64
	n, _ := rand.Int(rand.Reader, big.NewInt(m))

	return n.Int64()
}

// InRange checks if a number is within the specified range (inclusive).
//
// Parameters:
//   - value int: The value to check.
//   - min int: The minimum value of the range.
//   - max int: The maximum value of the range.
//
// Returns:
//   - bool: True if the value is within the range, false otherwise.
func InRange(value, min, max int) bool {
	return value >= min && value <= max
}

// Clamp constrains a value to a specified range.
//
// Parameters:
//   - value int: The value to constrain.
//   - min int: The minimum value of the range.
//   - max int: The maximum value of the range.
//
// Returns:
//   - int: The constrained value.
func Clamp(value, min, max int) int {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}

// Min returns the smaller of two integers.
//
// Parameters:
//   - a int: The first integer.
//   - b int: The second integer.
//
// Returns:
//   - int: The smaller of the two integers.
func Min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// Max returns the larger of two integers.
//
// Parameters:
//   - a int: The first integer.
//   - b int: The second integer.
//
// Returns:
//   - int: The larger of the two integers.
func Max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// SafeAtoi converts a string to an integer with error handling.
//
// Parameters:
//   - s string: The string to convert.
//   - defaultValue int: The default value to return if conversion fails.
//
// Returns:
//   - int: The converted integer or the default value if conversion fails.
func SafeAtoi(s string, defaultValue int) int {
	if s == "" {
		return defaultValue
	}

	value, err := strconv.Atoi(s)
	if err != nil {
		return defaultValue
	}

	return value
}

// SafeParseFloat converts a string to a float64 with error handling.
//
// Parameters:
//   - s string: The string to convert.
//   - defaultValue float64: The default value to return if conversion fails.
//
// Returns:
//   - float64: The converted float64 or the default value if conversion fails.
func SafeParseFloat(s string, defaultValue float64) float64 {
	if s == "" {
		return defaultValue
	}

	value, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return defaultValue
	}

	return value
}

// SafeInt32ToInt64 safely converts an int32 to int64.
//
// Parameters:
//   - value int32: The int32 value to convert.
//
// Returns:
//   - int64: The converted int64 value.
func SafeInt32ToInt64(value int32) int64 {
	return int64(value)
}

// SafeInt64ToInt32 safely converts an int64 to int32 with overflow checking.
//
// Parameters:
//   - value int64: The int64 value to convert.
//
// Returns:
//   - int32: The converted int32 value.
//   - error: An error if the value is outside the int32 range.
func SafeInt64ToInt32(value int64) (int32, error) {
	if value > math.MaxInt32 || value < math.MinInt32 {
		return 0, errors.New("value outside int32 range")
	}
	return int32(value), nil
}

// FormatThousands formats an integer with thousand separators.
//
// Parameters:
//   - value int: The integer to format.
//
// Returns:
//   - string: The formatted string with thousand separators.
func FormatThousands(value int) string {
	str := strconv.Itoa(value)
	n := len(str)
	if n <= 3 {
		return str
	}

	var result []byte
	for i := 0; i < n; i++ {
		if i > 0 && (n-i)%3 == 0 {
			result = append(result, ',')
		}
		result = append(result, str[i])
	}

	return string(result)
}
