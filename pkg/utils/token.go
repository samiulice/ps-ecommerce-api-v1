// Package utils provides cryptographic helpers.
package utils

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
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

// Sscanf is a wrapper around fmt.Sscanf that returns the number of items parsed and any error.
func Sscanf(str, format string, a ...any) (int, error) {
	return fmt.Sscanf(str, format, a...)
}
