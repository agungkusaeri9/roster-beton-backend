package handler

import (
	"go-arch/internal/entity"

	"github.com/gin-gonic/gin"
)

type APIResponse struct {
	Status  bool        `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type PaginatedAPIResponse struct {
	Status     bool                  `json:"status"`
	Message    string                `json:"message"`
	Data       interface{}           `json:"data"`
	Pagination entity.PaginationMeta `json:"pagination"`
}

// helper buat success
func SuccessResponse(c *gin.Context, statusCode int, message string, data interface{}) {
	c.JSON(statusCode, APIResponse{
		Status:  true,
		Message: message,
		Data:    data,
	})
}

// helper buat success dengan pagination
func SuccessPaginatedResponse(c *gin.Context, statusCode int, message string, data interface{}, pagination entity.PaginationMeta) {
	c.JSON(statusCode, PaginatedAPIResponse{
		Status:     true,
		Message:    message,
		Data:       data,
		Pagination: pagination,
	})
}

// helper buat error
func ErrorResponse(c *gin.Context, statusCode int, message string) {
	c.JSON(statusCode, APIResponse{
		Status:  false,
		Message: message,
		Data:    nil,
	})
}
