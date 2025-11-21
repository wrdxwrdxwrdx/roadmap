package user

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	userdto "roadmap/internal/domain/dto/user"
	userentity "roadmap/internal/domain/entities/user"
	jwtservice "roadmap/internal/pkg/jwt"
)

type RegisterUseCaseTestSuite struct {
	suite.Suite
	useCase      *RegisterUseCase
	mockRepo     *MockUserRepository
	jwtService   *jwtservice.JWTService
	validRequest userdto.RegisterRequest
	validUser    *userentity.User
	ctx          context.Context
}

func (s *RegisterUseCaseTestSuite) SetupTest() {
	s.mockRepo = new(MockUserRepository)
	s.jwtService = jwtservice.NewJWTService("test-secret-key", 24*3600*1000000000)
	s.useCase = NewRegisterUseCase(s.mockRepo, s.jwtService)
	s.ctx = context.Background()

	s.validRequest = userdto.RegisterRequest{
		Email:    "test@example.com",
		Username: "testuser",
		Password: "SecurePass123!",
	}

	now := time.Now()
	s.validUser = &userentity.User{
		ID:           uuid.New(),
		Username:     s.validRequest.Username,
		Email:        s.validRequest.Email,
		PasswordHash: "$2a$10$hashedpassword",
		CreatedAt:    now,
		UpdatedAt:    now,
	}
}

func (s *RegisterUseCaseTestSuite) TearDownTest() {
	s.mockRepo.AssertExpectations(s.T())
}

func (s *RegisterUseCaseTestSuite) TestRegister_Success() {
	s.mockRepo.On("EmailExists", s.ctx, s.validRequest.Email).Return(false, nil)
	s.mockRepo.On("UsernameExists", s.ctx, s.validRequest.Username).Return(false, nil)
	s.mockRepo.On("Create", s.ctx, mock.AnythingOfType("*user.User")).Return(s.validUser, nil)

	response, err := s.useCase.Execute(s.ctx, s.validRequest)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), s.validUser.ID, response.ID)
	assert.Equal(s.T(), s.validUser.Username, response.Username)
	assert.Equal(s.T(), s.validUser.Email, response.Email)
	assert.NotEmpty(s.T(), response.Token)
	assert.False(s.T(), response.CreatedAt.IsZero())
	assert.False(s.T(), response.UpdatedAt.IsZero())
}

func (s *RegisterUseCaseTestSuite) TestRegister_EmailAlreadyExists() {
	s.mockRepo.On("EmailExists", s.ctx, s.validRequest.Email).Return(true, nil)

	response, err := s.useCase.Execute(s.ctx, s.validRequest)

	assert.Error(s.T(), err)
	assert.Equal(s.T(), ErrEmailAlreadyExists, err)
	assert.Equal(s.T(), uuid.Nil, response.ID)
	s.mockRepo.AssertNotCalled(s.T(), "Create", mock.Anything, mock.Anything)
}

func (s *RegisterUseCaseTestSuite) TestRegister_UsernameAlreadyExists() {
	s.mockRepo.On("EmailExists", s.ctx, s.validRequest.Email).Return(false, nil)
	s.mockRepo.On("UsernameExists", s.ctx, s.validRequest.Username).Return(true, nil)

	response, err := s.useCase.Execute(s.ctx, s.validRequest)

	assert.Error(s.T(), err)
	assert.Equal(s.T(), ErrUsernameAlreadyExists, err)
	assert.Equal(s.T(), uuid.Nil, response.ID)
	s.mockRepo.AssertNotCalled(s.T(), "Create", mock.Anything, mock.Anything)
}

func (s *RegisterUseCaseTestSuite) TestRegister_InvalidPassword() {
	testCases := []struct {
		name     string
		password string
	}{
		{"too short", "Short1!"},
		{"no uppercase", "lowercase123!"},
		{"no lowercase", "UPPERCASE123!"},
		{"no number", "NoNumberHere!"},
		{"no special", "NoSpecialChar123"},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			s.mockRepo.ExpectedCalls = nil
			s.mockRepo.Calls = nil

			req := s.validRequest
			req.Password = tc.password

			s.mockRepo.On("EmailExists", s.ctx, req.Email).Return(false, nil)
			s.mockRepo.On("UsernameExists", s.ctx, req.Username).Return(false, nil)

			response, err := s.useCase.Execute(s.ctx, req)

			assert.Error(s.T(), err)
			var passwordErr *PasswordValidationError
			assert.ErrorAs(s.T(), err, &passwordErr)
			assert.Equal(s.T(), uuid.Nil, response.ID)
			s.mockRepo.AssertNotCalled(s.T(), "Create", mock.Anything, mock.Anything)
		})
	}
}

func (s *RegisterUseCaseTestSuite) TestRegister_EmailExistsError() {
	repoError := errors.New("database error")
	s.mockRepo.On("EmailExists", s.ctx, s.validRequest.Email).Return(false, repoError)

	response, err := s.useCase.Execute(s.ctx, s.validRequest)

	assert.Error(s.T(), err)
	assert.Equal(s.T(), repoError, err)
	assert.Equal(s.T(), uuid.Nil, response.ID)
}

func (s *RegisterUseCaseTestSuite) TestRegister_UsernameExistsError() {
	s.mockRepo.On("EmailExists", s.ctx, s.validRequest.Email).Return(false, nil)
	repoError := errors.New("database error")
	s.mockRepo.On("UsernameExists", s.ctx, s.validRequest.Username).Return(false, repoError)

	response, err := s.useCase.Execute(s.ctx, s.validRequest)

	assert.Error(s.T(), err)
	assert.Equal(s.T(), repoError, err)
	assert.Equal(s.T(), uuid.Nil, response.ID)
}

func (s *RegisterUseCaseTestSuite) TestRegister_CreateError() {
	s.mockRepo.On("EmailExists", s.ctx, s.validRequest.Email).Return(false, nil)
	s.mockRepo.On("UsernameExists", s.ctx, s.validRequest.Username).Return(false, nil)
	repoError := errors.New("create error")
	s.mockRepo.On("Create", s.ctx, mock.AnythingOfType("*user.User")).Return(nil, repoError)

	response, err := s.useCase.Execute(s.ctx, s.validRequest)

	assert.Error(s.T(), err)
	assert.Equal(s.T(), repoError, err)
	assert.Equal(s.T(), uuid.Nil, response.ID)
}

func (s *RegisterUseCaseTestSuite) TestRegister_JWTGenerationError() {
	s.mockRepo.On("EmailExists", s.ctx, s.validRequest.Email).Return(false, nil)
	s.mockRepo.On("UsernameExists", s.ctx, s.validRequest.Username).Return(false, nil)
	s.mockRepo.On("Create", s.ctx, mock.AnythingOfType("*user.User")).Return(s.validUser, nil)

	response, err := s.useCase.Execute(s.ctx, s.validRequest)

	assert.NoError(s.T(), err)
	assert.NotEmpty(s.T(), response.Token)
}

func TestRegisterUseCaseTestSuite(t *testing.T) {
	suite.Run(t, new(RegisterUseCaseTestSuite))
}

