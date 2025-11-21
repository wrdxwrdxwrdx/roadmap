package user

import (
	"context"
	"fmt"

	"golang.org/x/crypto/bcrypt"

	userdto "roadmap/internal/domain/dto/user"
	jwtservice "roadmap/internal/pkg/jwt"
	userrepo "roadmap/internal/repository/user"
)

type LoginUseCase struct {
	userRepository userrepo.UserRepository
	jwtService     *jwtservice.JWTService
}

func NewLoginUseCase(userRepository userrepo.UserRepository, jwtService *jwtservice.JWTService) *LoginUseCase {
	return &LoginUseCase{
		userRepository: userRepository,
		jwtService:     jwtService,
	}
}

func (u *LoginUseCase) Execute(ctx context.Context, req userdto.LoginRequest) (userdto.LoginResponse, error) {
	user, err := u.userRepository.GetByEmail(ctx, req.Email)
	if err != nil {
		return userdto.LoginResponse{}, ErrInvalidCredentials
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))
	if err != nil {
		return userdto.LoginResponse{}, ErrInvalidCredentials
	}

	token, err := u.jwtService.GenerateToken(
		user.ID.String(),
		user.Username,
		user.Email,
	)
	if err != nil {
		return userdto.LoginResponse{}, fmt.Errorf("failed to generate token: %w", err)
	}

	return userdto.LoginResponse{
		Token: token,
	}, nil
}
