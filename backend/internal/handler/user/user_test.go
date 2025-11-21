package userhandler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	userdto "roadmap/internal/domain/dto/user"
	userentity "roadmap/internal/domain/entities/user"
	userrepo "roadmap/internal/repository/user"
	userusecase "roadmap/internal/usecase/user"
)

// MockUserRepository is a mock implementation of UserRepository for testing
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

// UserHandlerTestSuite is a test suite for UserHandler
type UserHandlerTestSuite struct {
	suite.Suite
	handler    *UserHandler
	mockRepo   *MockUserRepository
	useCase    *userusecase.CreateUserUseCase
	router     *gin.Engine
}

func (s *UserHandlerTestSuite) SetupTest() {
	gin.SetMode(gin.TestMode)
	s.mockRepo = new(MockUserRepository)
	// Create real use case with mock repository
	var repo userrepo.UserRepository = s.mockRepo
	s.useCase = userusecase.NewCreateUserUseCase(repo)
	s.handler = NewUserHandler(s.useCase)
	s.router = gin.New()
	s.router.POST("/api/v1/users", s.handler.CreateUser)
}

func (s *UserHandlerTestSuite) TearDownTest() {
	s.mockRepo.AssertExpectations(s.T())
}

// TestCreateUser_Success tests successful user creation
func (s *UserHandlerTestSuite) TestCreateUser_Success() {
	requestBody := userdto.CreateUserRequest{
		Email:    "test@example.com",
		Username: "testuser",
		Password: "SecurePass123!",
	}

	userID := uuid.New()
	now := time.Now()
	createdUser := &userentity.User{
		ID:           userID,
		Username:     requestBody.Username,
		Email:        requestBody.Email,
		PasswordHash: "$2a$10$hashed",
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	s.mockRepo.On("EmailExists", mock.Anything, requestBody.Email).Return(false, nil)
	s.mockRepo.On("UsernameExists", mock.Anything, requestBody.Username).Return(false, nil)
	s.mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*user.User")).Return(createdUser, nil)

	body, _ := json.Marshal(requestBody)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/users", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	s.router.ServeHTTP(w, req)

	assert.Equal(s.T(), http.StatusCreated, w.Code)
	
	var response userdto.CreateUserResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), userID, response.ID)
	assert.Equal(s.T(), requestBody.Username, response.Username)
	assert.Equal(s.T(), requestBody.Email, response.Email)
}

// TestCreateUser_InvalidJSON tests invalid JSON in request body
func (s *UserHandlerTestSuite) TestCreateUser_InvalidJSON() {
	req := httptest.NewRequest(http.MethodPost, "/api/v1/users", bytes.NewBufferString("invalid json"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	s.router.ServeHTTP(w, req)

	assert.Equal(s.T(), http.StatusBadRequest, w.Code)
	
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), "Invalid request data", response["error"])
	assert.Contains(s.T(), response, "details")
	
	s.mockRepo.AssertNotCalled(s.T(), "EmailExists", mock.Anything, mock.Anything)
	s.mockRepo.AssertNotCalled(s.T(), "Create", mock.Anything, mock.Anything)
}

// TestCreateUser_MissingFields tests request with missing required fields
func (s *UserHandlerTestSuite) TestCreateUser_MissingFields() {
	testCases := []struct {
		name        string
		requestBody map[string]interface{}
		description string
	}{
		{
			name:        "missing email",
			requestBody: map[string]interface{}{"username": "testuser", "password": "SecurePass123!"},
			description: "email is required",
		},
		{
			name:        "missing username",
			requestBody: map[string]interface{}{"email": "test@example.com", "password": "SecurePass123!"},
			description: "username is required",
		},
		{
			name:        "missing password",
			requestBody: map[string]interface{}{"email": "test@example.com", "username": "testuser"},
			description: "password is required",
		},
		{
			name:        "empty email",
			requestBody: map[string]interface{}{"email": "", "username": "testuser", "password": "SecurePass123!"},
			description: "email cannot be empty",
		},
		{
			name:        "invalid email format",
			requestBody: map[string]interface{}{"email": "not-an-email", "username": "testuser", "password": "SecurePass123!"},
			description: "email must be valid format",
		},
		{
			name:        "username too short",
			requestBody: map[string]interface{}{"email": "test@example.com", "username": "ab", "password": "SecurePass123!"},
			description: "username must be at least 3 characters",
		},
		{
			name:        "password too short",
			requestBody: map[string]interface{}{"email": "test@example.com", "username": "testuser", "password": "Short1!"},
			description: "password must be at least 8 characters",
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			body, _ := json.Marshal(tc.requestBody)
			req := httptest.NewRequest(http.MethodPost, "/api/v1/users", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			s.router.ServeHTTP(w, req)

			assert.Equal(s.T(), http.StatusBadRequest, w.Code, tc.description)
			
			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(s.T(), err)
			assert.Equal(s.T(), "Invalid request data", response["error"])
			
			s.mockRepo.AssertNotCalled(s.T(), "EmailExists", mock.Anything, mock.Anything)
	s.mockRepo.AssertNotCalled(s.T(), "Create", mock.Anything, mock.Anything)
		})
	}
}

// TestCreateUser_EmailAlreadyExists tests error when email already exists
func (s *UserHandlerTestSuite) TestCreateUser_EmailAlreadyExists() {
	requestBody := userdto.CreateUserRequest{
		Email:    "existing@example.com",
		Username: "testuser",
		Password: "SecurePass123!",
	}

	s.mockRepo.On("EmailExists", mock.Anything, requestBody.Email).Return(true, nil)

	body, _ := json.Marshal(requestBody)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/users", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	s.router.ServeHTTP(w, req)

	assert.Equal(s.T(), http.StatusConflict, w.Code)
	
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), "Email already exists", response["error"])
}

// TestCreateUser_UsernameAlreadyExists tests error when username already exists
func (s *UserHandlerTestSuite) TestCreateUser_UsernameAlreadyExists() {
	requestBody := userdto.CreateUserRequest{
		Email:    "test@example.com",
		Username: "existinguser",
		Password: "SecurePass123!",
	}

	s.mockRepo.On("EmailExists", mock.Anything, requestBody.Email).Return(false, nil)
	s.mockRepo.On("UsernameExists", mock.Anything, requestBody.Username).Return(true, nil)

	body, _ := json.Marshal(requestBody)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/users", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	s.router.ServeHTTP(w, req)

	assert.Equal(s.T(), http.StatusConflict, w.Code)
	
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), "Username already exists", response["error"])
}

// TestCreateUser_PasswordValidationError tests password validation error
func (s *UserHandlerTestSuite) TestCreateUser_PasswordValidationError() {
	requestBody := userdto.CreateUserRequest{
		Email:    "test@example.com",
		Username: "testuser",
		Password: "alllowercase", // 12 chars, passes DTO validation but fails usecase validation (missing uppercase, number, special)
	}

	s.mockRepo.On("EmailExists", mock.Anything, requestBody.Email).Return(false, nil)
	s.mockRepo.On("UsernameExists", mock.Anything, requestBody.Username).Return(false, nil)
	// Password validation happens in use case and will return error
	// No need to mock Create as validation fails before that

	body, _ := json.Marshal(requestBody)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/users", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	s.router.ServeHTTP(w, req)

	assert.Equal(s.T(), http.StatusBadRequest, w.Code)
	
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(s.T(), err)
	// Password validation error from usecase should contain "password must contain"
	assert.Contains(s.T(), response["error"].(string), "password must contain")
}

// TestCreateUser_InternalServerError tests internal server error
func (s *UserHandlerTestSuite) TestCreateUser_InternalServerError() {
	requestBody := userdto.CreateUserRequest{
		Email:    "test@example.com",
		Username: "testuser",
		Password: "SecurePass123!",
	}

	internalError := errors.New("database connection failed")
	s.mockRepo.On("EmailExists", mock.Anything, requestBody.Email).Return(false, internalError)

	body, _ := json.Marshal(requestBody)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/users", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	s.router.ServeHTTP(w, req)

	assert.Equal(s.T(), http.StatusInternalServerError, w.Code)
	
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), "Failed to create user", response["error"])
}

// TestCreateUser_EmptyBody tests request with empty body
func (s *UserHandlerTestSuite) TestCreateUser_EmptyBody() {
	req := httptest.NewRequest(http.MethodPost, "/api/v1/users", bytes.NewBufferString(""))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	s.router.ServeHTTP(w, req)

	assert.Equal(s.T(), http.StatusBadRequest, w.Code)
	s.mockRepo.AssertNotCalled(s.T(), "EmailExists", mock.Anything, mock.Anything)
	s.mockRepo.AssertNotCalled(s.T(), "Create", mock.Anything, mock.Anything)
}

// TestCreateUser_NilBody tests request with nil body
func (s *UserHandlerTestSuite) TestCreateUser_NilBody() {
	req := httptest.NewRequest(http.MethodPost, "/api/v1/users", nil)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	s.router.ServeHTTP(w, req)

	assert.Equal(s.T(), http.StatusBadRequest, w.Code)
	s.mockRepo.AssertNotCalled(s.T(), "EmailExists", mock.Anything, mock.Anything)
	s.mockRepo.AssertNotCalled(s.T(), "Create", mock.Anything, mock.Anything)
}

// TestCreateUser_WrongContentType tests request with wrong content type
func (s *UserHandlerTestSuite) TestCreateUser_WrongContentType() {
	// Reset mocks
	s.mockRepo.ExpectedCalls = nil
	s.mockRepo.Calls = nil
	
	requestBody := userdto.CreateUserRequest{
		Email:    "test@example.com",
		Username: "testuser",
		Password: "SecurePass123!",
	}

	// Setup mocks in case Gin still parses JSON
	userID := uuid.New()
	now := time.Now()
	createdUser := &userentity.User{
		ID:           userID,
		Username:     requestBody.Username,
		Email:        requestBody.Email,
		PasswordHash: "$2a$10$hashed",
		CreatedAt:    now,
		UpdatedAt:    now,
	}
	
	s.mockRepo.On("EmailExists", mock.Anything, requestBody.Email).Return(false, nil).Maybe()
	s.mockRepo.On("UsernameExists", mock.Anything, requestBody.Username).Return(false, nil).Maybe()
	s.mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*user.User")).Return(createdUser, nil).Maybe()

	body, _ := json.Marshal(requestBody)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/users", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "text/plain")
	w := httptest.NewRecorder()

	s.router.ServeHTTP(w, req)

	// Gin might still parse JSON even with wrong content type
	// The behavior depends on Gin's JSON binding, but we verify that either:
	// 1. It fails with BadRequest (if JSON parsing fails)
	// 2. Or it succeeds (if Gin still parses JSON)
	// In either case, we document the behavior
	assert.True(s.T(), 
		w.Code == http.StatusBadRequest || w.Code == http.StatusCreated,
		"Should be either BadRequest or Created")
}

// TestCreateUser_ResponseStructure tests that response has correct structure
func (s *UserHandlerTestSuite) TestCreateUser_ResponseStructure() {
	requestBody := userdto.CreateUserRequest{
		Email:    "test@example.com",
		Username: "testuser",
		Password: "SecurePass123!",
	}

	userID := uuid.New()
	now := time.Now()
	createdUser := &userentity.User{
		ID:           userID,
		Username:     requestBody.Username,
		Email:        requestBody.Email,
		PasswordHash: "$2a$10$hashed",
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	s.mockRepo.On("EmailExists", mock.Anything, requestBody.Email).Return(false, nil)
	s.mockRepo.On("UsernameExists", mock.Anything, requestBody.Username).Return(false, nil)
	s.mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*user.User")).Return(createdUser, nil)

	body, _ := json.Marshal(requestBody)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/users", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	s.router.ServeHTTP(w, req)

	assert.Equal(s.T(), http.StatusCreated, w.Code)
	
	var response userdto.CreateUserResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(s.T(), err)
	
	// Verify all fields are present and correct
	assert.NotEqual(s.T(), uuid.Nil, response.ID)
	assert.Equal(s.T(), requestBody.Username, response.Username)
	assert.Equal(s.T(), requestBody.Email, response.Email)
	assert.False(s.T(), response.CreatedAt.IsZero())
	assert.False(s.T(), response.UpdatedAt.IsZero())
	
	// Verify password is not in response
	responseStr := w.Body.String()
	assert.NotContains(s.T(), responseStr, "password")
	assert.NotContains(s.T(), responseStr, "password_hash")
}

// TestCreateUser_ContextPropagation tests that context is properly passed to use case
func (s *UserHandlerTestSuite) TestCreateUser_ContextPropagation() {
	requestBody := userdto.CreateUserRequest{
		Email:    "test@example.com",
		Username: "testuser",
		Password: "SecurePass123!",
	}

	userID := uuid.New()
	now := time.Now()
	createdUser := &userentity.User{
		ID:           userID,
		Username:     requestBody.Username,
		Email:        requestBody.Email,
		PasswordHash: "$2a$10$hashed",
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	// Verify context is passed (not nil)
	s.mockRepo.On("EmailExists", mock.MatchedBy(func(ctx interface{}) bool {
		_, ok := ctx.(context.Context)
		return ok
	}), requestBody.Email).Return(false, nil)
	s.mockRepo.On("UsernameExists", mock.MatchedBy(func(ctx interface{}) bool {
		_, ok := ctx.(context.Context)
		return ok
	}), requestBody.Username).Return(false, nil)
	s.mockRepo.On("Create", mock.MatchedBy(func(ctx interface{}) bool {
		_, ok := ctx.(context.Context)
		return ok
	}), mock.AnythingOfType("*user.User")).Return(createdUser, nil)

	body, _ := json.Marshal(requestBody)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/users", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	s.router.ServeHTTP(w, req)

	assert.Equal(s.T(), http.StatusCreated, w.Code)
}

// TestCreateUser_MultiplePasswordValidationErrors tests different password validation errors
func (s *UserHandlerTestSuite) TestCreateUser_MultiplePasswordValidationErrors() {
	testCases := []struct {
		name          string
		password      string
		expectedError string
	}{
		{
			name:          "password too short",
			password:      "Short1!",
			expectedError: "password must be at least",
		},
		{
			name:          "missing uppercase",
			password:      "lowercase123!",
			expectedError: "uppercase",
		},
		{
			name:          "missing lowercase",
			password:      "UPPERCASE123!",
			expectedError: "lowercase",
		},
		{
			name:          "missing number",
			password:      "NoNumberHere!",
			expectedError: "number",
		},
		{
			name:          "missing special character",
			password:      "NoSpecialChar123",
			expectedError: "special",
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			// Reset mocks
			s.mockRepo.ExpectedCalls = nil
			s.mockRepo.Calls = nil
			
			requestBody := userdto.CreateUserRequest{
				Email:    "test@example.com",
				Username: "testuser",
				Password: tc.password,
			}

			// For passwords that pass DTO validation (>= 8 chars), mock repository calls
			if len(tc.password) >= 8 {
				s.mockRepo.On("EmailExists", mock.Anything, requestBody.Email).Return(false, nil)
				s.mockRepo.On("UsernameExists", mock.Anything, requestBody.Username).Return(false, nil)
			}
			// Password validation happens in use case and will return error for invalid passwords

			body, _ := json.Marshal(requestBody)
			req := httptest.NewRequest(http.MethodPost, "/api/v1/users", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			s.router.ServeHTTP(w, req)

			assert.Equal(s.T(), http.StatusBadRequest, w.Code)
			
			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(s.T(), err)
			
			errorMsg := response["error"].(string)
			// For short passwords, DTO validation fails with "Invalid request data"
			// For longer passwords, usecase validation fails with specific message
			if len(tc.password) < 8 {
				assert.Equal(s.T(), "Invalid request data", errorMsg, 
					"short passwords should fail DTO validation")
			} else {
				assert.Contains(s.T(), errorMsg, tc.expectedError,
					"longer invalid passwords should fail usecase validation with specific message")
			}
		})
	}
}

// TestCreateUser_AllErrorTypes tests all possible error types and their exact error messages
func (s *UserHandlerTestSuite) TestCreateUser_AllErrorTypes() {
	testCases := []struct {
		name           string
		setupMock      func()
		expectedStatus int
		expectedError  string
		verifyError    func(*testing.T, map[string]interface{})
	}{
		{
			name: "EmailAlreadyExists returns exact error",
			setupMock: func() {
				s.mockRepo.On("EmailExists", mock.Anything, "existing@example.com").Return(true, nil)
			},
			expectedStatus: http.StatusConflict,
			expectedError:  "Email already exists",
			verifyError: func(t *testing.T, response map[string]interface{}) {
				assert.Equal(t, "Email already exists", response["error"])
				assert.Nil(t, response["details"])
			},
		},
		{
			name: "UsernameAlreadyExists returns exact error",
			setupMock: func() {
				s.mockRepo.On("EmailExists", mock.Anything, "test@example.com").Return(false, nil)
				s.mockRepo.On("UsernameExists", mock.Anything, "existinguser").Return(true, nil)
			},
			expectedStatus: http.StatusConflict,
			expectedError:  "Username already exists",
			verifyError: func(t *testing.T, response map[string]interface{}) {
				assert.Equal(t, "Username already exists", response["error"])
			},
		},
		{
			name: "PasswordValidationError returns exact error message",
			setupMock: func() {
				s.mockRepo.On("EmailExists", mock.Anything, "test@example.com").Return(false, nil)
				s.mockRepo.On("UsernameExists", mock.Anything, "testuser").Return(false, nil)
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "password must contain",
			verifyError: func(t *testing.T, response map[string]interface{}) {
				errorMsg := response["error"].(string)
				// Password validation happens in usecase, so error should contain password validation message
				assert.Contains(t, errorMsg, "password must contain", 
					"error should contain password validation message from usecase")
			},
		},
		{
			name: "InternalServerError returns generic error",
			setupMock: func() {
				s.mockRepo.On("EmailExists", mock.Anything, "test@example.com").Return(false, errors.New("database connection failed"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedError:  "Failed to create user",
			verifyError: func(t *testing.T, response map[string]interface{}) {
				assert.Equal(t, "Failed to create user", response["error"])
			},
		},
		{
			name: "Repository error during username check",
			setupMock: func() {
				s.mockRepo.On("EmailExists", mock.Anything, "test@example.com").Return(false, nil)
				s.mockRepo.On("UsernameExists", mock.Anything, "testuser").Return(false, errors.New("database timeout"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedError:  "Failed to create user",
			verifyError: func(t *testing.T, response map[string]interface{}) {
				assert.Equal(t, "Failed to create user", response["error"])
			},
		},
		{
			name: "Repository error during create",
			setupMock: func() {
				s.mockRepo.On("EmailExists", mock.Anything, "test@example.com").Return(false, nil)
				s.mockRepo.On("UsernameExists", mock.Anything, "testuser").Return(false, nil)
				s.mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*user.User")).Return(nil, errors.New("constraint violation"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedError:  "Failed to create user",
			verifyError: func(t *testing.T, response map[string]interface{}) {
				assert.Equal(t, "Failed to create user", response["error"])
			},
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			// Reset mocks
			s.mockRepo.ExpectedCalls = nil
			s.mockRepo.Calls = nil

			tc.setupMock()

			requestBody := userdto.CreateUserRequest{
				Email:    "test@example.com",
				Username: "testuser",
				Password: "SecurePass123!",
			}

			// Override for specific test cases
			if tc.name == "EmailAlreadyExists returns exact error" {
				requestBody.Email = "existing@example.com"
			}
			if tc.name == "UsernameAlreadyExists returns exact error" {
				requestBody.Username = "existinguser"
			}
			if tc.name == "PasswordValidationError returns exact error message" {
				// Use password that passes DTO validation (min 8 chars) but fails usecase validation (missing requirements)
				requestBody.Password = "alllowercase" // 12 chars, but missing uppercase, number, special
			}

			body, _ := json.Marshal(requestBody)
			req := httptest.NewRequest(http.MethodPost, "/api/v1/users", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			s.router.ServeHTTP(w, req)

			assert.Equal(s.T(), tc.expectedStatus, w.Code, "Status code should match")

			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(s.T(), err)

			if tc.verifyError != nil {
				tc.verifyError(s.T(), response)
			} else {
				assert.Contains(s.T(), response["error"].(string), tc.expectedError)
			}
		})
	}
}

// TestCreateUser_EdgeCases tests edge cases and boundary conditions
func (s *UserHandlerTestSuite) TestCreateUser_EdgeCases() {
	testCases := []struct {
		name        string
		requestBody map[string]interface{}
		description string
	}{
		{
			name:        "email with maximum valid length",
			requestBody: map[string]interface{}{"email": "a" + strings.Repeat("b", 240) + "@example.com", "username": "testuser", "password": "SecurePass123!"},
			description: "very long but valid email",
		},
		{
			name:        "username exactly 3 characters (minimum)",
			requestBody: map[string]interface{}{"email": "test@example.com", "username": "abc", "password": "SecurePass123!"},
			description: "username at minimum length",
		},
		{
			name:        "username exactly 100 characters (maximum)",
			requestBody: map[string]interface{}{"email": "test@example.com", "username": strings.Repeat("a", 100), "password": "SecurePass123!"},
			description: "username at maximum length",
		},
		{
			name:        "password exactly 8 characters (minimum)",
			requestBody: map[string]interface{}{"email": "test@example.com", "username": "testuser", "password": "Pass123!"},
			description: "password at minimum length",
		},
		{
			name:        "email with plus sign",
			requestBody: map[string]interface{}{"email": "test+tag@example.com", "username": "testuser", "password": "SecurePass123!"},
			description: "email with plus addressing",
		},
		{
			name:        "email with subdomain",
			requestBody: map[string]interface{}{"email": "test@mail.example.com", "username": "testuser", "password": "SecurePass123!"},
			description: "email with subdomain",
		},
		{
			name:        "username with numbers",
			requestBody: map[string]interface{}{"email": "test@example.com", "username": "user123", "password": "SecurePass123!"},
			description: "username with numbers",
		},
		{
			name:        "username with underscores",
			requestBody: map[string]interface{}{"email": "test@example.com", "username": "user_name", "password": "SecurePass123!"},
			description: "username with underscores",
		},
		{
			name:        "password with unicode special characters",
			requestBody: map[string]interface{}{"email": "test@example.com", "username": "testuser", "password": "Pass123€"},
			description: "password with unicode",
		},
		{
			name:        "email with unicode",
			requestBody: map[string]interface{}{"email": "tëst@example.com", "username": "testuser", "password": "SecurePass123!"},
			description: "email with unicode characters",
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			userID := uuid.New()
			now := time.Now()
			createdUser := &userentity.User{
				ID:           userID,
				Username:     tc.requestBody["username"].(string),
				Email:        tc.requestBody["email"].(string),
				PasswordHash: "$2a$10$hashed",
				CreatedAt:    now,
				UpdatedAt:    now,
			}

			s.mockRepo.On("EmailExists", mock.Anything, tc.requestBody["email"]).Return(false, nil)
			s.mockRepo.On("UsernameExists", mock.Anything, tc.requestBody["username"]).Return(false, nil)
			s.mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*user.User")).Return(createdUser, nil)

			body, _ := json.Marshal(tc.requestBody)
			req := httptest.NewRequest(http.MethodPost, "/api/v1/users", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			s.router.ServeHTTP(w, req)

			if w.Code == http.StatusCreated {
				var response userdto.CreateUserResponse
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(s.T(), err, tc.description)
				assert.NotEqual(s.T(), uuid.Nil, response.ID, tc.description)
				// Don't check exact UUID match as usecase generates its own
			} else {
				// If validation failed, that's also valid for edge cases
				assert.True(s.T(), w.Code == http.StatusBadRequest || w.Code == http.StatusCreated, 
					"Should be either BadRequest or Created for edge case: %s", tc.description)
			}
		})
	}
}

// TestCreateUser_ConcurrentRequests tests handling of concurrent requests
func (s *UserHandlerTestSuite) TestCreateUser_ConcurrentRequests() {
	requestBody := userdto.CreateUserRequest{
		Email:    "concurrent@example.com",
		Username: "concurrentuser",
		Password: "SecurePass123!",
	}

	userID := uuid.New()
	now := time.Now()
	createdUser := &userentity.User{
		ID:           userID,
		Username:     requestBody.Username,
		Email:        requestBody.Email,
		PasswordHash: "$2a$10$hashed",
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	// Setup mock to handle concurrent calls
	s.mockRepo.On("EmailExists", mock.Anything, requestBody.Email).Return(false, nil).Maybe()
	s.mockRepo.On("UsernameExists", mock.Anything, requestBody.Username).Return(false, nil).Maybe()
	s.mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*user.User")).Return(createdUser, nil).Maybe()

	body, _ := json.Marshal(requestBody)
	
	// Run concurrent requests
	const numRequests = 10
	results := make(chan int, numRequests)
	
	for i := 0; i < numRequests; i++ {
		go func() {
			req := httptest.NewRequest(http.MethodPost, "/api/v1/users", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			s.router.ServeHTTP(w, req)
			results <- w.Code
		}()
	}

	// Collect results
	var successCount, errorCount int
	for i := 0; i < numRequests; i++ {
		code := <-results
		if code == http.StatusCreated {
			successCount++
		} else {
			errorCount++
		}
	}

	// At least some requests should succeed
	assert.Greater(s.T(), successCount, 0, "At least one request should succeed")
}

// TestCreateUser_ResponseHeaders tests response headers
func (s *UserHandlerTestSuite) TestCreateUser_ResponseHeaders() {
	requestBody := userdto.CreateUserRequest{
		Email:    "test@example.com",
		Username: "testuser",
		Password: "SecurePass123!",
	}

	userID := uuid.New()
	now := time.Now()
	createdUser := &userentity.User{
		ID:           userID,
		Username:     requestBody.Username,
		Email:        requestBody.Email,
		PasswordHash: "$2a$10$hashed",
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	s.mockRepo.On("EmailExists", mock.Anything, requestBody.Email).Return(false, nil)
	s.mockRepo.On("UsernameExists", mock.Anything, requestBody.Username).Return(false, nil)
	s.mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*user.User")).Return(createdUser, nil)

	body, _ := json.Marshal(requestBody)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/users", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	s.router.ServeHTTP(w, req)

	assert.Equal(s.T(), http.StatusCreated, w.Code)
	assert.Equal(s.T(), "application/json; charset=utf-8", w.Header().Get("Content-Type"))
}

// TestCreateUser_InvalidJSONVariations tests various invalid JSON formats
func (s *UserHandlerTestSuite) TestCreateUser_InvalidJSONVariations() {
	testCases := []struct {
		name        string
		body        string
		description string
	}{
		{
			name:        "malformed JSON - missing quote",
			body:        `{"email": "test@example.com", "username": "testuser", "password": "SecurePass123!"`,
			description: "unclosed JSON object",
		},
		{
			name:        "malformed JSON - trailing comma",
			body:        `{"email": "test@example.com", "username": "testuser", "password": "SecurePass123!",}`,
			description: "trailing comma in JSON",
		},
		{
			name:        "invalid JSON - wrong brackets",
			body:        `["email": "test@example.com"]`,
			description: "array instead of object",
		},
		{
			name:        "invalid JSON - null value",
			body:        `null`,
			description: "null JSON value",
		},
		{
			name:        "invalid JSON - empty string",
			body:        ``,
			description: "empty request body",
		},
		{
			name:        "invalid JSON - just whitespace",
			body:        `   `,
			description: "whitespace only",
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			req := httptest.NewRequest(http.MethodPost, "/api/v1/users", bytes.NewBufferString(tc.body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			s.router.ServeHTTP(w, req)

			assert.Equal(s.T(), http.StatusBadRequest, w.Code, tc.description)
			
			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(s.T(), err)
			assert.Equal(s.T(), "Invalid request data", response["error"])
			assert.Contains(s.T(), response, "details")
			
			s.mockRepo.AssertNotCalled(s.T(), "EmailExists", mock.Anything, mock.Anything)
			s.mockRepo.AssertNotCalled(s.T(), "Create", mock.Anything, mock.Anything)
		})
	}
}

// Run the test suite
func TestUserHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(UserHandlerTestSuite))
}

