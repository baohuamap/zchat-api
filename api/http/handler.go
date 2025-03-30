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
	UploadImage(ctx *gin.Context)
	GetProfileImage(ctx *gin.Context)
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

func (h *handler) UploadImage(c *gin.Context) {
    file, err := c.FormFile("file")
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file"})
        return
    }

    openedFile, err := file.Open()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open file"})
        return
    }
    defer openedFile.Close()

    // Gọi service để upload ảnh
    message, err := h.userService.UploadImage(c.Request.Context(), openedFile, file.Filename)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": message})
}

func (h *handler) GetProfileImage(c *gin.Context) {
    userID := c.Param("userID")

    // Gọi service để lấy dữ liệu ảnh
    imageData, err := h.userService.GetProfileImage(c.Request.Context(), userID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    // Trả về dữ liệu ảnh dưới dạng binary
    c.Data(http.StatusOK, "image/jpeg", imageData)
}
