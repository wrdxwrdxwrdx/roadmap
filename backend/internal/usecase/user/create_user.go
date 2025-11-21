package user

import (
	"context"
	"fmt"
	"strings"
	"time"
	"unicode"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	userdto "roadmap/internal/domain/dto/user"
	userentity "roadmap/internal/domain/entities/user"
	userrepo "roadmap/internal/repository/user"
)

type CreateUserUseCase struct {
	userRepository userrepo.UserRepository
}

func NewCreateUserUseCase(userRepository userrepo.UserRepository) *CreateUserUseCase {
	return &CreateUserUseCase{userRepository: userRepository}
}

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

func (u *CreateUserUseCase) Execute(ctx context.Context, req userdto.CreateUserRequest) (userdto.CreateUserResponse, error) {
	emailExists, err := u.userRepository.EmailExists(ctx, req.Email)
	if err != nil {
		return userdto.CreateUserResponse{}, err
	}
	if emailExists {
		return userdto.CreateUserResponse{}, ErrEmailAlreadyExists
	}

	usernameExists, err := u.userRepository.UsernameExists(ctx, req.Username)
	if err != nil {
		return userdto.CreateUserResponse{}, err
	}
	if usernameExists {
		return userdto.CreateUserResponse{}, ErrUsernameAlreadyExists
	}

	if err := validatePassword(req.Password); err != nil {
		return userdto.CreateUserResponse{}, err
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return userdto.CreateUserResponse{}, err
	}

	now := time.Now()
	user := &userentity.User{
		ID:           uuid.New(),
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: string(passwordHash),
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	createdUser, err := u.userRepository.Create(ctx, user)
	if err != nil {
		return userdto.CreateUserResponse{}, err
	}

	return userdto.CreateUserResponse{
		ID:        createdUser.ID,
		Username:  createdUser.Username,
		Email:     createdUser.Email,
		CreatedAt: createdUser.CreatedAt,
		UpdatedAt: createdUser.UpdatedAt,
	}, nil
}
