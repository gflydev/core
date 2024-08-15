package utils

import (
	crand "crypto/rand"
	"github.com/valyala/bytebufferpool"
)

var randBytesPool = bytebufferpool.Pool{}

const (
	charset        = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	charsetIdxBits = 6                     // 6 bits to represent a charset index
	charsetIdxMask = 1<<charsetIdxBits - 1 // All 1-bits, as many as charsetIdxBits
)

// RandByte returns dst with a cryptographically secure string random bytes.
//
// NOTE: Make sure that dst has the length you need.
func RandByte(dst []byte) []byte {
	buf := randBytesPool.Get()
	buf.B = ExtendByte(buf.B, len(dst))

	if _, err := crand.Read(buf.B); err != nil {
		panic(err)
	}

	size := len(dst)

	for i, j := 0, 0; i < size; j++ {
		// Mask bytes to get an index into the character slice.
		if idx := int(buf.B[j%size] & charsetIdxMask); idx < len(charset) {
			dst[i] = charset[idx]
			i++
		}
	}

	randBytesPool.Put(buf)

	return dst
}

// ExtendByte extends b to needLen bytes.
//
// Example b have [10 23 34 129], needLen 2. So return [10 23]
// Example b have [10 23 34 129], needLen 10. So return [10 23 34 129 0 0 0 0 0 0]
func ExtendByte(b []byte, needLen int) []byte {
	b = b[:cap(b)]
	if n := needLen - cap(b); n > 0 {
		b = append(b, make([]byte, n)...)
	}

	return b[:needLen]
}

// PrependByte prepends bytes into a given byte slice.
//
// Example dst have [10 23 34] src is 44 54. So return [44 54 10 23 34]
func PrependByte(dst []byte, src ...byte) []byte {
	dstLen := len(dst)
	srcLen := len(src)

	dst = ExtendByte(dst, dstLen+srcLen)
	copy(dst[srcLen:], dst[:dstLen])
	copy(dst[:srcLen], src)

	return dst
}

// PrependByteStr prepends a string into a given byte slice.
func PrependByteStr(dst []byte, src string) []byte {
	return PrependByte(dst, UnsafeBytes(src)...)
}

// CopyByte returns a copy of byte slice in a new pointer.
func CopyByte(b []byte) []byte {
	return []byte(UnsafeStr(b))
}

// EqualByte reports whether a and b
// are the same length and contain the same bytes.
// A nil argument is equivalent to an empty slice.
func EqualByte(a, b []byte) bool {
	return UnsafeStr(a) == UnsafeStr(b)
}
