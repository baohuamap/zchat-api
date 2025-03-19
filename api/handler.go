package api

import (
	"net/http"

	"github.com/baohuamap/zchat-api/dto"
	"github.com/gin-gonic/gin"
)

type handler struct {
}

func (h *handler) ServerStatus(ctx *gin.Context) {
	response := dto.ServerStatusResponse{
		Status: "healthy",
	}
	ctx.JSON(http.StatusOK, response)
}
