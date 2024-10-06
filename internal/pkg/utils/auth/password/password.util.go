package password

import (
	"regexp"

	"github.com/pkg/errors"
)

// Constants for password policies
const (
	minPasswordLength = 8
)

// ErrPasswordTooShort is returned when the password does not meet the minimum length requirement.
var ErrPasswordTooShort = errors.New("password is too short")

// ErrPasswordWeak is returned when the password does not meet complexity requirements.
var ErrPasswordWeak = errors.New("password is too weak")

// ValidatePasswordStrength ensures the password meets complexity requirements
func ValidatePasswordStrength(password string) error {
	if len(password) < minPasswordLength {
		return ErrPasswordTooShort
	}

	// Ensure the password contains at least one uppercase letter, one lowercase letter, one number, and one special character
	var (
		hasUpper  = regexp.MustCompile(`[A-Z]`).MatchString
		hasLower  = regexp.MustCompile(`[a-z]`).MatchString
		hasNumber = regexp.MustCompile(`[0-9]`).MatchString
		hasSymbol = regexp.MustCompile(`[!@#~$%^&*()_+|<>,.?/:;'"[\]{}-]`).MatchString
	)

	if !hasUpper(password) || !hasLower(password) || !hasNumber(password) || !hasSymbol(password) {
		return ErrPasswordWeak
	}

	return nil
}
