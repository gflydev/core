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

// RandByte returns a byte slice filled with cryptographically secure random
// bytes mapped to a specific character set.
//
// Parameters:
//   - dst []byte: The destination slice that will be filled with random bytes.
//
// Returns:
//   - []byte: The modified destination slice containing random bytes.
func RandByte(dst []byte) []byte {
	// Retrieve a buffer from the pool.
	buf := randBytesPool.Get()

	// Extend the buffer to match the required length.
	buf.B = ExtendByte(buf.B, len(dst))

	// Read cryptographically secure random bytes into the buffer.
	if _, err := crand.Read(buf.B); err != nil {
		panic(err) // Panic if random byte generation fails.
	}

	// The length of the destination slice.
	size := len(dst)

	// Fill the destination slice using the random bytes and mask them to fit the charset.
	for i, j := 0, 0; i < size; j++ {
		// Mask bytes to get an index into the character slice.
		if idx := int(buf.B[j%size] & charsetIdxMask); idx < len(charset) {
			dst[i] = charset[idx] // Map the random byte to a character in the charset.
			i++
		}
	}

	// Return the buffer to the pool.
	randBytesPool.Put(buf)

	return dst
}

// ExtendByte extends or truncates the provided slice to match the required length.
//
// Parameters:
//   - b []byte: The input slice to be extended or truncated.
//   - needLen int: The desired length of the resulting slice.
//
// Returns:
//   - []byte: The modified slice extended or truncated to the requested length.
func ExtendByte(b []byte, needLen int) []byte {
	// Extend the slice capacity if necessary.
	b = b[:cap(b)]
	if n := needLen - cap(b); n > 0 {
		b = append(b, make([]byte, n)...) // Append extra bytes if needed.
	}

	return b[:needLen] // Truncate or return the slice at the exact needed length.
}

// PrependByte prepends the provided bytes to the destination slice.
//
// Parameters:
//   - dst []byte: The destination slice to prepend bytes to.
//   - src ...byte: The source bytes to prepend to the destination slice.
//
// Returns:
//   - []byte: The new slice with the source bytes prepended to the destination slice.
func PrependByte(dst []byte, src ...byte) []byte {
	dstLen := len(dst) // Length of the destination slice.
	srcLen := len(src) // Length of the source bytes.

	// Extend the destination slice to accommodate the source bytes.
	dst = ExtendByte(dst, dstLen+srcLen)

	// Shift the original destination slice to make room for the source bytes.
	copy(dst[srcLen:], dst[:dstLen])

	// Copy the source bytes into the start of the new slice.
	copy(dst[:srcLen], src)

	return dst
}

// PrependByteStr prepends a string to the given byte slice.
//
// Parameters:
//   - dst []byte: The destination slice to prepend the string to.
//   - src string: The source string to prepend to the destination slice.
//
// Returns:
//   - []byte: The new slice with the source string prepended to the destination slice.
func PrependByteStr(dst []byte, src string) []byte {
	return PrependByte(dst, UnsafeBytes(src)...) // Convert the string to bytes and prepend.
}

// CopyByte creates a new copy of the given byte slice.
//
// Parameters:
//   - b []byte: The input slice to copy.
//
// Returns:
//   - []byte: A new slice with copied data from the input slice.
func CopyByte(b []byte) []byte {
	return []byte(UnsafeStr(b)) // Convert the slice to string and back to a new slice.
}

// EqualByte compares two byte slices for equality.
//
// Parameters:
//   - a []byte: The first byte slice.
//   - b []byte: The second byte slice.
//
// Returns:
//   - bool: True if the slices have the same length and identical contents, false otherwise.
func EqualByte(a, b []byte) bool {
	return UnsafeStr(a) == UnsafeStr(b) // Compare slices as strings.
}
