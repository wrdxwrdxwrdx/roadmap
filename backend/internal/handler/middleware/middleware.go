package middleware

import (
	"github.com/gin-gonic/gin"
)

func SetupMiddleware(router *gin.Engine) {
	router.Use(RecoveryMiddleware())

	router.Use(CORSMiddleware())

	router.Use(LoggingMiddleware())
}
