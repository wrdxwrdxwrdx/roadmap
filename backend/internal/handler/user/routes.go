package userhandler

import (
	"github.com/gin-gonic/gin"
)

func SetupUserRoutes(router *gin.RouterGroup, handler *UserHandler, authMiddleware gin.HandlerFunc) {
	users := router.Group("/users")
	{
		users.POST("create", handler.CreateUser)
		users.POST("register", handler.Register)
		users.POST("login", handler.Login)

		protected := users.Group("")
		protected.Use(authMiddleware)
		{
			protected.GET("profile", handler.GetProfile)
		}
	}
}
