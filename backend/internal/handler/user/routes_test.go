package userhandler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	userusecase "roadmap/internal/usecase/user"
	jwtservice "roadmap/internal/pkg/jwt"
)

func TestSetupUserRoutes(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	
	// Create real use cases with nil repositories (they won't be called in this test)
	createUseCase := userusecase.NewCreateUserUseCase(nil)
	registerUseCase := userusecase.NewRegisterUseCase(nil, jwtservice.NewJWTService("test-secret", 24*3600*1000000000))
	loginUseCase := userusecase.NewLoginUseCase(nil, jwtservice.NewJWTService("test-secret", 24*3600*1000000000))
	
	handler := NewUserHandler(createUseCase, registerUseCase, loginUseCase)
	authMiddleware := func(c *gin.Context) {
		c.Set("user_id", "test-user-id")
		c.Next()
	}

	api := router.Group("/api/v1")
	SetupUserRoutes(api, handler, authMiddleware)

	// Test create route exists (will fail validation but route exists)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/users/create", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.NotEqual(t, http.StatusNotFound, w.Code, "create route should exist")

	// Test register route exists
	req = httptest.NewRequest(http.MethodPost, "/api/v1/users/register", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.NotEqual(t, http.StatusNotFound, w.Code, "register route should exist")

	// Test login route exists
	req = httptest.NewRequest(http.MethodPost, "/api/v1/users/login", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.NotEqual(t, http.StatusNotFound, w.Code, "login route should exist")

	// Test protected profile route with auth middleware
	req = httptest.NewRequest(http.MethodGet, "/api/v1/users/profile", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	// Should be 200 OK because auth middleware sets user_id
	assert.Equal(t, http.StatusOK, w.Code, "profile route should work with auth")
}

