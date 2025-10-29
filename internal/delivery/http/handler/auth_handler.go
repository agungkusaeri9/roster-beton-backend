package handler

import (
	"go-arch/internal/usecase"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	uc usecase.AuthUsecase
}

func NewAuthHandler(r *gin.Engine, uc usecase.AuthUsecase) {
	h := &AuthHandler{uc: uc}
	r.POST("/register", h.Register)
	r.POST("/login", h.Login)
}

func (h *AuthHandler) Register(c *gin.Context) {
	var params RegisterRequest

	if err := c.ShouldBindJSON(&params); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	userID, err := h.uc.Register(params.Name, params.Username, params.Password)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	data := gin.H{"user_id": userID}
	SuccessResponse(c, 201, "User registered successfully", data)
}

func (h *AuthHandler) Login(c *gin.Context) {
	var params LoginRequest

	err := c.ShouldBindBodyWithJSON(&params)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	token, err := h.uc.Login(params.Username, params.Password)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	data := gin.H{"token": token}
	SuccessResponse(c, 200, "User registered successfully", data)
	// if err != nil {
	// 	c.JSON(400, gin.H{"error": err.Error()})
	// 	return
	// }
}

type RegisterRequest struct {
	Name     string `json:"name"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
