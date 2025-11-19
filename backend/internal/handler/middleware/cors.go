package middleware

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func CORSMiddleware() gin.HandlerFunc {
	config := cors.DefaultConfig()

	config.AllowAllOrigins = true

	config.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"}

	config.AllowHeaders = []string{
		"Origin",
		"Content-Length",
		"Content-Type",
		"Authorization",
		"Accept",
		"X-Requested-With",
	}

	config.ExposeHeaders = []string{
		"Content-Length",
		"Content-Type",
		"Authorization",
	}

	config.AllowCredentials = true

	config.MaxAge = 43200

	return cors.New(config)
}
