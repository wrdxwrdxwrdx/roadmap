package user

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	userdto "roadmap/internal/domain/dto/user"
	userentity "roadmap/internal/domain/entities/user"
	jwtservice "roadmap/internal/pkg/jwt"
	userrepo "roadmap/internal/repository/user"
)

type RegisterUseCase struct {
	userRepository userrepo.UserRepository
	jwtService     *jwtservice.JWTService
}

func NewRegisterUseCase(userRepository userrepo.UserRepository, jwtService *jwtservice.JWTService) *RegisterUseCase {
	return &RegisterUseCase{
		userRepository: userRepository,
		jwtService:     jwtService,
	}
}

func (u *RegisterUseCase) Execute(ctx context.Context, req userdto.RegisterRequest) (userdto.RegisterResponse, error) {
	emailExists, err := u.userRepository.EmailExists(ctx, req.Email)
	if err != nil {
		return userdto.RegisterResponse{}, err
	}
	if emailExists {
		return userdto.RegisterResponse{}, ErrEmailAlreadyExists
	}

	usernameExists, err := u.userRepository.UsernameExists(ctx, req.Username)
	if err != nil {
		return userdto.RegisterResponse{}, err
	}
	if usernameExists {
		return userdto.RegisterResponse{}, ErrUsernameAlreadyExists
	}

	if err := validatePassword(req.Password); err != nil {
		return userdto.RegisterResponse{}, err
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return userdto.RegisterResponse{}, err
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
		return userdto.RegisterResponse{}, err
	}

	token, err := u.jwtService.GenerateToken(
		createdUser.ID.String(),
		createdUser.Username,
		createdUser.Email,
	)
	if err != nil {
		return userdto.RegisterResponse{}, fmt.Errorf("failed to generate token: %w", err)
	}

	return userdto.RegisterResponse{
		ID:        createdUser.ID,
		Username:  createdUser.Username,
		Email:     createdUser.Email,
		Token:     token,
		CreatedAt: createdUser.CreatedAt,
		UpdatedAt: createdUser.UpdatedAt,
	}, nil
}
