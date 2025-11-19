package middleware

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

func LoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		c.Next()

		latency := time.Since(start)
		clientIP := c.ClientIP()
		method := c.Request.Method
		statusCode := c.Writer.Status()
		errorMessage := c.Errors.ByType(gin.ErrorTypePrivate).String()

		if raw != "" {
			path = path + "?" + raw
		}

		if errorMessage != "" {
			log.Printf("[%s] %s %s %d %v %s - Error: %s",
				clientIP, method, path, statusCode, latency, c.Request.UserAgent(), errorMessage)
		} else {
			log.Printf("[%s] %s %s %d %v %s",
				clientIP, method, path, statusCode, latency, c.Request.UserAgent())
		}
	}
}
