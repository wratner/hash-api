package utils

import (
	"crypto/sha512"
	"encoding/base64"
)

// HashPassword generates a SHA512 hash of the password
func HashPassword(password []byte) ([]byte, error) {
	sha512 := sha512.New()
	_, err := sha512.Write(password)
	if err != nil {
		return nil, err
	}
	return sha512.Sum(nil), nil
}

// Base64 encodes the data to base64 format.
func Base64(data []byte) string {
	return base64.StdEncoding.EncodeToString(data)
}
