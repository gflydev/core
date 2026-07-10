package utils

import (
	"golang.org/x/crypto/bcrypt"
)

// GeneratePassword func for a making hash & salt with user password.
//
// It uses bcrypt.DefaultCost (10). MinCost (4) must never be used for real
// passwords: it is trivially brute-forceable. On failure it returns an empty
// string (never the error text) so a bogus value can never be persisted as a
// hash; ComparePasswords rejects an empty hash, i.e. it fails safe. Use
// GeneratePasswordE when you need to inspect the error.
func GeneratePassword(p string) string {
	hash, _ := GeneratePasswordE(p)

	return hash
}

// GeneratePasswordE hashes the password and returns any error from bcrypt.
func GeneratePasswordE(p string) (string, error) {
	// Normalize password from string to []byte.
	bytePwd := UnsafeBytes(p)

	hash, err := bcrypt.GenerateFromPassword(bytePwd, bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	// GenerateFromPassword returns a byte slice so we need to
	// convert the bytes to a string and return it.
	return string(hash), nil
}

// ComparePasswords func for a comparing password.
func ComparePasswords(hashedPwd, inputPwd string) bool {
	// Since we'll be getting the hashed password from the DB it will be a string,
	// so we'll need to convert it to a byte slice.
	byteHash := UnsafeBytes(hashedPwd)
	byteInput := UnsafeBytes(inputPwd)

	// Return result.
	if err := bcrypt.CompareHashAndPassword(byteHash, byteInput); err != nil {
		return false
	}

	return true
}
