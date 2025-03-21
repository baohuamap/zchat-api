package http

import (
	"net/http"

	"github.com/baohuamap/zchat-api/dto"
	"github.com/baohuamap/zchat-api/service"
	"github.com/gin-gonic/gin"
)

type Handler interface {
	ServerStatus(ctx *gin.Context)
	CreateUser(ctx *gin.Context)
	Login(ctx *gin.Context)
	Logout(ctx *gin.Context)
}

type handler struct {
	userService service.User
}

func NewHandler(u service.User) Handler {
	return &handler{userService: u}
}

func (h *handler) ServerStatus(ctx *gin.Context) {
	response := dto.ServerStatusResponse{
		Status: "healthy",
	}
	ctx.JSON(http.StatusOK, response)
}

func (h *handler) CreateUser(c *gin.Context) {
	var u dto.CreateUserReq
	if err := c.ShouldBindJSON(&u); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := h.userService.CreateUser(c.Request.Context(), &u)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}

func (h *handler) Login(c *gin.Context) {
	var user dto.LoginUserReq
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	u, err := h.userService.Login(c.Request.Context(), &user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.SetCookie("jwt", u.AccessToken, 60*60*24, "/", "localhost", false, true)
	c.JSON(http.StatusOK, u)
}

func (h *handler) Logout(c *gin.Context) {
	c.SetCookie("jwt", "", -1, "", "", false, true)
	c.JSON(http.StatusOK, gin.H{"message": "logout successful"})
}
