package user

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPasswordValidationError_Error(t *testing.T) {
	err := &PasswordValidationError{
		Message: "password must contain at least 8 characters",
	}

	assert.Equal(t, "password must contain at least 8 characters", err.Error())
}

func TestPasswordValidationError_IsError(t *testing.T) {
	err := &PasswordValidationError{
		Message: "test error",
	}

	var passwordErr *PasswordValidationError
	assert.True(t, errors.As(err, &passwordErr))
	assert.Equal(t, "test error", passwordErr.Message)
}

func TestErrors_Constants(t *testing.T) {
	assert.NotNil(t, ErrEmailAlreadyExists)
	assert.NotNil(t, ErrUsernameAlreadyExists)
	assert.NotNil(t, ErrInvalidCredentials)

	assert.Equal(t, "email already exists", ErrEmailAlreadyExists.Error())
	assert.Equal(t, "username already exists", ErrUsernameAlreadyExists.Error())
	assert.Equal(t, "invalid email or password", ErrInvalidCredentials.Error())
}

func TestErrors_Is(t *testing.T) {
	err1 := ErrEmailAlreadyExists
	err2 := ErrEmailAlreadyExists

	assert.True(t, errors.Is(err1, err2))
	assert.True(t, errors.Is(err1, ErrEmailAlreadyExists))

	assert.False(t, errors.Is(ErrEmailAlreadyExists, ErrUsernameAlreadyExists))
	assert.False(t, errors.Is(ErrEmailAlreadyExists, ErrInvalidCredentials))
}
