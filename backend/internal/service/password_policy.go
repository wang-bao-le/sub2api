package service

import (
	"unicode"
	"unicode/utf8"

	infraerrors "github.com/Wei-Shaw/sub2api/internal/pkg/errors"
)

// ValidatePassword enforces the policy for newly set local passwords.
func ValidatePassword(password string) error {
	if utf8.RuneCountInString(password) < 8 {
		return infraerrors.BadRequest("PASSWORD_TOO_SHORT", "password must be at least 8 characters")
	}

	hasLetter := false
	hasDigit := false
	for _, r := range password {
		hasLetter = hasLetter || unicode.IsLetter(r)
		hasDigit = hasDigit || unicode.IsDigit(r)
	}
	if !hasLetter || !hasDigit {
		return infraerrors.BadRequest("PASSWORD_WEAK", "password must contain at least one letter and one number")
	}
	return nil
}
