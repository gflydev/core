// Package utils provides a collection of utility functions for common operations
// in the gFly framework. This package is organized into multiple files, each focusing
// on a specific type of utility.
//
// The strings.go file contains utilities for string manipulation, including:
//   - Converting between strings and byte slices without allocations
//   - Checking if a string is present in a slice of strings
//   - Finding the index of a string in a slice
//   - Escaping special characters in strings
//   - Truncating strings with optional ellipsis
//   - Checking if strings are empty or blank
//   - Converting strings to different case formats (camelCase, PascalCase, snake_case)
//
// Many of these utilities are optimized for performance, using techniques like
// unsafe pointers for zero-allocation conversions where appropriate.
//
// Usage Examples:
//
//	// Convert between strings and bytes without allocation
//	bytes := utils.UnsafeBytes(myString)
//	str := utils.UnsafeStr(myBytes)
//
//	// Create a copy of a string
//	strCopy := utils.CopyStr(myString)
//
//	// Check if a string is in a slice
//	isIncluded := utils.IncludeStr(mySlice, "search")
//
//	// Truncate a string
//	truncated := utils.Truncate(longString, 50, true)
//
//	// Check if a string is empty or contains only whitespace
//	isEmpty := utils.IsEmpty(myString)
//	isBlank := utils.IsBlank(myString)
//
//	// Convert string case
//	camel := utils.ToCamelCase("my-variable-name")
//	pascal := utils.ToPascalCase("my-class-name")
//	snake := utils.ToSnakeCase("MyMethodName")
package utils

import (
	"github.com/valyala/bytebufferpool"
	"github.com/valyala/fasthttp"
	"strings"
	"unicode"
	"unsafe"
)

// UnsafeBytes returns a byte pointer without allocation.
//
// Parameters:
//   - s string: The input string to convert to a byte slice.
//
// Returns:
//   - []byte: A byte slice pointing to the same memory as the input string.
func UnsafeBytes(s string) []byte {
	return unsafe.Slice(unsafe.StringData(s), len(s))
}

// UnsafeStr returns a string pointer without allocation.
//
// Parameters:
//   - b []byte: The input byte slice to convert to a string.
//
// Returns:
//   - string: A string pointing to the same memory as the input byte slice.
func UnsafeStr(b []byte) string {
	return unsafe.String(unsafe.SliceData(b), len(b))
}

// CopyStr creates a new copy of the input string to ensure immutability.
//
// Parameters:
//   - s string: The input string to copy.
//
// Returns:
//   - string: A new immutable copy of the input string.
func CopyStr(s string) string {
	return string(UnsafeBytes(s))
}

// IncludeStr checks if a given string is present in a slice of strings.
//
// Parameters:
//   - slice []string: The slice of strings to search in.
//   - s string: The string to look for in the slice.
//
// Returns:
//   - bool: True if the string is present in the slice, false otherwise.
func IncludeStr(slice []string, s string) bool {
	return IndexOfStr(slice, s) != -1
}

// IndexOfStr finds the index of a given string in a slice of strings.
//
// Parameters:
//   - slice []string: The slice of strings to search in.
//   - s string: The string to look for in the slice.
//
// Returns:
//   - int: The index of the string in the slice, or -1 if the string is not found.
func IndexOfStr(slice []string, s string) int {
	for i, v := range slice {
		if v == s {
			return i
		}
	}

	return -1
}

// QuoteStr escapes special characters in a given string.
//
// Parameters:
//   - raw string: The input string to escape.
//
// Returns:
//   - string: The escaped string with special characters properly quoted.
//
// Variables:
//   - bb *bytebufferpool.ByteBuffer: A byte buffer obtained from the bytebufferpool for temporary storage.
//   - quoted string: The resulting escaped string.
func QuoteStr(raw string) string {
	bb := bytebufferpool.Get()
	quoted := UnsafeStr(fasthttp.AppendQuotedArg(bb.B, UnsafeBytes(raw)))
	bytebufferpool.Put(bb)

	return quoted
}

// Truncate truncates a string to the specified length and adds an ellipsis if truncated.
//
// Parameters:
//   - s string: The input string to truncate.
//   - length int: The maximum length of the truncated string (excluding ellipsis).
//   - withEllipsis bool: Whether to add an ellipsis ("...") if the string is truncated.
//
// Returns:
//   - string: The truncated string, with an ellipsis if requested and if truncation occurred.
func Truncate(s string, length int, withEllipsis bool) string {
	if length <= 0 {
		return ""
	}

	if len(s) <= length {
		return s
	}

	if withEllipsis {
		if length > 3 {
			return s[:length-3] + "..."
		}
		return s[:length]
	}

	return s[:length]
}

// IsEmpty checks if a string is empty.
//
// Parameters:
//   - s string: The input string to check.
//
// Returns:
//   - bool: True if the string is empty, false otherwise.
func IsEmpty(s string) bool {
	return len(s) == 0
}

// IsBlank checks if a string is empty or contains only whitespace.
//
// Parameters:
//   - s string: The input string to check.
//
// Returns:
//   - bool: True if the string is empty or contains only whitespace, false otherwise.
func IsBlank(s string) bool {
	if len(s) == 0 {
		return true
	}

	for _, r := range s {
		if !unicode.IsSpace(r) {
			return false
		}
	}

	return true
}

// ToCamelCase converts a string to camelCase.
//
// Parameters:
//   - s string: The input string to convert.
//
// Returns:
//   - string: The camelCase version of the input string.
func ToCamelCase(s string) string {
	if s == "" {
		return s
	}

	words := strings.FieldsFunc(s, func(r rune) bool {
		return r == '_' || r == '-' || r == ' ' || unicode.IsPunct(r)
	})

	if len(words) == 0 {
		return ""
	}

	result := strings.ToLower(words[0])
	for i := 1; i < len(words); i++ {
		if words[i] == "" {
			continue
		}
		result += strings.ToUpper(words[i][:1]) + strings.ToLower(words[i][1:])
	}

	return result
}

// ToPascalCase converts a string to PascalCase.
//
// Parameters:
//   - s string: The input string to convert.
//
// Returns:
//   - string: The PascalCase version of the input string.
func ToPascalCase(s string) string {
	if s == "" {
		return s
	}

	words := strings.FieldsFunc(s, func(r rune) bool {
		return r == '_' || r == '-' || r == ' ' || unicode.IsPunct(r)
	})

	if len(words) == 0 {
		return ""
	}

	var result string
	for _, word := range words {
		if word == "" {
			continue
		}
		result += strings.ToUpper(word[:1]) + strings.ToLower(word[1:])
	}

	return result
}

// ToSnakeCase converts a string to snake_case.
//
// Parameters:
//   - s string: The input string to convert.
//
// Returns:
//   - string: The snake_case version of the input string.
func ToSnakeCase(s string) string {
	if s == "" {
		return s
	}

	var result strings.Builder
	for i, r := range s {
		if unicode.IsUpper(r) {
			if i > 0 {
				result.WriteRune('_')
			}
			result.WriteRune(unicode.ToLower(r))
		} else if r == ' ' || r == '-' {
			result.WriteRune('_')
		} else {
			result.WriteRune(r)
		}
	}

	return result.String()
}
