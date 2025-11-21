package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	jwtservice "roadmap/internal/pkg/jwt"
)

const (
	UserIDKey   = "user_id"
	UsernameKey = "username"
	EmailKey    = "email"
)

func AuthMiddleware(jwtService *jwtservice.JWTService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Authorization header is required",
			})
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid authorization header format. Expected: Bearer <token>",
			})
			c.Abort()
			return
		}

		token := parts[1]

		claims, err := jwtService.ValidateToken(token)
		if err != nil {
			statusCode := http.StatusUnauthorized
			errorMessage := "Invalid or expired token"

			switch err {
			case jwtservice.ErrExpiredToken:
				errorMessage = "Token has expired"
			case jwtservice.ErrInvalidToken:
				errorMessage = "Invalid token"
			}

			c.JSON(statusCode, gin.H{
				"error": errorMessage,
			})
			c.Abort()
			return
		}

		c.Set(UserIDKey, claims.UserID)
		c.Set(UsernameKey, claims.Username)
		c.Set(EmailKey, claims.Email)

		c.Next()
	}
}

func GetUserID(c *gin.Context) (string, bool) {
	userID, exists := c.Get(UserIDKey)
	if !exists {
		return "", false
	}
	return userID.(string), true
}

func GetUsername(c *gin.Context) (string, bool) {
	username, exists := c.Get(UsernameKey)
	if !exists {
		return "", false
	}
	return username.(string), true
}

func GetEmail(c *gin.Context) (string, bool) {
	email, exists := c.Get(EmailKey)
	if !exists {
		return "", false
	}
	return email.(string), true
}
