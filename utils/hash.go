package utils

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
)

// Sha256 generates a SHA256 hash from a list of arguments.
//
// Parameters:
//   - args ...any: Variadic arguments to be hashed. The arguments are concatenated
//     into a single string, separated by hyphens ("-").
//
// Returns:
//   - string: The resulting SHA256 hash as a hexadecimal-encoded string.
func Sha256(args ...any) string {
	// Slice to hold format specifiers for each argument ("%v").
	var strSlice []string
	for i := 0; i < len(args); i++ {
		strSlice = append(strSlice, "%v")
	}

	// Create a new SHA256 hash.
	hash := sha256.New()

	// Format the arguments into a single string, separated by hyphens ("-").
	// Expand args with "..." so each "%v" verb consumes one argument; passing
	// the slice directly would render the whole slice into the first verb and
	// leave the rest as "%!v(MISSING)".
	code := fmt.Sprintf(strings.Join(strSlice, "-"), args...)

	// Write the formatted string into the hash.
	hash.Write([]byte(code))

	// Return the hexadecimal-encoded hash value.
	return hex.EncodeToString(hash.Sum(nil))
}

// MD5 generates an MD5 hash from a string.
//
// Parameters:
//   - s string: The input string to be hashed.
//
// Returns:
//   - string: The resulting MD5 hash as a hexadecimal-encoded string.
func MD5(s string) string {
	// Create a new MD5 hash.
	hash := md5.New()

	// Write the input string into the hash.
	hash.Write([]byte(s))

	// Return the hexadecimal-encoded hash value.
	return hex.EncodeToString(hash.Sum(nil))
}
