package utils

import (
	"crypto/rand"
	"math/big"
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
