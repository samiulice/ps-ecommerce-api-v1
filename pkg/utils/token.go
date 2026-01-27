// Package utils provides cryptographic helpers.
package utils

import (
	"crypto/rand"
	"encoding/base64"
)

// GenerateRandomToken creates a secure random string for refresh tokens.
func GenerateRandomToken() (string, error) {
	b := make([]byte, 32) // 256-bit token
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	return base64.RawURLEncoding.EncodeToString(b), nil
}
