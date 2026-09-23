package handler

import (
	"database/sql"
	"strconv"

	"go-arch/internal/usecase"

	"github.com/gin-gonic/gin"
)

type ProductHandler struct {
	uc usecase.ProductUsecase
}

func NewProductHandler(r *gin.RouterGroup, uc usecase.ProductUsecase) {
	h := &ProductHandler{uc: uc}

	r.GET("/products", h.GetAll)
	r.GET("/products/:slug", h.GetBySlug)
}

func (h *ProductHandler) GetAll(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "10")
	category := c.Query("category")
	search := c.Query("search")
	if search == "" {
		search = c.Query("q")
	}

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		limit = 10
	}

	products, pagination, err := h.uc.GetAll(page, limit, category, search)
	if err != nil {
		ErrorResponse(c, 500, err.Error())
		return
	}

	SuccessPaginatedResponse(c, 200, "Products retrieved successfully", products, pagination)
}

func (h *ProductHandler) GetBySlug(c *gin.Context) {
	slug := c.Param("slug")
	if slug == "" {
		ErrorResponse(c, 400, "Product slug is required")
		return
	}

	product, err := h.uc.GetBySlug(slug)
	if err != nil {
		if err == sql.ErrNoRows {
			ErrorResponse(c, 404, "Product not found")
			return
		}
		ErrorResponse(c, 500, err.Error())
		return
	}

	SuccessResponse(c, 200, "Product detail retrieved successfully", product)
}
