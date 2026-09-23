package handler

import (
	"database/sql"
	"strconv"

	"go-arch/internal/usecase"

	"github.com/gin-gonic/gin"
)

type GalleryHandler struct {
	uc usecase.GalleryUsecase
}

func NewGalleryHandler(r *gin.RouterGroup, uc usecase.GalleryUsecase) {
	h := &GalleryHandler{uc: uc}

	r.GET("/galleries", h.GetAll)
	r.GET("/galleries/:id", h.GetByID)
}

func (h *GalleryHandler) GetAll(c *gin.Context) {
	category := c.Query("category")
	galleries, err := h.uc.GetAll(category)
	if err != nil {
		ErrorResponse(c, 500, err.Error())
		return
	}

	SuccessResponse(c, 200, "Galleries retrieved successfully", galleries)
}

func (h *GalleryHandler) GetByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		ErrorResponse(c, 400, "Invalid gallery id")
		return
	}

	item, err := h.uc.GetByID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			ErrorResponse(c, 404, "Gallery item not found")
			return
		}
		ErrorResponse(c, 500, err.Error())
		return
	}

	SuccessResponse(c, 200, "Gallery detail retrieved successfully", item)
}
