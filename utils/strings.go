package utils

import (
	"github.com/valyala/bytebufferpool"
	"github.com/valyala/fasthttp"
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
