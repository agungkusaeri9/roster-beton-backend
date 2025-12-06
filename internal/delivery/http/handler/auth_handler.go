package handler

import (
	"fmt"
	"strings"
	"time"

	"go-arch/internal/usecase"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type AuthHandler struct {
	uc usecase.AuthUsecase
}

func (h *AuthHandler) Health(c *gin.Context) {
	SuccessResponse(c, 200, "Service is healthy", gin.H{
		"status": "ok",
		"service": "go-arch",
	})
}

func NewAuthHandler(r *gin.Engine, uc usecase.AuthUsecase) {
	h := &AuthHandler{uc: uc}
	
	// Health check endpoint
	r.GET("/health", h.Health)
	
	// Auth endpoints
	r.POST("/register", h.Register)
	r.POST("/login", h.Login)
}

func (h *AuthHandler) Register(c *gin.Context) {
	var params RegisterRequest

	// Bind JSON request
	if err := c.ShouldBindJSON(&params); err != nil {
		// Handle validation errors from gin validator
		if validationErr, ok := err.(validator.ValidationErrors); ok {
			ErrorResponse(c, 400, formatValidationErrors(validationErr))
			return
		}
		ErrorResponse(c, 400, err.Error())
		return
	}

	// Validation will be handled in usecase layer
	userID, err := h.uc.Register(params.Name, params.Username, params.Password)
	if err != nil {
		ErrorResponse(c, 400, err.Error())
		return
	}

	data := gin.H{"user_id": userID}
	SuccessResponse(c, 201, "User registered successfully", data)
}

func (h *AuthHandler) Login(c *gin.Context) {
	var params LoginRequest

	// Bind JSON request
	err := c.ShouldBindBodyWithJSON(&params)
	if err != nil {
		// Handle validation errors from gin validator
		if validationErr, ok := err.(validator.ValidationErrors); ok {
			ErrorResponse(c, 400, formatValidationErrors(validationErr))
			return
		}
		ErrorResponse(c, 400, err.Error())
		return
	}

	// Validation will be handled in usecase layer
	result, err := h.uc.Login(params.Username, params.Password)
	if err != nil {
		ErrorResponse(c, 400, err.Error())
		return
	}

	data := gin.H{
		"token": gin.H{
			"access_token": result.Token,
			"expired_at":  result.ExpiredAt.Format(time.RFC3339),
		},
		"user": gin.H{
			"username": result.User.Username,
			"name":     result.User.Name,
		},
	}
	SuccessResponse(c, 200, "Login successful", data)
}

type RegisterRequest struct {
	Name     string `json:"name" binding:"required,min=2,max=100"`
	Username string `json:"username" binding:"required,min=3,max=150"`
	Password string `json:"password" binding:"required,min=6,max=255"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// formatValidationErrors formats gin validator errors into readable string
func formatValidationErrors(err validator.ValidationErrors) string {
	var messages []string
	for _, e := range err {
		field := e.Field()
		tag := e.Tag()
		
		var message string
		switch tag {
		case "required":
			message = fmt.Sprintf("%s is required", field)
		case "min":
			message = fmt.Sprintf("%s must be at least %s characters", field, e.Param())
		case "max":
			message = fmt.Sprintf("%s must not exceed %s characters", field, e.Param())
		case "email":
			message = fmt.Sprintf("%s must be a valid email", field)
		default:
			message = fmt.Sprintf("%s is invalid", field)
		}
		messages = append(messages, message)
	}
	return strings.Join(messages, "; ")
}
