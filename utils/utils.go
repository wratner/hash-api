package utils

import (
	"crypto/sha512"
	"encoding/base64"
)

// HashPassword generates a SHA512 hash of the password
func HashPassword(password []byte) []byte {
	sha512 := sha512.New()
	sha512.Write(password)
	return sha512.Sum(nil)
}

// Base64 will encode the data to base64 format
func Base64(data []byte) string {
	return base64.StdEncoding.EncodeToString(data)
}
