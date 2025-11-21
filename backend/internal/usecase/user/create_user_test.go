package user

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	userdto "roadmap/internal/domain/dto/user"
	userentity "roadmap/internal/domain/entities/user"
)

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(ctx context.Context, user *userentity.User) (*userentity.User, error) {
	args := m.Called(ctx, user)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*userentity.User), args.Error(1)
}

func (m *MockUserRepository) GetByID(ctx context.Context, id uuid.UUID) (*userentity.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*userentity.User), args.Error(1)
}

func (m *MockUserRepository) GetByEmail(ctx context.Context, email string) (*userentity.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*userentity.User), args.Error(1)
}

func (m *MockUserRepository) EmailExists(ctx context.Context, email string) (bool, error) {
	args := m.Called(ctx, email)
	return args.Bool(0), args.Error(1)
}

func (m *MockUserRepository) UsernameExists(ctx context.Context, username string) (bool, error) {
	args := m.Called(ctx, username)
	return args.Bool(0), args.Error(1)
}

type CreateUserUseCaseTestSuite struct {
	suite.Suite
	useCase      *CreateUserUseCase
	mockRepo     *MockUserRepository
	validRequest userdto.CreateUserRequest
	validUser    *userentity.User
	ctx          context.Context
}

func (s *CreateUserUseCaseTestSuite) SetupTest() {
	s.mockRepo = new(MockUserRepository)
	s.useCase = NewCreateUserUseCase(s.mockRepo)
	s.ctx = context.Background()

	s.validRequest = userdto.CreateUserRequest{
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

func (s *CreateUserUseCaseTestSuite) TearDownTest() {
	s.mockRepo.AssertExpectations(s.T())
}

func (s *CreateUserUseCaseTestSuite) TestCreateUser_Success() {
	s.mockRepo.On("EmailExists", s.ctx, s.validRequest.Email).Return(false, nil)
	s.mockRepo.On("UsernameExists", s.ctx, s.validRequest.Username).Return(false, nil)
	s.mockRepo.On("Create", s.ctx, mock.AnythingOfType("*user.User")).Return(s.validUser, nil)

	response, err := s.useCase.Execute(s.ctx, s.validRequest)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), s.validUser.ID, response.ID)
	assert.Equal(s.T(), s.validUser.Username, response.Username)
	assert.Equal(s.T(), s.validUser.Email, response.Email)
	assert.Equal(s.T(), s.validUser.CreatedAt, response.CreatedAt)
	assert.Equal(s.T(), s.validUser.UpdatedAt, response.UpdatedAt)
}

func (s *CreateUserUseCaseTestSuite) TestCreateUser_EmailAlreadyExists() {
	s.mockRepo.On("EmailExists", s.ctx, s.validRequest.Email).Return(true, nil)

	response, err := s.useCase.Execute(s.ctx, s.validRequest)

	assert.Error(s.T(), err)
	assert.True(s.T(), errors.Is(err, ErrEmailAlreadyExists))
	assert.Equal(s.T(), userdto.CreateUserResponse{}, response)

	s.mockRepo.AssertNotCalled(s.T(), "UsernameExists", s.ctx, s.validRequest.Username)
}

func (s *CreateUserUseCaseTestSuite) TestCreateUser_UsernameAlreadyExists() {
	s.mockRepo.On("EmailExists", s.ctx, s.validRequest.Email).Return(false, nil)
	s.mockRepo.On("UsernameExists", s.ctx, s.validRequest.Username).Return(true, nil)

	response, err := s.useCase.Execute(s.ctx, s.validRequest)

	assert.Error(s.T(), err)
	assert.True(s.T(), errors.Is(err, ErrUsernameAlreadyExists))
	assert.Equal(s.T(), userdto.CreateUserResponse{}, response)
}

func (s *CreateUserUseCaseTestSuite) TestCreateUser_EmailExistsError() {
	repoError := errors.New("database connection error")
	s.mockRepo.On("EmailExists", s.ctx, s.validRequest.Email).Return(false, repoError)

	response, err := s.useCase.Execute(s.ctx, s.validRequest)

	assert.Error(s.T(), err)
	assert.Equal(s.T(), repoError, err)
	assert.Equal(s.T(), userdto.CreateUserResponse{}, response)
}

func (s *CreateUserUseCaseTestSuite) TestCreateUser_UsernameExistsError() {
	repoError := errors.New("database connection error")
	s.mockRepo.On("EmailExists", s.ctx, s.validRequest.Email).Return(false, nil)
	s.mockRepo.On("UsernameExists", s.ctx, s.validRequest.Username).Return(false, repoError)

	response, err := s.useCase.Execute(s.ctx, s.validRequest)

	assert.Error(s.T(), err)
	assert.Equal(s.T(), repoError, err)
	assert.Equal(s.T(), userdto.CreateUserResponse{}, response)
}

func (s *CreateUserUseCaseTestSuite) TestCreateUser_CreateError() {
	repoError := errors.New("failed to insert user")
	s.mockRepo.On("EmailExists", s.ctx, s.validRequest.Email).Return(false, nil)
	s.mockRepo.On("UsernameExists", s.ctx, s.validRequest.Username).Return(false, nil)
	s.mockRepo.On("Create", s.ctx, mock.AnythingOfType("*user.User")).Return(nil, repoError)

	response, err := s.useCase.Execute(s.ctx, s.validRequest)

	assert.Error(s.T(), err)
	assert.Equal(s.T(), repoError, err)
	assert.Equal(s.T(), userdto.CreateUserResponse{}, response)
}

func TestValidatePassword(t *testing.T) {
	tests := []struct {
		name      string
		password  string
		wantError bool
		errorType string
		errorMsg  string
	}{
		{
			name:      "valid password with all requirements",
			password:  "SecurePass123!",
			wantError: false,
		},
		{
			name:      "valid password with different special chars",
			password:  "MyP@ssw0rd#",
			wantError: false,
		},
		{
			name:      "password too short",
			password:  "Short1!",
			wantError: true,
			errorType: "length",
		},
		{
			name:      "password missing uppercase",
			password:  "lowercase123!",
			wantError: true,
			errorType: "uppercase",
		},
		{
			name:      "password missing lowercase",
			password:  "UPPERCASE123!",
			wantError: true,
			errorType: "lowercase",
		},
		{
			name:      "password missing number",
			password:  "NoNumberHere!",
			wantError: true,
			errorType: "number",
		},
		{
			name:      "password missing special character",
			password:  "NoSpecialChar123",
			wantError: true,
			errorType: "special",
		},
		{
			name:      "password missing multiple requirements",
			password:  "alllowercase",
			wantError: true,
			errorType: "multiple",
		},
		{
			name:      "empty password",
			password:  "",
			wantError: true,
			errorType: "length",
		},
		{
			name:      "password with unicode special chars",
			password:  "Pass123€",
			wantError: false,
		},
		{
			name:      "password exactly 8 characters",
			password:  "Pass123!",
			wantError: false,
		},
		{
			name:      "password with only spaces",
			password:  "        ",
			wantError: true,
			errorType: "multiple",
		},
		{
			name:      "password with various special characters",
			password:  "Test123@#$%^&*()",
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validatePassword(tt.password)

			if tt.wantError {
				assert.Error(t, err)
				var passwordErr *PasswordValidationError
				assert.True(t, errors.As(err, &passwordErr), "error should be PasswordValidationError")
				assert.NotEmpty(t, passwordErr.Message)

				if tt.errorType == "length" {
					assert.Contains(t, passwordErr.Message, "at least")
					assert.Contains(t, passwordErr.Message, "characters long")
				} else if tt.errorType == "uppercase" {
					assert.Contains(t, passwordErr.Message, "uppercase")
				} else if tt.errorType == "lowercase" {
					assert.Contains(t, passwordErr.Message, "lowercase")
				} else if tt.errorType == "number" {
					assert.Contains(t, passwordErr.Message, "number")
				} else if tt.errorType == "special" {
					assert.Contains(t, passwordErr.Message, "special")
				} else if tt.errorType == "multiple" {
					assert.Contains(t, passwordErr.Message, "must contain")
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func (s *CreateUserUseCaseTestSuite) TestCreateUser_PasswordValidation() {
	testCases := []struct {
		name        string
		password    string
		expectError bool
		errorType   string
	}{
		{
			name:        "password too short",
			password:    "Short1!",
			expectError: true,
			errorType:   "PasswordValidationError",
		},
		{
			name:        "password missing uppercase",
			password:    "lowercase123!",
			expectError: true,
			errorType:   "PasswordValidationError",
		},
		{
			name:        "password missing lowercase",
			password:    "UPPERCASE123!",
			expectError: true,
			errorType:   "PasswordValidationError",
		},
		{
			name:        "password missing number",
			password:    "NoNumberHere!",
			expectError: true,
			errorType:   "PasswordValidationError",
		},
		{
			name:        "password missing special character",
			password:    "NoSpecialChar123",
			expectError: true,
			errorType:   "PasswordValidationError",
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			req := s.validRequest
			req.Password = tc.password

			s.mockRepo.On("EmailExists", s.ctx, req.Email).Return(false, nil)
			s.mockRepo.On("UsernameExists", s.ctx, req.Username).Return(false, nil)

			response, err := s.useCase.Execute(s.ctx, req)

			if tc.expectError {
				assert.Error(s.T(), err)
				var passwordErr *PasswordValidationError
				assert.True(s.T(), errors.As(err, &passwordErr), "error should be PasswordValidationError")
				assert.Equal(s.T(), userdto.CreateUserResponse{}, response)
				s.mockRepo.AssertNotCalled(s.T(), "Create", s.ctx, mock.Anything)
			}
		})
	}
}

func (s *CreateUserUseCaseTestSuite) TestCreateUser_PasswordHashing() {
	s.mockRepo.On("EmailExists", s.ctx, s.validRequest.Email).Return(false, nil)
	s.mockRepo.On("UsernameExists", s.ctx, s.validRequest.Username).Return(false, nil)

	var capturedUser *userentity.User
	s.mockRepo.On("Create", s.ctx, mock.AnythingOfType("*user.User")).Run(func(args mock.Arguments) {
		capturedUser = args.Get(1).(*userentity.User)
	}).Return(s.validUser, nil)

	_, err := s.useCase.Execute(s.ctx, s.validRequest)

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), capturedUser)
	assert.NotEqual(s.T(), s.validRequest.Password, capturedUser.PasswordHash, "password should be hashed")
	assert.NotEmpty(s.T(), capturedUser.PasswordHash, "password hash should not be empty")
	assert.True(s.T(), len(capturedUser.PasswordHash) > 20, "bcrypt hash should be long enough")
}

func (s *CreateUserUseCaseTestSuite) TestCreateUser_ContextPropagation() {
	ctxWithValue := context.WithValue(s.ctx, "test-key", "test-value")

	s.mockRepo.On("EmailExists", ctxWithValue, s.validRequest.Email).Return(false, nil)
	s.mockRepo.On("UsernameExists", ctxWithValue, s.validRequest.Username).Return(false, nil)
	s.mockRepo.On("Create", ctxWithValue, mock.AnythingOfType("*user.User")).Return(s.validUser, nil)

	_, err := s.useCase.Execute(ctxWithValue, s.validRequest)

	assert.NoError(s.T(), err)
}

func (s *CreateUserUseCaseTestSuite) TestCreateUser_UserIDGeneration() {
	s.mockRepo.On("EmailExists", s.ctx, s.validRequest.Email).Return(false, nil)
	s.mockRepo.On("UsernameExists", s.ctx, s.validRequest.Username).Return(false, nil)

	var capturedUser *userentity.User
	s.mockRepo.On("Create", s.ctx, mock.AnythingOfType("*user.User")).Run(func(args mock.Arguments) {
		capturedUser = args.Get(1).(*userentity.User)
	}).Return(s.validUser, nil)

	_, err := s.useCase.Execute(s.ctx, s.validRequest)

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), capturedUser)
	assert.NotEqual(s.T(), uuid.Nil, capturedUser.ID, "user ID should be generated")
}

func (s *CreateUserUseCaseTestSuite) TestCreateUser_Timestamps() {
	s.mockRepo.On("EmailExists", s.ctx, s.validRequest.Email).Return(false, nil)
	s.mockRepo.On("UsernameExists", s.ctx, s.validRequest.Username).Return(false, nil)

	var capturedUser *userentity.User
	beforeTime := time.Now().Add(-time.Millisecond) 
	s.mockRepo.On("Create", s.ctx, mock.AnythingOfType("*user.User")).Run(func(args mock.Arguments) {
		capturedUser = args.Get(1).(*userentity.User)
	}).Return(s.validUser, nil)
	afterTime := time.Now().Add(time.Millisecond) 

	_, err := s.useCase.Execute(s.ctx, s.validRequest)

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), capturedUser)

	assert.WithinDuration(s.T(), beforeTime, capturedUser.CreatedAt, 2*time.Second,
		"created_at should be close to beforeTime")
	assert.WithinDuration(s.T(), afterTime, capturedUser.CreatedAt, 2*time.Second,
		"created_at should be close to afterTime")
	assert.WithinDuration(s.T(), beforeTime, capturedUser.UpdatedAt, 2*time.Second,
		"updated_at should be close to beforeTime")
	assert.WithinDuration(s.T(), afterTime, capturedUser.UpdatedAt, 2*time.Second,
		"updated_at should be close to afterTime")
	assert.Equal(s.T(), capturedUser.CreatedAt, capturedUser.UpdatedAt,
		"created_at and updated_at should be equal for new user")
}

func FuzzValidatePassword(f *testing.F) {
	seedCases := []string{
		"SecurePass123!",
		"Short1!",
		"lowercase123!",
		"UPPERCASE123!",
		"NoNumberHere!",
		"NoSpecialChar123",
		"",
		"Pass123!",
		"VeryLongPassword123!@#$%^&*()",
	}

	for _, seed := range seedCases {
		f.Add(seed)
	}

	f.Fuzz(func(t *testing.T, password string) {
		err := validatePassword(password)

		if err != nil {
			var passwordErr *PasswordValidationError
			assert.True(t, errors.As(err, &passwordErr), "error should be PasswordValidationError")
			assert.NotEmpty(t, passwordErr.Message)

			if len(password) < 8 {
				assert.Contains(t, passwordErr.Message, "at least")
				assert.Contains(t, passwordErr.Message, "characters long")
			} else {
				assert.True(t,
					strings.Contains(passwordErr.Message, "uppercase") ||
						strings.Contains(passwordErr.Message, "lowercase") ||
						strings.Contains(passwordErr.Message, "number") ||
						strings.Contains(passwordErr.Message, "special"),
					"error message should mention missing requirement")
			}
		} else {
			assert.GreaterOrEqual(t, len(password), 8, "valid password must be at least 8 characters")

			hasUpper := false
			hasLower := false
			hasNumber := false
			hasSpecial := false

			for _, char := range password {
				switch {
				case 'A' <= char && char <= 'Z':
					hasUpper = true
				case 'a' <= char && char <= 'z':
					hasLower = true
				case '0' <= char && char <= '9':
					hasNumber = true
				default:
					if strings.ContainsRune("!@#$%^&*()_+-=[]{}|;:,.<>?", char) {
						hasSpecial = true
					}
				}
			}

			assert.True(t, hasUpper, "valid password must have uppercase")
			assert.True(t, hasLower, "valid password must have lowercase")
			assert.True(t, hasNumber, "valid password must have number")
			assert.True(t, hasSpecial, "valid password must have special character")
		}
	})
}

func (s *CreateUserUseCaseTestSuite) TestCreateUser_AllErrorScenarios() {
	testCases := []struct {
		name             string
		setupMock        func()
		request          userdto.CreateUserRequest
		expectedError    error
		verifyError      func(*testing.T, error)
		shouldCallCreate bool
	}{
		{
			name: "EmailAlreadyExists - exact error type",
			setupMock: func() {
				s.mockRepo.On("EmailExists", s.ctx, "existing@example.com").Return(true, nil)
			},
			request: userdto.CreateUserRequest{
				Email:    "existing@example.com",
				Username: "testuser",
				Password: "SecurePass123!",
			},
			expectedError: ErrEmailAlreadyExists,
			verifyError: func(t *testing.T, err error) {
				assert.True(t, errors.Is(err, ErrEmailAlreadyExists), "should be ErrEmailAlreadyExists")
				assert.False(t, errors.Is(err, ErrUsernameAlreadyExists), "should not be ErrUsernameAlreadyExists")
			},
			shouldCallCreate: false,
		},
		{
			name: "UsernameAlreadyExists - exact error type",
			setupMock: func() {
				s.mockRepo.On("EmailExists", s.ctx, "test@example.com").Return(false, nil)
				s.mockRepo.On("UsernameExists", s.ctx, "existinguser").Return(true, nil)
			},
			request: userdto.CreateUserRequest{
				Email:    "test@example.com",
				Username: "existinguser",
				Password: "SecurePass123!",
			},
			expectedError: ErrUsernameAlreadyExists,
			verifyError: func(t *testing.T, err error) {
				assert.True(t, errors.Is(err, ErrUsernameAlreadyExists), "should be ErrUsernameAlreadyExists")
				assert.False(t, errors.Is(err, ErrEmailAlreadyExists), "should not be ErrEmailAlreadyExists")
			},
			shouldCallCreate: false,
		},
		{
			name: "PasswordValidationError - too short",
			setupMock: func() {
				s.mockRepo.On("EmailExists", s.ctx, "test@example.com").Return(false, nil)
				s.mockRepo.On("UsernameExists", s.ctx, "testuser").Return(false, nil)
			},
			request: userdto.CreateUserRequest{
				Email:    "test@example.com",
				Username: "testuser",
				Password: "Short1!",
			},
			expectedError: nil, 
			verifyError: func(t *testing.T, err error) {
				var passwordErr *PasswordValidationError
				assert.True(t, errors.As(err, &passwordErr), "should be PasswordValidationError")
				assert.Contains(t, passwordErr.Message, "at least")
				assert.Contains(t, passwordErr.Message, "characters long")
			},
			shouldCallCreate: false,
		},
		{
			name: "PasswordValidationError - missing uppercase",
			setupMock: func() {
				s.mockRepo.On("EmailExists", s.ctx, "test@example.com").Return(false, nil)
				s.mockRepo.On("UsernameExists", s.ctx, "testuser").Return(false, nil)
			},
			request: userdto.CreateUserRequest{
				Email:    "test@example.com",
				Username: "testuser",
				Password: "lowercase123!",
			},
			expectedError: nil,
			verifyError: func(t *testing.T, err error) {
				var passwordErr *PasswordValidationError
				assert.True(t, errors.As(err, &passwordErr))
				assert.Contains(t, passwordErr.Message, "uppercase")
			},
			shouldCallCreate: false,
		},
		{
			name: "EmailExists database error",
			setupMock: func() {
				s.mockRepo.On("EmailExists", s.ctx, "test@example.com").Return(false, errors.New("database connection failed"))
			},
			request: userdto.CreateUserRequest{
				Email:    "test@example.com",
				Username: "testuser",
				Password: "SecurePass123!",
			},
			expectedError: errors.New("database connection failed"),
			verifyError: func(t *testing.T, err error) {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "database connection failed")
			},
			shouldCallCreate: false,
		},
		{
			name: "UsernameExists database error",
			setupMock: func() {
				s.mockRepo.On("EmailExists", s.ctx, "test@example.com").Return(false, nil)
				s.mockRepo.On("UsernameExists", s.ctx, "testuser").Return(false, errors.New("database timeout"))
			},
			request: userdto.CreateUserRequest{
				Email:    "test@example.com",
				Username: "testuser",
				Password: "SecurePass123!",
			},
			expectedError: errors.New("database timeout"),
			verifyError: func(t *testing.T, err error) {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "database timeout")
			},
			shouldCallCreate: false,
		},
		{
			name: "Create database error",
			setupMock: func() {
				s.mockRepo.On("EmailExists", s.ctx, "test@example.com").Return(false, nil)
				s.mockRepo.On("UsernameExists", s.ctx, "testuser").Return(false, nil)
				s.mockRepo.On("Create", s.ctx, mock.AnythingOfType("*user.User")).Return(nil, errors.New("constraint violation"))
			},
			request: userdto.CreateUserRequest{
				Email:    "test@example.com",
				Username: "testuser",
				Password: "SecurePass123!",
			},
			expectedError: errors.New("constraint violation"),
			verifyError: func(t *testing.T, err error) {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "constraint violation")
			},
			shouldCallCreate: true,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			s.mockRepo.ExpectedCalls = nil
			s.mockRepo.Calls = nil

			tc.setupMock()

			response, err := s.useCase.Execute(s.ctx, tc.request)

			if tc.expectedError != nil || tc.verifyError != nil {
				assert.Error(s.T(), err)

				if tc.verifyError != nil {
					tc.verifyError(s.T(), err)
				} else if tc.expectedError != nil {
					assert.Equal(s.T(), tc.expectedError, err)
				}

				assert.Equal(s.T(), userdto.CreateUserResponse{}, response)
			}

			if !tc.shouldCallCreate {
				s.mockRepo.AssertNotCalled(s.T(), "Create", s.ctx, mock.Anything)
			}
		})
	}
}

func (s *CreateUserUseCaseTestSuite) TestCreateUser_PasswordValidation_Comprehensive() {
	testCases := []struct {
		name          string
		password      string
		shouldPass    bool
		expectedError string
	}{
		{
			name:          "password exactly 7 characters",
			password:      "Pass12!",
			shouldPass:    false,
			expectedError: "at least 8 characters",
		},
		{
			name:          "password exactly 8 characters - valid",
			password:      "Pass123!",
			shouldPass:    true,
			expectedError: "",
		},
		{
			name:          "password 72 characters (bcrypt max) - valid",
			password:      "A" + strings.Repeat("a", 69) + "1!", 
			shouldPass:    true,
			expectedError: "",
		},

		{
			name:          "only uppercase, number, special",
			password:      "UPPERCASE123!",
			shouldPass:    false,
			expectedError: "lowercase",
		},
		{
			name:          "only lowercase, number, special",
			password:      "lowercase123!",
			shouldPass:    false,
			expectedError: "uppercase",
		},
		{
			name:          "only uppercase, lowercase, special",
			password:      "NoNumbers!",
			shouldPass:    false,
			expectedError: "number",
		},
		{
			name:          "only uppercase, lowercase, number",
			password:      "NoSpecial123",
			shouldPass:    false,
			expectedError: "special",
		},

		{
			name:          "missing uppercase and number",
			password:      "lowercaseonly!",
			shouldPass:    false,
			expectedError: "uppercase",
		},
		{
			name:          "missing lowercase and special",
			password:      "UPPERCASE123",
			shouldPass:    false,
			expectedError: "lowercase",
		},

		{
			name:          "password with @ symbol",
			password:      "Pass123@",
			shouldPass:    true,
			expectedError: "",
		},
		{
			name:          "password with # symbol",
			password:      "Pass123#",
			shouldPass:    true,
			expectedError: "",
		},
		{
			name:          "password with $ symbol",
			password:      "Pass123$",
			shouldPass:    true,
			expectedError: "",
		},
		{
			name:          "password with % symbol",
			password:      "Pass123%",
			shouldPass:    true,
			expectedError: "",
		},
		{
			name:          "password with & symbol",
			password:      "Pass123&",
			shouldPass:    true,
			expectedError: "",
		},
		{
			name:          "password with * symbol",
			password:      "Pass123*",
			shouldPass:    true,
			expectedError: "",
		},
		{
			name:          "password with ( symbol",
			password:      "Pass123(",
			shouldPass:    true,
			expectedError: "",
		},
		{
			name:          "password with ) symbol",
			password:      "Pass123)",
			shouldPass:    true,
			expectedError: "",
		},
		{
			name:          "password with _ symbol",
			password:      "Pass123_",
			shouldPass:    true,
			expectedError: "",
		},
		{
			name:          "password with - symbol",
			password:      "Pass123-",
			shouldPass:    true,
			expectedError: "",
		},
		{
			name:          "password with + symbol",
			password:      "Pass123+",
			shouldPass:    true,
			expectedError: "",
		},
		{
			name:          "password with = symbol",
			password:      "Pass123=",
			shouldPass:    true,
			expectedError: "",
		},

		{
			name:          "password with spaces",
			password:      "Pass 123!",
			shouldPass:    true, 
			expectedError: "",
		},
		{
			name:          "password with unicode",
			password:      "Pass123€",
			shouldPass:    true,
			expectedError: "",
		},
		{
			name:          "password with emoji",
			password:      "Pass123😀",
			shouldPass:    true,
			expectedError: "",
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			s.mockRepo.ExpectedCalls = nil
			s.mockRepo.Calls = nil

			req := s.validRequest
			req.Password = tc.password

			s.mockRepo.On("EmailExists", s.ctx, req.Email).Return(false, nil)
			s.mockRepo.On("UsernameExists", s.ctx, req.Username).Return(false, nil)

			if tc.shouldPass {
				s.mockRepo.On("Create", s.ctx, mock.AnythingOfType("*user.User")).Return(s.validUser, nil)
			}

			response, err := s.useCase.Execute(s.ctx, req)

			if tc.shouldPass {
				assert.NoError(s.T(), err, "password should be valid: %s", tc.password)
				assert.NotEqual(s.T(), uuid.Nil, response.ID, "response should have valid ID")
			} else {
				assert.Error(s.T(), err, "password should be invalid: %s", tc.password)
				var passwordErr *PasswordValidationError
				assert.True(s.T(), errors.As(err, &passwordErr), "error should be PasswordValidationError")
				if tc.expectedError != "" {
					assert.Contains(s.T(), passwordErr.Message, tc.expectedError,
						"error message should contain: %s", tc.expectedError)
				}
				assert.Equal(s.T(), userdto.CreateUserResponse{}, response)
				s.mockRepo.AssertNotCalled(s.T(), "Create", s.ctx, mock.Anything)
			}
		})
	}
}

func (s *CreateUserUseCaseTestSuite) TestCreateUser_PasswordHashing_Verification() {
	testCases := []struct {
		name     string
		password string
	}{
		{"short password", "Pass123!"},
		{"medium password", "SecurePass123!"},
		{"long password", "VeryLongSecurePassword123!@#$%^&*()"},
		{"password with unicode", "Pass123€"},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			s.mockRepo.ExpectedCalls = nil
			s.mockRepo.Calls = nil

			req := s.validRequest
			req.Password = tc.password

			s.mockRepo.On("EmailExists", s.ctx, req.Email).Return(false, nil)
			s.mockRepo.On("UsernameExists", s.ctx, req.Username).Return(false, nil)

			var capturedUser *userentity.User
			s.mockRepo.On("Create", s.ctx, mock.AnythingOfType("*user.User")).Run(func(args mock.Arguments) {
				capturedUser = args.Get(1).(*userentity.User)
			}).Return(s.validUser, nil)

			_, err := s.useCase.Execute(s.ctx, req)

			assert.NoError(s.T(), err)
			assert.NotNil(s.T(), capturedUser)

			assert.NotEqual(s.T(), tc.password, capturedUser.PasswordHash,
				"password should be hashed, not stored as plain text")

			assert.True(s.T(),
				strings.HasPrefix(capturedUser.PasswordHash, "$2a$") ||
					strings.HasPrefix(capturedUser.PasswordHash, "$2b$") ||
					strings.HasPrefix(capturedUser.PasswordHash, "$2y$"),
				"password hash should be bcrypt format")

			assert.Equal(s.T(), 60, len(capturedUser.PasswordHash),
				"bcrypt hash should be 60 characters long")

			s.mockRepo.ExpectedCalls = nil
			s.mockRepo.Calls = nil

			req2 := req
			var capturedUser2 *userentity.User
			s.mockRepo.On("EmailExists", s.ctx, req2.Email).Return(false, nil)
			s.mockRepo.On("UsernameExists", s.ctx, req2.Username).Return(false, nil)
			s.mockRepo.On("Create", s.ctx, mock.AnythingOfType("*user.User")).Run(func(args mock.Arguments) {
				capturedUser2 = args.Get(1).(*userentity.User)
			}).Return(s.validUser, nil)

			_, err2 := s.useCase.Execute(s.ctx, req2)
			assert.NoError(s.T(), err2)
			assert.NotNil(s.T(), capturedUser2)

			assert.NotEqual(s.T(), capturedUser.PasswordHash, capturedUser2.PasswordHash,
				"same password should produce different hashes due to salt")
		})
	}
}

func (s *CreateUserUseCaseTestSuite) TestCreateUser_ConcurrentExecution() {
	const numGoroutines = 20

	results := make(chan struct {
		response userdto.CreateUserResponse
		err      error
	}, numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func(index int) {
			req := userdto.CreateUserRequest{
				Email:    "concurrent" + string(rune(index)) + "@example.com",
				Username: "concurrentuser" + string(rune(index)),
				Password: "SecurePass123!",
			}

			userID := uuid.New()
			now := time.Now()
			createdUser := &userentity.User{
				ID:           userID,
				Username:     req.Username,
				Email:        req.Email,
				PasswordHash: "$2a$10$hashed",
				CreatedAt:    now,
				UpdatedAt:    now,
			}

			s.mockRepo.On("EmailExists", s.ctx, req.Email).Return(false, nil).Maybe()
			s.mockRepo.On("UsernameExists", s.ctx, req.Username).Return(false, nil).Maybe()
			s.mockRepo.On("Create", s.ctx, mock.AnythingOfType("*user.User")).Return(createdUser, nil).Maybe()

			response, err := s.useCase.Execute(s.ctx, req)
			results <- struct {
				response userdto.CreateUserResponse
				err      error
			}{response, err}
		}(i)
	}

	var successCount, errorCount int
	for i := 0; i < numGoroutines; i++ {
		result := <-results
		if result.err == nil {
			successCount++
			assert.NotEqual(s.T(), uuid.Nil, result.response.ID)
		} else {
			errorCount++
		}
	}

	assert.Greater(s.T(), successCount, 0, "at least some requests should succeed")
}

func TestCreateUserUseCaseTestSuite(t *testing.T) {
	suite.Run(t, new(CreateUserUseCaseTestSuite))
}
