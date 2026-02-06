package identity

import "strings"

// NormalizeEmail normalizes an email address
func NormalizeEmail(input string) string {
	return strings.ToLower(strings.TrimSpace(input))
}

// NormalizePhone ensures phone number is trimmed
// (E.164 should already be enforced by frontend)
func NormalizePhone(input string) string {
	return strings.TrimSpace(input)
}
