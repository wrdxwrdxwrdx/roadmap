package user

import (
	"fmt"
	"strings"
	"unicode"
)

func validatePassword(password string) error {
	const minLength = 8

	if len(password) < minLength {
		return &PasswordValidationError{
			Message: fmt.Sprintf("password must be at least %d characters long", minLength),
		}
	}

	var (
		hasUpper   = false
		hasLower   = false
		hasNumber  = false
		hasSpecial = false
	)

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char) || strings.ContainsRune("!@#$%^&*()_+-=[]{}|;:,.<>?", char):
			hasSpecial = true
		}
	}

	var missing []string
	if !hasUpper {
		missing = append(missing, "uppercase letter")
	}
	if !hasLower {
		missing = append(missing, "lowercase letter")
	}
	if !hasNumber {
		missing = append(missing, "number")
	}
	if !hasSpecial {
		missing = append(missing, "special character")
	}

	if len(missing) > 0 {
		return &PasswordValidationError{
			Message: fmt.Sprintf("password must contain at least one: %s", strings.Join(missing, ", ")),
		}
	}

	return nil
}
