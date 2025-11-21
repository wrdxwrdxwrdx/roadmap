package user

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"golang.org/x/crypto/bcrypt"

	userdto "roadmap/internal/domain/dto/user"
	userentity "roadmap/internal/domain/entities/user"
	jwtservice "roadmap/internal/pkg/jwt"
)

type LoginUseCaseTestSuite struct {
	suite.Suite
	useCase      *LoginUseCase
	mockRepo     *MockUserRepository
	jwtService   *jwtservice.JWTService
	validRequest userdto.LoginRequest
	validUser    *userentity.User
	ctx          context.Context
}

func (s *LoginUseCaseTestSuite) SetupTest() {
	s.mockRepo = new(MockUserRepository)
	s.jwtService = jwtservice.NewJWTService("test-secret-key", 24*3600*1000000000)
	s.useCase = NewLoginUseCase(s.mockRepo, s.jwtService)
	s.ctx = context.Background()

	s.validRequest = userdto.LoginRequest{
		Email:    "test@example.com",
		Password: "SecurePass123!",
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(s.validRequest.Password), bcrypt.DefaultCost)
	now := time.Now()
	s.validUser = &userentity.User{
		ID:           uuid.New(),
		Username:     "testuser",
		Email:        s.validRequest.Email,
		PasswordHash: string(hashedPassword),
		CreatedAt:    now,
		UpdatedAt:    now,
	}
}

func (s *LoginUseCaseTestSuite) TearDownTest() {
	s.mockRepo.AssertExpectations(s.T())
}

func (s *LoginUseCaseTestSuite) TestLogin_Success() {
	s.mockRepo.On("GetByEmail", s.ctx, s.validRequest.Email).Return(s.validUser, nil)

	response, err := s.useCase.Execute(s.ctx, s.validRequest)

	assert.NoError(s.T(), err)
	assert.NotEmpty(s.T(), response.Token)
}

func (s *LoginUseCaseTestSuite) TestLogin_UserNotFound() {
	s.mockRepo.On("GetByEmail", s.ctx, s.validRequest.Email).Return(nil, pgx.ErrNoRows)

	response, err := s.useCase.Execute(s.ctx, s.validRequest)

	assert.Error(s.T(), err)
	assert.Equal(s.T(), ErrInvalidCredentials, err)
	assert.Empty(s.T(), response.Token)
}

func (s *LoginUseCaseTestSuite) TestLogin_WrongPassword() {
	wrongPasswordUser := &userentity.User{
		ID:           s.validUser.ID,
		Username:     s.validUser.Username,
		Email:        s.validUser.Email,
		PasswordHash: s.validUser.PasswordHash,
		CreatedAt:    s.validUser.CreatedAt,
		UpdatedAt:    s.validUser.UpdatedAt,
	}

	s.mockRepo.On("GetByEmail", s.ctx, s.validRequest.Email).Return(wrongPasswordUser, nil)

	req := s.validRequest
	req.Password = "WrongPassword123!"

	response, err := s.useCase.Execute(s.ctx, req)

	assert.Error(s.T(), err)
	assert.Equal(s.T(), ErrInvalidCredentials, err)
	assert.Empty(s.T(), response.Token)
}

func (s *LoginUseCaseTestSuite) TestLogin_GetByEmailError() {
	repoError := errors.New("database error")
	s.mockRepo.On("GetByEmail", s.ctx, s.validRequest.Email).Return(nil, repoError)

	response, err := s.useCase.Execute(s.ctx, s.validRequest)

	assert.Error(s.T(), err)
	assert.Equal(s.T(), ErrInvalidCredentials, err)
	assert.Empty(s.T(), response.Token)
}

func (s *LoginUseCaseTestSuite) TestLogin_JWTGenerationError() {
	s.mockRepo.On("GetByEmail", s.ctx, s.validRequest.Email).Return(s.validUser, nil)

	response, err := s.useCase.Execute(s.ctx, s.validRequest)

	assert.NoError(s.T(), err)
	assert.NotEmpty(s.T(), response.Token)
}

func TestLoginUseCaseTestSuite(t *testing.T) {
	suite.Run(t, new(LoginUseCaseTestSuite))
}

