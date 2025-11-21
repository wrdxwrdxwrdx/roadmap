package user

import "errors"

var (
	ErrEmailAlreadyExists    = errors.New("email already exists")
	ErrUsernameAlreadyExists = errors.New("username already exists")
)

type PasswordValidationError struct {
	Message string
}

func (e *PasswordValidationError) Error() string {
	return e.Message
}
