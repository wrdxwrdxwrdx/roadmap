package user

import (
	"context"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	userentity "roadmap/internal/domain/entities/user"
	"roadmap/internal/infrastructure/database"
)

type UserRepositoryIntegrationTestSuite struct {
	suite.Suite
	repo *userRepository
	db   *database.Database
	ctx  context.Context
}

func (s *UserRepositoryIntegrationTestSuite) SetupSuite() {
	dsn := os.Getenv("TEST_DB_DSN")
	if dsn == "" {
		s.T().Skip("Skipping integration tests: TEST_DB_DSN not set")
		return
	}

	cfg := &database.Config{
		Host:     "localhost",
		Port:     "5432",
		User:     "postgres",
		Password: "postgres",
		DBName:   "roadmap_test",
		SSLMode:  "disable",
	}

	var err error
	s.db, err = database.NewDatabase(cfg)
	require.NoError(s.T(), err, "Failed to connect to test database")

	s.repo = NewUserRepository(s.db).(*userRepository)
	s.ctx = context.Background()
}

func (s *UserRepositoryIntegrationTestSuite) TearDownSuite() {
	if s.db != nil {
		s.db.Close()
	}
}

func (s *UserRepositoryIntegrationTestSuite) SetupTest() {
	s.cleanupTestData()
}

func (s *UserRepositoryIntegrationTestSuite) TearDownTest() {
	s.cleanupTestData()
}

func (s *UserRepositoryIntegrationTestSuite) cleanupTestData() {
	if s.db == nil {
		return
	}

	_, err := s.db.Pool.Exec(s.ctx, "DELETE FROM users WHERE email LIKE 'test%@example.com' OR username LIKE 'testuser%'")
	if err != nil {
		s.T().Logf("Warning: Failed to cleanup test data: %v", err)
	}
}

func (s *UserRepositoryIntegrationTestSuite) TestUserRepository_Create_Success() {
	if s.db == nil {
		s.T().Skip("Database not available")
	}

	user := &userentity.User{
		ID:           uuid.New(),
		Username:     "testuser_create",
		Email:        "test_create@example.com",
		PasswordHash: "$2a$10$testhash",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	createdUser, err := s.repo.Create(s.ctx, user)

	require.NoError(s.T(), err)
	assert.NotNil(s.T(), createdUser)
	assert.Equal(s.T(), user.ID, createdUser.ID)
	assert.Equal(s.T(), user.Username, createdUser.Username)
	assert.Equal(s.T(), user.Email, createdUser.Email)
	assert.Equal(s.T(), user.PasswordHash, createdUser.PasswordHash)
	assert.False(s.T(), createdUser.CreatedAt.IsZero())
	assert.False(s.T(), createdUser.UpdatedAt.IsZero())
}

func (s *UserRepositoryIntegrationTestSuite) TestUserRepository_Create_DuplicateEmail() {
	if s.db == nil {
		s.T().Skip("Database not available")
	}

	user1 := &userentity.User{
		ID:           uuid.New(),
		Username:     "testuser1",
		Email:        "duplicate@example.com",
		PasswordHash: "$2a$10$testhash1",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	user2 := &userentity.User{
		ID:           uuid.New(),
		Username:     "testuser2",
<<<<<<< HEAD
		Email:        "duplicate@example.com",
=======
		Email:        "duplicate@example.com", 
>>>>>>> 1962dcb (feat: register + login + tests)
		PasswordHash: "$2a$10$testhash2",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	_, err1 := s.repo.Create(s.ctx, user1)
	require.NoError(s.T(), err1)

	_, err2 := s.repo.Create(s.ctx, user2)
	assert.Error(s.T(), err2)
	assert.Contains(s.T(), err2.Error(), "duplicate key")
}

func (s *UserRepositoryIntegrationTestSuite) TestUserRepository_GetByID_Success() {
	if s.db == nil {
		s.T().Skip("Database not available")
	}

	user := &userentity.User{
		ID:           uuid.New(),
		Username:     "testuser_get",
		Email:        "test_get@example.com",
		PasswordHash: "$2a$10$testhash",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	createdUser, err := s.repo.Create(s.ctx, user)
	require.NoError(s.T(), err)

	retrievedUser, err := s.repo.GetByID(s.ctx, createdUser.ID)

	require.NoError(s.T(), err)
	assert.Equal(s.T(), createdUser.ID, retrievedUser.ID)
	assert.Equal(s.T(), createdUser.Username, retrievedUser.Username)
	assert.Equal(s.T(), createdUser.Email, retrievedUser.Email)
	assert.Equal(s.T(), createdUser.PasswordHash, retrievedUser.PasswordHash)
}

func (s *UserRepositoryIntegrationTestSuite) TestUserRepository_GetByID_NotFound() {
	if s.db == nil {
		s.T().Skip("Database not available")
	}

	nonExistentID := uuid.New()
	_, err := s.repo.GetByID(s.ctx, nonExistentID)

	assert.Error(s.T(), err)
	assert.Contains(s.T(), err.Error(), "user not found")
}

func (s *UserRepositoryIntegrationTestSuite) TestUserRepository_GetByEmail_Success() {
	if s.db == nil {
		s.T().Skip("Database not available")
	}

	user := &userentity.User{
		ID:           uuid.New(),
		Username:     "testuser_email",
		Email:        "test_email@example.com",
		PasswordHash: "$2a$10$testhash",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	createdUser, err := s.repo.Create(s.ctx, user)
	require.NoError(s.T(), err)

	retrievedUser, err := s.repo.GetByEmail(s.ctx, createdUser.Email)

	require.NoError(s.T(), err)
	assert.Equal(s.T(), createdUser.ID, retrievedUser.ID)
	assert.Equal(s.T(), createdUser.Username, retrievedUser.Username)
	assert.Equal(s.T(), createdUser.Email, retrievedUser.Email)
}

func (s *UserRepositoryIntegrationTestSuite) TestUserRepository_GetByEmail_NotFound() {
	if s.db == nil {
		s.T().Skip("Database not available")
	}

	_, err := s.repo.GetByEmail(s.ctx, "nonexistent@example.com")

	assert.Error(s.T(), err)
	assert.Contains(s.T(), err.Error(), "user not found")
}

func (s *UserRepositoryIntegrationTestSuite) TestUserRepository_EmailExists_True() {
	if s.db == nil {
		s.T().Skip("Database not available")
	}

	user := &userentity.User{
		ID:           uuid.New(),
		Username:     "testuser_exists",
		Email:        "test_exists@example.com",
		PasswordHash: "$2a$10$testhash",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	_, err := s.repo.Create(s.ctx, user)
	require.NoError(s.T(), err)

	exists, err := s.repo.EmailExists(s.ctx, user.Email)

	require.NoError(s.T(), err)
	assert.True(s.T(), exists)
}

func (s *UserRepositoryIntegrationTestSuite) TestUserRepository_EmailExists_False() {
	if s.db == nil {
		s.T().Skip("Database not available")
	}

	exists, err := s.repo.EmailExists(s.ctx, "nonexistent@example.com")

	require.NoError(s.T(), err)
	assert.False(s.T(), exists)
}

func (s *UserRepositoryIntegrationTestSuite) TestUserRepository_UsernameExists_True() {
	if s.db == nil {
		s.T().Skip("Database not available")
	}

	user := &userentity.User{
		ID:           uuid.New(),
		Username:     "testuser_username",
		Email:        "test_username@example.com",
		PasswordHash: "$2a$10$testhash",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	_, err := s.repo.Create(s.ctx, user)
	require.NoError(s.T(), err)

	exists, err := s.repo.UsernameExists(s.ctx, user.Username)

	require.NoError(s.T(), err)
	assert.True(s.T(), exists)
}

func (s *UserRepositoryIntegrationTestSuite) TestUserRepository_UsernameExists_False() {
	if s.db == nil {
		s.T().Skip("Database not available")
	}

	exists, err := s.repo.UsernameExists(s.ctx, "nonexistentuser")

	require.NoError(s.T(), err)
	assert.False(s.T(), exists)
}

func (s *UserRepositoryIntegrationTestSuite) TestUserRepository_ConcurrentAccess() {
	if s.db == nil {
		s.T().Skip("Database not available")
	}

	userCount := 10
	errors := make(chan error, userCount)

	for i := 0; i < userCount; i++ {
		go func(index int) {
			user := &userentity.User{
				ID:           uuid.New(),
				Username:     "testuser_concurrent_" + string(rune(index)),
				Email:        "test_concurrent_" + string(rune(index)) + "@example.com",
				PasswordHash: "$2a$10$testhash",
				CreatedAt:    time.Now(),
				UpdatedAt:    time.Now(),
			}
			_, err := s.repo.Create(s.ctx, user)
			errors <- err
		}(i)
	}

	var errorCount int
	for i := 0; i < userCount; i++ {
		if err := <-errors; err != nil {
			errorCount++
		}
	}

	assert.Equal(s.T(), 0, errorCount)
}

func (s *UserRepositoryIntegrationTestSuite) TestUserRepository_Create_ErrorHandling() {
	if s.db == nil {
		s.T().Skip("Database not available")
	}

	testCases := []struct {
		name        string
		setupUser   func() *userentity.User
		expectError bool
		errorCheck  func(*testing.T, error)
	}{
		{
			name: "duplicate email returns error",
			setupUser: func() *userentity.User {
				user1 := &userentity.User{
					ID:           uuid.New(),
					Username:     "uniqueuser1",
					Email:        "duplicate_email@example.com",
					PasswordHash: "$2a$10$testhash1",
					CreatedAt:    time.Now(),
					UpdatedAt:    time.Now(),
				}
				_, err := s.repo.Create(s.ctx, user1)
				require.NoError(s.T(), err)

				return &userentity.User{
					ID:           uuid.New(),
					Username:     "uniqueuser2",
<<<<<<< HEAD
					Email:        "duplicate_email@example.com",
=======
					Email:        "duplicate_email@example.com", 
>>>>>>> 1962dcb (feat: register + login + tests)
					PasswordHash: "$2a$10$testhash2",
					CreatedAt:    time.Now(),
					UpdatedAt:    time.Now(),
				}
			},
			expectError: true,
			errorCheck: func(t *testing.T, err error) {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "duplicate")
			},
		},
		{
			name: "duplicate username returns error",
			setupUser: func() *userentity.User {
				user1 := &userentity.User{
					ID:           uuid.New(),
					Username:     "duplicate_username",
					Email:        "unique1@example.com",
					PasswordHash: "$2a$10$testhash1",
					CreatedAt:    time.Now(),
					UpdatedAt:    time.Now(),
				}
				_, err := s.repo.Create(s.ctx, user1)
				require.NoError(s.T(), err)

				return &userentity.User{
					ID:           uuid.New(),
<<<<<<< HEAD
					Username:     "duplicate_username",
=======
					Username:     "duplicate_username", 
>>>>>>> 1962dcb (feat: register + login + tests)
					Email:        "unique2@example.com",
					PasswordHash: "$2a$10$testhash2",
					CreatedAt:    time.Now(),
					UpdatedAt:    time.Now(),
				}
			},
			expectError: true,
			errorCheck: func(t *testing.T, err error) {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "duplicate")
			},
		},
		{
			name: "nil user returns error",
			setupUser: func() *userentity.User {
				return nil
			},
			expectError: true,
			errorCheck: func(t *testing.T, err error) {
				assert.Error(t, err)
			},
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			user := tc.setupUser()
			if user == nil {
				_, err := s.repo.Create(s.ctx, nil)
				if tc.expectError {
					assert.Error(s.T(), err)
					if tc.errorCheck != nil {
						tc.errorCheck(s.T(), err)
					}
				}
				return
			}

			_, err := s.repo.Create(s.ctx, user)

			if tc.expectError {
				assert.Error(s.T(), err)
				if tc.errorCheck != nil {
					tc.errorCheck(s.T(), err)
				}
			} else {
				assert.NoError(s.T(), err)
			}
		})
	}
}

func (s *UserRepositoryIntegrationTestSuite) TestUserRepository_GetByID_ErrorHandling() {
	if s.db == nil {
		s.T().Skip("Database not available")
	}

	testCases := []struct {
		name        string
		setupID     func() uuid.UUID
		expectError bool
		errorCheck  func(*testing.T, error)
	}{
		{
			name: "non-existent ID returns not found error",
			setupID: func() uuid.UUID {
<<<<<<< HEAD
				return uuid.New()
=======
				return uuid.New() 
>>>>>>> 1962dcb (feat: register + login + tests)
			},
			expectError: true,
			errorCheck: func(t *testing.T, err error) {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "user not found")
			},
		},
		{
			name: "nil UUID returns error",
			setupID: func() uuid.UUID {
				return uuid.Nil
			},
			expectError: true,
			errorCheck: func(t *testing.T, err error) {
				assert.Error(t, err)
			},
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			id := tc.setupID()
			_, err := s.repo.GetByID(s.ctx, id)

			if tc.expectError {
				assert.Error(s.T(), err)
				if tc.errorCheck != nil {
					tc.errorCheck(s.T(), err)
				}
			} else {
				assert.NoError(s.T(), err)
			}
		})
	}
}

func (s *UserRepositoryIntegrationTestSuite) TestUserRepository_GetByEmail_ErrorHandling() {
	if s.db == nil {
		s.T().Skip("Database not available")
	}

	testCases := []struct {
		name        string
		email       string
		expectError bool
		errorCheck  func(*testing.T, error)
	}{
		{
			name:        "non-existent email returns not found error",
			email:       "nonexistent@example.com",
			expectError: true,
			errorCheck: func(t *testing.T, err error) {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "user not found")
			},
		},
		{
			name:        "empty email returns error",
			email:       "",
			expectError: true,
			errorCheck: func(t *testing.T, err error) {
				assert.Error(t, err)
			},
		},
		{
			name:        "invalid email format still queries",
			email:       "not-an-email",
			expectError: true,
			errorCheck: func(t *testing.T, err error) {
				assert.Error(t, err)
			},
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			_, err := s.repo.GetByEmail(s.ctx, tc.email)

			if tc.expectError {
				assert.Error(s.T(), err)
				if tc.errorCheck != nil {
					tc.errorCheck(s.T(), err)
				}
			} else {
				assert.NoError(s.T(), err)
			}
		})
	}
}

func (s *UserRepositoryIntegrationTestSuite) TestUserRepository_EmailExists_EdgeCases() {
	if s.db == nil {
		s.T().Skip("Database not available")
	}

	testCases := []struct {
		name      string
		email     string
		setupUser bool
		expected  bool
	}{
		{
			name:      "empty email returns false",
			email:     "",
			setupUser: false,
			expected:  false,
		},
		{
			name:      "case sensitive email check",
			email:     "Test@Example.com",
			setupUser: true,
<<<<<<< HEAD
			expected:  false,
=======
			expected:  false, 
>>>>>>> 1962dcb (feat: register + login + tests)
		},
		{
			name:      "email with special characters",
			email:     "test+tag@example.com",
			setupUser: true,
			expected:  true,
		},
		{
			name:      "very long email",
			email:     "a" + strings.Repeat("b", 200) + "@example.com",
			setupUser: false,
			expected:  false,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			if tc.setupUser {
				user := &userentity.User{
					ID:           uuid.New(),
					Username:     "testuser_" + tc.name,
					Email:        tc.email,
					PasswordHash: "$2a$10$testhash",
					CreatedAt:    time.Now(),
					UpdatedAt:    time.Now(),
				}
				_, err := s.repo.Create(s.ctx, user)
				require.NoError(s.T(), err)
			}

			exists, err := s.repo.EmailExists(s.ctx, tc.email)
			require.NoError(s.T(), err)
			assert.Equal(s.T(), tc.expected, exists)
		})
	}
}

func (s *UserRepositoryIntegrationTestSuite) TestUserRepository_UsernameExists_EdgeCases() {
	if s.db == nil {
		s.T().Skip("Database not available")
	}

	testCases := []struct {
		name      string
		username  string
		setupUser bool
		expected  bool
	}{
		{
			name:      "empty username returns false",
			username:  "",
			setupUser: false,
			expected:  false,
		},
		{
			name:      "case sensitive username check",
			username:  "TestUser",
			setupUser: true,
<<<<<<< HEAD
			expected:  false,
=======
			expected:  false, 
>>>>>>> 1962dcb (feat: register + login + tests)
		},
		{
			name:      "username with numbers",
			username:  "user123",
			setupUser: true,
			expected:  true,
		},
		{
			name:      "username with underscores",
			username:  "user_name",
			setupUser: true,
			expected:  true,
		},
		{
			name:      "very long username",
			username:  strings.Repeat("a", 200),
			setupUser: false,
			expected:  false,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			if tc.setupUser {
				user := &userentity.User{
					ID:           uuid.New(),
					Username:     tc.username,
					Email:        "test_" + tc.name + "@example.com",
					PasswordHash: "$2a$10$testhash",
					CreatedAt:    time.Now(),
					UpdatedAt:    time.Now(),
				}
				_, err := s.repo.Create(s.ctx, user)
				require.NoError(s.T(), err)
			}

			exists, err := s.repo.UsernameExists(s.ctx, tc.username)
			require.NoError(s.T(), err)
			assert.Equal(s.T(), tc.expected, exists)
		})
	}
}

func (s *UserRepositoryIntegrationTestSuite) TestUserRepository_DataIntegrity() {
	if s.db == nil {
		s.T().Skip("Database not available")
	}

	user := &userentity.User{
		ID:           uuid.New(),
		Username:     "integrity_test",
		Email:        "integrity@example.com",
		PasswordHash: "$2a$10$testhash",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	createdUser, err := s.repo.Create(s.ctx, user)
	require.NoError(s.T(), err)

	assert.Equal(s.T(), user.ID, createdUser.ID)
	assert.Equal(s.T(), user.Username, createdUser.Username)
	assert.Equal(s.T(), user.Email, createdUser.Email)
	assert.Equal(s.T(), user.PasswordHash, createdUser.PasswordHash)
	assert.WithinDuration(s.T(), user.CreatedAt, createdUser.CreatedAt, time.Second)
	assert.WithinDuration(s.T(), user.UpdatedAt, createdUser.UpdatedAt, time.Second)

	retrievedByID, err := s.repo.GetByID(s.ctx, createdUser.ID)
	require.NoError(s.T(), err)
	assert.Equal(s.T(), createdUser.ID, retrievedByID.ID)
	assert.Equal(s.T(), createdUser.Username, retrievedByID.Username)
	assert.Equal(s.T(), createdUser.Email, retrievedByID.Email)

	retrievedByEmail, err := s.repo.GetByEmail(s.ctx, createdUser.Email)
	require.NoError(s.T(), err)
	assert.Equal(s.T(), createdUser.ID, retrievedByEmail.ID)
	assert.Equal(s.T(), createdUser.Username, retrievedByEmail.Username)
	assert.Equal(s.T(), createdUser.Email, retrievedByEmail.Email)

	assert.Equal(s.T(), retrievedByID.ID, retrievedByEmail.ID)
	assert.Equal(s.T(), retrievedByID.Username, retrievedByEmail.Username)
	assert.Equal(s.T(), retrievedByID.Email, retrievedByEmail.Email)
}

func (s *UserRepositoryIntegrationTestSuite) TestUserRepository_TransactionIsolation() {
	if s.db == nil {
		s.T().Skip("Database not available")
	}

	user1 := &userentity.User{
		ID:           uuid.New(),
		Username:     "isolation_user1",
		Email:        "isolation1@example.com",
		PasswordHash: "$2a$10$testhash1",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	user2 := &userentity.User{
		ID:           uuid.New(),
		Username:     "isolation_user2",
		Email:        "isolation2@example.com",
		PasswordHash: "$2a$10$testhash2",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	created1, err1 := s.repo.Create(s.ctx, user1)
	require.NoError(s.T(), err1)

	created2, err2 := s.repo.Create(s.ctx, user2)
	require.NoError(s.T(), err2)

	assert.NotEqual(s.T(), created1.ID, created2.ID)
	assert.NotEqual(s.T(), created1.Username, created2.Username)
	assert.NotEqual(s.T(), created1.Email, created2.Email)

	retrieved1, err := s.repo.GetByID(s.ctx, created1.ID)
	require.NoError(s.T(), err)
	assert.Equal(s.T(), created1.ID, retrieved1.ID)

	retrieved2, err := s.repo.GetByID(s.ctx, created2.ID)
	require.NoError(s.T(), err)
	assert.Equal(s.T(), created2.ID, retrieved2.ID)
}

func TestUserRepositoryIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(UserRepositoryIntegrationTestSuite))
}
