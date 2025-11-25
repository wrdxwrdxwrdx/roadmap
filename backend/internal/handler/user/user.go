package userhandler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	userdto "roadmap/internal/domain/dto/user"
	"roadmap/internal/handler/middleware"
	userusecase "roadmap/internal/usecase/user"
)

type UserHandler struct {
	createUserUseCase *userusecase.CreateUserUseCase
	registerUseCase   *userusecase.RegisterUseCase
	loginUseCase      *userusecase.LoginUseCase
}

func NewUserHandler(
	createUserUseCase *userusecase.CreateUserUseCase,
	registerUseCase *userusecase.RegisterUseCase,
	loginUseCase *userusecase.LoginUseCase,
) *UserHandler {
	return &UserHandler{
		createUserUseCase: createUserUseCase,
		registerUseCase:   registerUseCase,
		loginUseCase:      loginUseCase,
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

func (h *UserHandler) Register(c *gin.Context) {
	var req userdto.RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request data",
			"details": err.Error(),
		})
		return
	}

	response, err := h.registerUseCase.Execute(c.Request.Context(), req)
	if err != nil {
		statusCode := http.StatusInternalServerError
		errorMessage := "Failed to register user"

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

func (h *UserHandler) Login(c *gin.Context) {
	var req userdto.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request data",
			"details": err.Error(),
		})
		return
	}

	response, err := h.loginUseCase.Execute(c.Request.Context(), req)
	if err != nil {
		var statusCode int
		var errorMessage string

		if errors.Is(err, userusecase.ErrInvalidCredentials) {
			statusCode = http.StatusUnauthorized
			errorMessage = "Invalid email or password"
		} else {
			statusCode = http.StatusInternalServerError
			errorMessage = "Failed to login"
		}

		c.JSON(statusCode, gin.H{
			"error": errorMessage,
		})
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *UserHandler) GetProfile(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User ID not found in context",
		})
		return
	}

	username, _ := middleware.GetUsername(c)
	email, _ := middleware.GetEmail(c)

	c.JSON(http.StatusOK, gin.H{
		"user_id":  userID,
		"username": username,
		"email":    email,
	})
}
