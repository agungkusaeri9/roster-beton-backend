package handler

import (
	"database/sql"

	"go-arch/internal/usecase"

	"github.com/gin-gonic/gin"
)

type ConfigHandler struct {
	uc usecase.ConfigUsecase
}

func NewConfigHandler(r *gin.RouterGroup, uc usecase.ConfigUsecase) {
	h := &ConfigHandler{uc: uc}

	r.GET("/configs", h.GetAll)
	r.GET("/configs/:key", h.GetByKey)
}

func (h *ConfigHandler) GetAll(c *gin.Context) {
	format := c.Query("format")
	if format == "list" {
		configs, err := h.uc.GetAll()
		if err != nil {
			ErrorResponse(c, 500, err.Error())
			return
		}
		SuccessResponse(c, 200, "Configs retrieved successfully", configs)
		return
	}

	configMap, err := h.uc.GetMap()
	if err != nil {
		ErrorResponse(c, 500, err.Error())
		return
	}

	SuccessResponse(c, 200, "Configs retrieved successfully", configMap)
}

func (h *ConfigHandler) GetByKey(c *gin.Context) {
	key := c.Param("key")
	if key == "" {
		ErrorResponse(c, 400, "Config key is required")
		return
	}

	config, err := h.uc.GetByKey(key)
	if err != nil {
		if err == sql.ErrNoRows {
			ErrorResponse(c, 404, "Config not found")
			return
		}
		ErrorResponse(c, 500, err.Error())
		return
	}

	SuccessResponse(c, 200, "Config detail retrieved successfully", config)
}
