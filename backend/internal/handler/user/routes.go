package userhandler

import (
	"github.com/gin-gonic/gin"
)

func SetupUserRoutes(router *gin.RouterGroup, handler *UserHandler) {
	users := router.Group("/users")
	{
		users.POST("create", handler.CreateUser)
	}
}
