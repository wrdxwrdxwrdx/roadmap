package userhandler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	userdto "roadmap/internal/domain/dto/user"
	userusecase "roadmap/internal/usecase/user"
)

type UserHandler struct {
	createUserUseCase *userusecase.CreateUserUseCase
}

func NewUserHandler(createUserUseCase *userusecase.CreateUserUseCase) *UserHandler {
	return &UserHandler{
		createUserUseCase: createUserUseCase,
	}
}

func (h *UserHandler) CreateUser(c *gin.Context) {
	var req userdto.CreateUserRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request data",
			"details": err.Error(),
		})
		return
	}

	response, err := h.createUserUseCase.Execute(c.Request.Context(), req)
	if err != nil {
		statusCode := http.StatusInternalServerError
		errorMessage := "Failed to create user"

		if errors.Is(err, userusecase.ErrEmailAlreadyExists) {
			statusCode = http.StatusConflict
			errorMessage = "Email already exists"
		} else if errors.Is(err, userusecase.ErrUsernameAlreadyExists) {
			statusCode = http.StatusConflict
			errorMessage = "Username already exists"
		} else {
			var passwordErr *userusecase.PasswordValidationError
			if errors.As(err, &passwordErr) {
				statusCode = http.StatusBadRequest
				errorMessage = passwordErr.Error()
			}
		}

		c.JSON(statusCode, gin.H{
			"error": errorMessage,
		})
		return
	}

	c.JSON(http.StatusCreated, response)
}
