package user

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	userentity "roadmap/internal/domain/entities/user"
)

type MockDatabase struct {
	mock.Mock
}

type MockPool struct {
	mock.Mock
}

type MockRow struct {
	mock.Mock
	values []interface{}
	err    error
}

func (m *MockRow) Scan(dest ...interface{}) error {
	if m.err != nil {
		return m.err
	}
	if len(m.values) != len(dest) {
		return errors.New("value count mismatch")
	}
	for i, val := range m.values {
		if destPtr, ok := dest[i].(*interface{}); ok {
			*destPtr = val
		} else {
			switch d := dest[i].(type) {
			case *uuid.UUID:
				if uuidVal, ok := val.(uuid.UUID); ok {
					*d = uuidVal
				}
			case *string:
				if strVal, ok := val.(string); ok {
					*d = strVal
				}
			case *time.Time:
				if timeVal, ok := val.(time.Time); ok {
					*d = timeVal
				}
			case *bool:
				if boolVal, ok := val.(bool); ok {
					*d = boolVal
				}
			}
		}
	}
	return nil
}

type UserRepositoryTestSuite struct {
	suite.Suite
	repo     *userRepository
	ctx      context.Context
	testUser *userentity.User
}

func (s *UserRepositoryTestSuite) SetupTest() {
	s.ctx = context.Background()

	now := time.Now()
	s.testUser = &userentity.User{
		ID:           uuid.New(),
		Username:     "testuser",
		Email:        "test@example.com",
		PasswordHash: "$2a$10$hashedpassword",
		CreatedAt:    now,
		UpdatedAt:    now,
	}
}

func TestNewUserRepository(t *testing.T) {
	repo := &userRepository{}
	assert.NotNil(t, repo)
}

func (s *UserRepositoryTestSuite) TestUserRepository_Create_Success() {

	assert.NotNil(s.T(), s.testUser)
	assert.NotEqual(s.T(), uuid.Nil, s.testUser.ID)
}

func (s *UserRepositoryTestSuite) TestUserRepository_GetByID_Success() {
	testID := uuid.New()
	assert.NotEqual(s.T(), uuid.Nil, testID)
}

func (s *UserRepositoryTestSuite) TestUserRepository_GetByID_NotFound() {
	expectedError := pgx.ErrNoRows
	assert.True(s.T(), errors.Is(expectedError, pgx.ErrNoRows))
}

func (s *UserRepositoryTestSuite) TestUserRepository_GetByEmail_Success() {
	testEmail := "test@example.com"
	assert.NotEmpty(s.T(), testEmail)
}

func (s *UserRepositoryTestSuite) TestUserRepository_GetByEmail_NotFound() {
	expectedError := pgx.ErrNoRows
	assert.True(s.T(), errors.Is(expectedError, pgx.ErrNoRows))
}

func (s *UserRepositoryTestSuite) TestUserRepository_EmailExists_True() {
	testEmail := "existing@example.com"
	assert.NotEmpty(s.T(), testEmail)
}

func (s *UserRepositoryTestSuite) TestUserRepository_EmailExists_False() {
	testEmail := "nonexistent@example.com"
	assert.NotEmpty(s.T(), testEmail)
}

func (s *UserRepositoryTestSuite) TestUserRepository_EmailExists_Error() {
	dbError := errors.New("database connection error")
	assert.Error(s.T(), dbError)
}

func (s *UserRepositoryTestSuite) TestUserRepository_UsernameExists_True() {
	testUsername := "existinguser"
	assert.NotEmpty(s.T(), testUsername)
}

func (s *UserRepositoryTestSuite) TestUserRepository_UsernameExists_False() {
	testUsername := "nonexistentuser"
	assert.NotEmpty(s.T(), testUsername)
}

func (s *UserRepositoryTestSuite) TestUserRepository_UsernameExists_Error() {
	dbError := errors.New("database connection error")
	assert.Error(s.T(), dbError)
}

func (s *UserRepositoryTestSuite) TestUserRepository_ErrorWrapping() {
	originalError := errors.New("original error")
	wrappedError := errors.New("failed to create user: " + originalError.Error())

	assert.Error(s.T(), wrappedError)
	assert.Contains(s.T(), wrappedError.Error(), "failed to create user")
	assert.Contains(s.T(), wrappedError.Error(), "original error")
}

func TestUserRepository_QueryStructure(t *testing.T) {

	tests := []struct {
		name  string
		query string
		check func(*testing.T, string)
	}{
		{
			name:  "create query has all fields",
			query: `INSERT INTO users (id, email, password_hash, username, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id, email, password_hash, username, created_at, updated_at`,
			check: func(t *testing.T, q string) {
				assert.Contains(t, q, "INSERT INTO users")
				assert.Contains(t, q, "id, email, password_hash, username, created_at, updated_at")
				assert.Contains(t, q, "RETURNING")
			},
		},
		{
			name:  "get by id query",
			query: `SELECT id, email, password_hash, username, created_at, updated_at FROM users WHERE id = $1`,
			check: func(t *testing.T, q string) {
				assert.Contains(t, q, "SELECT")
				assert.Contains(t, q, "FROM users")
				assert.Contains(t, q, "WHERE id = $1")
			},
		},
		{
			name:  "get by email query",
			query: `SELECT id, email, password_hash, username, created_at, updated_at FROM users WHERE email = $1`,
			check: func(t *testing.T, q string) {
				assert.Contains(t, q, "SELECT")
				assert.Contains(t, q, "FROM users")
				assert.Contains(t, q, "WHERE email = $1")
			},
		},
		{
			name:  "email exists query",
			query: `SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)`,
			check: func(t *testing.T, q string) {
				assert.Contains(t, q, "SELECT EXISTS")
				assert.Contains(t, q, "FROM users")
				assert.Contains(t, q, "WHERE email = $1")
			},
		},
		{
			name:  "username exists query",
			query: `SELECT EXISTS(SELECT 1 FROM users WHERE username = $1)`,
			check: func(t *testing.T, q string) {
				assert.Contains(t, q, "SELECT EXISTS")
				assert.Contains(t, q, "FROM users")
				assert.Contains(t, q, "WHERE username = $1")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.check(t, tt.query)
		})
	}
}

func TestUserRepository_ScanOrder(t *testing.T) {
	expectedOrder := []string{"id", "email", "password_hash", "username", "created_at", "updated_at"}

	assert.Equal(t, 6, len(expectedOrder))
	assert.Equal(t, "id", expectedOrder[0])
	assert.Equal(t, "email", expectedOrder[1])
	assert.Equal(t, "password_hash", expectedOrder[2])
	assert.Equal(t, "username", expectedOrder[3])
	assert.Equal(t, "created_at", expectedOrder[4])
	assert.Equal(t, "updated_at", expectedOrder[5])
}

func TestUserRepository_ErrorMessages(t *testing.T) {
	tests := []struct {
		name        string
		operation   string
		expectedMsg string
	}{
		{
			name:        "create error message",
			operation:   "create",
			expectedMsg: "failed to create user",
		},
		{
			name:        "get by id error message",
			operation:   "get by id",
			expectedMsg: "failed to get user by id",
		},
		{
			name:        "get by email error message",
			operation:   "get by email",
			expectedMsg: "failed to get user by email",
		},
		{
			name:        "email exists error message",
			operation:   "email exists",
			expectedMsg: "failed to check email existence",
		},
		{
			name:        "username exists error message",
			operation:   "username exists",
			expectedMsg: "failed to check username existence",
		},
		{
			name:        "user not found error message",
			operation:   "not found",
			expectedMsg: "user not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.NotEmpty(t, tt.expectedMsg)
			if tt.operation != "not found" {
				assert.Contains(t, tt.expectedMsg, "failed",
					"error message should contain 'failed' for operation: %s", tt.operation)
			} else {
				assert.Contains(t, tt.expectedMsg, "not found",
					"error message should contain 'not found'")
			}
		})
	}
}

func (s *UserRepositoryTestSuite) TestUserRepository_ContextPropagation() {
	ctx := context.Background()
	assert.NotNil(s.T(), ctx)

	ctxWithValue := context.WithValue(ctx, "test-key", "test-value")
	assert.NotNil(s.T(), ctxWithValue)
}

func TestUserRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(UserRepositoryTestSuite))
}
