package user

import (
	"context"
	"time"

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
