package api

import (
	"net/http"

	"github.com/baohuamap/zchat-api/dto"
	"github.com/gin-gonic/gin"

	"github.com/baohuamap/zchat-api/models"
)

type handler struct {
	Service models.Service
}

func (h *handler) ServerStatus(ctx *gin.Context) {
	response := dto.ServerStatusResponse{
		Status: "healthy",
	}
	ctx.JSON(http.StatusOK, response)
}

func NewHandler(s models.Service) *handler {
	return &handler{}
}

func (h *handler) Create(c *gin.Context) {
	var u models.CreateUserReq
	if err := c.ShouldBindJSON(&u); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := h.Service.Create(c.Request.Context(), &u)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}

func (h *handler) Login(c *gin.Context) {
	var user models.LoginUserReq
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	u, err := h.Service.Login(c.Request.Context(), &user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// c.SetCookie("jwt", u.accessToken, 60*60*24, "/", "localhost", false, true)
	c.JSON(http.StatusOK, u)
}

func (h *handler) Logout(c *gin.Context) {
	// c.SetCookie("jwt", "", -1, "", "", false, true)
	c.JSON(http.StatusOK, gin.H{"message": "logout successful"})
}
