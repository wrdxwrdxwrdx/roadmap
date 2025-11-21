package user

import "errors"

var (
	ErrEmailAlreadyExists    = errors.New("email already exists")
	ErrUsernameAlreadyExists = errors.New("username already exists")
	ErrInvalidCredentials    = errors.New("invalid email or password")
)

type PasswordValidationError struct {
	Message string
}

func (e *PasswordValidationError) Error() string {
	return e.Message
}
