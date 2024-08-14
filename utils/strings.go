package utils

import (
	"unsafe"
)

// UnsafeBytes returns a byte pointer without allocation.
func UnsafeBytes(s string) []byte {
	return unsafe.Slice(unsafe.StringData(s), len(s))
}

// UnsafeStr returns a string pointer without allocation.
func UnsafeStr(b []byte) string {
	return unsafe.String(unsafe.SliceData(b), len(b))
}

// CopyStr copies a string to make it immutable
func CopyStr(s string) string {
	return string(UnsafeBytes(s))
}

// CopyBytes copies a slice to make it immutable
func CopyBytes(b []byte) []byte {
	tmp := make([]byte, len(b))
	copy(tmp, b)
	return tmp
}

// IncludeStr returns true or false if given string is in slice.
func IncludeStr(slice []string, s string) bool {
	return IndexOfStr(slice, s) != -1
}

// IndexOfStr returns index position in slice from given string
// If value is -1, the string does not found.
func IndexOfStr(slice []string, s string) int {
	for i, v := range slice {
		if v == s {
			return i
		}
	}

	return -1
}
