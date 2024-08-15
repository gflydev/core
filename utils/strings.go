package utils

import (
	"github.com/valyala/bytebufferpool"
	"github.com/valyala/fasthttp"
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

// QuoteStr escape special characters in a given string
func QuoteStr(raw string) string {
	bb := bytebufferpool.Get()
	quoted := UnsafeStr(fasthttp.AppendQuotedArg(bb.B, UnsafeBytes(raw)))
	bytebufferpool.Put(bb)

	return quoted
}
