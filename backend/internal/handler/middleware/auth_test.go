package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	jwtservice "roadmap/internal/pkg/jwt"
)

func TestAuthMiddleware_MissingAuthorizationHeader(t *testing.T) {
	gin.SetMode(gin.TestMode)
	jwtService := jwtservice.NewJWTService("test-secret-key", 24*3600*1000000000)
	middleware := AuthMiddleware(jwtService)

	router := gin.New()
	router.Use(middleware)
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestAuthMiddleware_InvalidHeaderFormat(t *testing.T) {
	gin.SetMode(gin.TestMode)
	jwtService := jwtservice.NewJWTService("test-secret-key", 24*3600*1000000000)
	middleware := AuthMiddleware(jwtService)

	router := gin.New()
	router.Use(middleware)
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	testCases := []struct {
		name   string
		header string
	}{
		{"no bearer prefix", "token123"},
		{"empty token", "Bearer "},
		{"multiple spaces", "Bearer  token123"},
		{"wrong prefix", "Token token123"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/test", nil)
			req.Header.Set("Authorization", tc.header)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusUnauthorized, w.Code)
		})
	}
}

func TestAuthMiddleware_InvalidToken(t *testing.T) {
	gin.SetMode(gin.TestMode)
	jwtService := jwtservice.NewJWTService("test-secret-key", 24*3600*1000000000)
	middleware := AuthMiddleware(jwtService)

	router := gin.New()
	router.Use(middleware)
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", "Bearer invalid.token.here")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestAuthMiddleware_ExpiredToken(t *testing.T) {
	gin.SetMode(gin.TestMode)
	jwtService := jwtservice.NewJWTService("test-secret-key", -3600*1000000000)

	token, _ := jwtService.GenerateToken("user1", "user1", "user1@example.com")

	middleware := AuthMiddleware(jwtService)

	router := gin.New()
	router.Use(middleware)
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestAuthMiddleware_ValidToken(t *testing.T) {
	gin.SetMode(gin.TestMode)
	jwtService := jwtservice.NewJWTService("test-secret-key", 24*3600*1000000000)

	userID := "550e8400-e29b-41d4-a716-446655440000"
	username := "testuser"
	email := "test@example.com"

	token, _ := jwtService.GenerateToken(userID, username, email)

	middleware := AuthMiddleware(jwtService)

	router := gin.New()
	router.Use(middleware)
	router.GET("/test", func(c *gin.Context) {
		ctxUserID, exists := GetUserID(c)
		assert.True(t, exists)
		assert.Equal(t, userID, ctxUserID)

		ctxUsername, exists := GetUsername(c)
		assert.True(t, exists)
		assert.Equal(t, username, ctxUsername)

		ctxEmail, exists := GetEmail(c)
		assert.True(t, exists)
		assert.Equal(t, email, ctxEmail)

		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestGetUserID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(httptest.NewRecorder())

	userID, exists := GetUserID(c)
	assert.False(t, exists)
	assert.Empty(t, userID)

	c.Set(UserIDKey, "test-user-id")
	userID, exists = GetUserID(c)
	assert.True(t, exists)
	assert.Equal(t, "test-user-id", userID)
}

func TestGetUsername(t *testing.T) {
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(httptest.NewRecorder())

	username, exists := GetUsername(c)
	assert.False(t, exists)
	assert.Empty(t, username)

	c.Set(UsernameKey, "testuser")
	username, exists = GetUsername(c)
	assert.True(t, exists)
	assert.Equal(t, "testuser", username)
}

func TestGetEmail(t *testing.T) {
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(httptest.NewRecorder())

	email, exists := GetEmail(c)
	assert.False(t, exists)
	assert.Empty(t, email)

	c.Set(EmailKey, "test@example.com")
	email, exists = GetEmail(c)
	assert.True(t, exists)
	assert.Equal(t, "test@example.com", email)
}

func TestAuthMiddleware_WrongSecretKey(t *testing.T) {
	gin.SetMode(gin.TestMode)
	jwtService1 := jwtservice.NewJWTService("secret-key-1", 24*3600*1000000000)
	jwtService2 := jwtservice.NewJWTService("secret-key-2", 24*3600*1000000000)

	token, _ := jwtService1.GenerateToken("user1", "user1", "user1@example.com")

	middleware := AuthMiddleware(jwtService2)

	router := gin.New()
	router.Use(middleware)
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

