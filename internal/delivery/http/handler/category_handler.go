package handler

import (
	"database/sql"

	"go-arch/internal/usecase"

	"github.com/gin-gonic/gin"
)

type CategoryHandler struct {
	uc usecase.CategoryUsecase
}

func NewCategoryHandler(r *gin.RouterGroup, uc usecase.CategoryUsecase) {
	h := &CategoryHandler{uc: uc}

	r.GET("/categories", h.GetAll)
	r.GET("/categories/:slug", h.GetBySlug)
}

func (h *CategoryHandler) GetAll(c *gin.Context) {
	categories, err := h.uc.GetAll()
	if err != nil {
		ErrorResponse(c, 500, err.Error())
		return
	}

	SuccessResponse(c, 200, "Categories retrieved successfully", categories)
}

func (h *CategoryHandler) GetBySlug(c *gin.Context) {
	slug := c.Param("slug")
	if slug == "" {
		ErrorResponse(c, 400, "Category slug is required")
		return
	}

	category, err := h.uc.GetBySlug(slug)
	if err != nil {
		if err == sql.ErrNoRows {
			ErrorResponse(c, 404, "Category not found")
			return
		}
		ErrorResponse(c, 500, err.Error())
		return
	}

	SuccessResponse(c, 200, "Category detail retrieved successfully", category)
}
