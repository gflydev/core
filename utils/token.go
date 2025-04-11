package utils

import "time"

// Token generates a unique token by combining input strings, the current timestamp,
// a random integer, and a random byte slice, hashed using SHA256.
//
// Parameters:
//   - object ...string: Optional variadic input strings that will be included in the token generation.
//
// Returns:
//   - string: A unique SHA256 hash representing the generated token.
func Token(object ...string) string {
	// Generate the current time as a formatted string.
	currentTime := time.Now().Format("20060102150405") // Format: YYYYMMDDHHMMSS

	// Generate a cryptographically secure random integer.
	randomNum := RandInt64(20) // Random integer in range [0, 20).

	// Generate a cryptographically secure random byte slice.
	randomByte := RandByte(make([]byte, 50)) // Random byte slice of length 50.

	// Combine the input strings, current time, random number, and random bytes to create a unique token.
	return Sha256(object, currentTime, randomNum, randomByte)
}
