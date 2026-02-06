package identity

import (
	"regexp"
	"strings"
)

var (
	// RFC-5322 simplified (safe for backend validation)
	emailRegex = regexp.MustCompile(
		`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`,
	)

	// E.164 international phone number
	// +1234567890 (10–15 digits)
	phoneRegex = regexp.MustCompile(
		`^\+?[1-9][0-9]{9,14}$`,
	)
)

// IsEmail checks if the input is an email address
func IsEmail(input string) bool {
	input = strings.TrimSpace(strings.ToLower(input))
	return emailRegex.MatchString(input)
}

// IsMobile checks if the input is an international mobile number
func IsMobile(input string) bool {
	input = strings.TrimSpace(input)
	return phoneRegex.MatchString(input)
}
