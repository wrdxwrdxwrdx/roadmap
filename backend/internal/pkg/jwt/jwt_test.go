package jwt

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

func TestNewJWTService(t *testing.T) {
	secretKey := "test-secret-key"
	expiresIn := 24 * time.Hour

	service := NewJWTService(secretKey, expiresIn)

	assert.NotNil(t, service)
	assert.Equal(t, []byte(secretKey), service.secretKey)
	assert.Equal(t, expiresIn, service.expiresIn)
}

func TestJWTService_GenerateToken(t *testing.T) {
	service := NewJWTService("test-secret-key", 24*time.Hour)

	userID := "550e8400-e29b-41d4-a716-446655440000"
	username := "testuser"
	email := "test@example.com"

	token, err := service.GenerateToken(userID, username, email)

	assert.NoError(t, err)
	assert.NotEmpty(t, token)
}

func TestJWTService_GenerateToken_DifferentUsers(t *testing.T) {
	service := NewJWTService("test-secret-key", 24*time.Hour)

	token1, _ := service.GenerateToken("user1", "user1", "user1@example.com")
	token2, _ := service.GenerateToken("user2", "user2", "user2@example.com")

	assert.NotEqual(t, token1, token2)
}

func TestJWTService_ValidateToken_Success(t *testing.T) {
	service := NewJWTService("test-secret-key", 24*time.Hour)

	userID := "550e8400-e29b-41d4-a716-446655440000"
	username := "testuser"
	email := "test@example.com"

	token, _ := service.GenerateToken(userID, username, email)

	claims, err := service.ValidateToken(token)

	assert.NoError(t, err)
	assert.NotNil(t, claims)
	assert.Equal(t, userID, claims.UserID)
	assert.Equal(t, username, claims.Username)
	assert.Equal(t, email, claims.Email)
}

func TestJWTService_ValidateToken_InvalidToken(t *testing.T) {
	service := NewJWTService("test-secret-key", 24*time.Hour)

	invalidToken := "invalid.token.here"

	claims, err := service.ValidateToken(invalidToken)

	assert.Error(t, err)
	assert.Nil(t, claims)
	assert.Equal(t, ErrInvalidToken, err)
}

func TestJWTService_ValidateToken_WrongSecretKey(t *testing.T) {
	service1 := NewJWTService("secret-key-1", 24*time.Hour)
	service2 := NewJWTService("secret-key-2", 24*time.Hour)

	token, _ := service1.GenerateToken("user1", "user1", "user1@example.com")

	claims, err := service2.ValidateToken(token)

	assert.Error(t, err)
	assert.Nil(t, claims)
	assert.Equal(t, ErrInvalidToken, err)
}

func TestJWTService_ValidateToken_ExpiredToken(t *testing.T) {
	service := NewJWTService("test-secret-key", -time.Hour)

	token, _ := service.GenerateToken("user1", "user1", "user1@example.com")

	time.Sleep(100 * time.Millisecond)

	claims, err := service.ValidateToken(token)

	assert.Error(t, err)
	assert.Nil(t, claims)
	assert.Equal(t, ErrExpiredToken, err)
}

func TestJWTService_ValidateToken_ManipulatedToken(t *testing.T) {
	service := NewJWTService("test-secret-key", 24*time.Hour)

	token, _ := service.GenerateToken("user1", "user1", "user1@example.com")
	manipulatedToken := token[:len(token)-5] + "xxxxx"

	claims, err := service.ValidateToken(manipulatedToken)

	assert.Error(t, err)
	assert.Nil(t, claims)
	assert.Equal(t, ErrInvalidToken, err)
}

func TestJWTService_ValidateToken_EmptyToken(t *testing.T) {
	service := NewJWTService("test-secret-key", 24*time.Hour)

	claims, err := service.ValidateToken("")

	assert.Error(t, err)
	assert.Nil(t, claims)
	assert.Equal(t, ErrInvalidToken, err)
}

func TestJWTService_TokenContainsCorrectClaims(t *testing.T) {
	service := NewJWTService("test-secret-key", 24*time.Hour)

	userID := "550e8400-e29b-41d4-a716-446655440000"
	username := "testuser"
	email := "test@example.com"

	token, _ := service.GenerateToken(userID, username, email)
	claims, _ := service.ValidateToken(token)

	assert.Equal(t, userID, claims.UserID)
	assert.Equal(t, username, claims.Username)
	assert.Equal(t, email, claims.Email)
	assert.NotNil(t, claims.ExpiresAt)
	assert.NotNil(t, claims.IssuedAt)
	assert.NotNil(t, claims.NotBefore)
}

func TestJWTService_TokenExpiration(t *testing.T) {
	service := NewJWTService("test-secret-key", 1*time.Second)

	token, _ := service.GenerateToken("user1", "user1", "user1@example.com")

	claims1, err1 := service.ValidateToken(token)
	assert.NoError(t, err1)
	assert.NotNil(t, claims1)

	time.Sleep(2 * time.Second)

	claims2, err2 := service.ValidateToken(token)
	assert.Error(t, err2)
	assert.Nil(t, claims2)
	assert.Equal(t, ErrExpiredToken, err2)
}

func TestJWTService_ValidateToken_WrongSigningMethod(t *testing.T) {
	service := NewJWTService("test-secret-key", 24*time.Hour)

	wrongMethodToken := jwt.NewWithClaims(jwt.SigningMethodRS256, &Claims{
		UserID:   "user1",
		Username: "user1",
		Email:    "user1@example.com",
	})
	tokenString, _ := wrongMethodToken.SignedString([]byte("wrong-key"))

	claims, err := service.ValidateToken(tokenString)

	assert.Error(t, err)
	assert.Nil(t, claims)
	assert.Equal(t, ErrInvalidToken, err)
}

func TestJWTService_ValidateToken_InvalidClaims(t *testing.T) {
	service := NewJWTService("test-secret-key", 24*time.Hour)

	// Create a token with invalid claims structure
	invalidToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.invalid.claims"

	claims, err := service.ValidateToken(invalidToken)

	assert.Error(t, err)
	assert.Nil(t, claims)
	assert.Equal(t, ErrInvalidToken, err)
}

func TestJWTService_ValidateToken_NonStandardError(t *testing.T) {
	service := NewJWTService("test-secret-key", 24*time.Hour)

	// Token with malformed structure
	malformedToken := "not.a.valid.token"

	claims, err := service.ValidateToken(malformedToken)

	assert.Error(t, err)
	assert.Nil(t, claims)
	assert.Equal(t, ErrInvalidToken, err)
}

