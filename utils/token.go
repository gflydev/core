package utils

import "time"

// Token generate unique token
func Token(object ...string) string {
	// Make random data
	currentTime := time.Now().Format("20060102150405")
	randomNum := RandInt64(20)
	randomByte := RandByte(make([]byte, 50))

	// Token
	return Sha256(object, currentTime, randomNum, randomByte)
}
