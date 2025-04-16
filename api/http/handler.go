package http

import (
	"net/http"
	"strconv"

	"github.com/baohuamap/zchat-api/dto"
	"github.com/baohuamap/zchat-api/service"
	"github.com/gin-gonic/gin"
)

type Handler interface {
	ServerStatus(ctx *gin.Context)
	CreateUser(ctx *gin.Context)
	Login(ctx *gin.Context)
	Logout(ctx *gin.Context)
	AddFriend(ctx *gin.Context)
	CancelFriendRequest(ctx *gin.Context)
	AcceptFriend(ctx *gin.Context)
	RejectFriend(ctx *gin.Context)
	GetSentFriendRequests(ctx *gin.Context)
	GetReceivedFriendRequests(ctx *gin.Context)
	GetFriends(ctx *gin.Context)
	LoadConversations(ctx *gin.Context)
	LoadMessages(ctx *gin.Context)
	UploadAvatar(ctx *gin.Context)
	FindUsers(ctx *gin.Context)
	GetUser(ctx *gin.Context)
}

type handler struct {
	userService service.User
	msgService  service.Message
}

func NewHandler(u service.User, m service.Message) Handler {
	return &handler{userService: u, msgService: m}
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

func (h *handler) AddFriend(c *gin.Context) {
	userID := c.Param("userId")
	friendID := c.Param("friendId")

	if userID == "" || friendID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "userId and friendId are required"})
		return
	}
	// Convert userID and friendID to uint
	// Assuming they are valid uints for simplicity
	userIDUint, err := strconv.ParseUint(userID, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid userId"})
		return
	}
	friendIDUint, err := strconv.ParseUint(friendID, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid friendId"})
		return
	}

	err = h.userService.AddFriend(c.Request.Context(), userIDUint, friendIDUint)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "friend request sent successfully"})
}

func (h *handler) CancelFriendRequest(c *gin.Context) {
	userID := c.Param("userId")
	friendID := c.Param("friendId")

	if userID == "" || friendID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "userId and friendId are required"})
		return
	}
	// Convert userID and friendID to uint
	// Assuming they are valid uints for simplicity
	userIDUint, err := strconv.ParseUint(userID, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid userId"})
		return
	}
	friendIDUint, err := strconv.ParseUint(friendID, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid friendId"})
		return
	}

	err = h.userService.CancelFriendRequest(c.Request.Context(), userIDUint, friendIDUint)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "cancel friend request successfully"})
}

func (h *handler) AcceptFriend(c *gin.Context) {
	userID := c.Param("userId")
	friendID := c.Param("friendId")

	if userID == "" || friendID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "userId and friendId are required"})
		return
	}
	// Convert userID and friendID to uint
	// Assuming they are valid uints for simplicity
	userIDUint, err := strconv.ParseUint(userID, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid userId"})
		return
	}
	friendIDUint, err := strconv.ParseUint(friendID, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid friendId"})
		return
	}

	err = h.userService.AcceptFriend(c.Request.Context(), userIDUint, friendIDUint)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "accept friend successfully"})
}

func (h *handler) RejectFriend(c *gin.Context) {
	userID := c.Param("userId")
	friendID := c.Param("friendId")

	if userID == "" || friendID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "userId and friendId are required"})
		return
	}
	// Convert userID and friendID to uint
	// Assuming they are valid uints for simplicity
	userIDUint, err := strconv.ParseUint(userID, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid userId"})
		return
	}
	friendIDUint, err := strconv.ParseUint(friendID, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid friendId"})
		return
	}

	err = h.userService.RejectFriend(c.Request.Context(), userIDUint, friendIDUint)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "reject friend successfully"})
}

func (h *handler) GetSentFriendRequests(c *gin.Context) {
	userID := c.Param("userId")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "userId is required"})
		return
	}
	// Convert userID to uint
	// Assuming they are valid uints for simplicity
	userIDUint, err := strconv.ParseUint(userID, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid userId"})
		return
	}

	friendRequests, err := h.userService.GetSentFriendRequests(c.Request.Context(), userIDUint)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, friendRequests)
}

func (h *handler) GetReceivedFriendRequests(c *gin.Context) {
	userID := c.Param("friendId")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "userId is required"})
		return
	}
	// Convert userID to uint
	// Assuming they are valid uints for simplicity
	userIDUint, err := strconv.ParseUint(userID, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid userId"})
		return
	}

	friendRequests, err := h.userService.GetReceivedFriendRequests(c.Request.Context(), userIDUint)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, friendRequests)
}

func (h *handler) GetFriends(c *gin.Context) {
	userID := c.Param("userId")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "userId is required"})
		return
	}

	search := c.Query("search")
	if search == "" {
		// c.JSON(http.StatusBadRequest, gin.H{"error": "search is required"})
		// return
	}

	// Convert userID to uint
	// Assuming they are valid uints for simplicity
	userIDUint, err := strconv.ParseUint(userID, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid userId"})
		return
	}

	friends, err := h.userService.GetFriends(c.Request.Context(), userIDUint, search)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, friends)
}

func (h *handler) LoadConversations(c *gin.Context) {
	userID := c.Param("userId")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "userId is required"})
		return
	}
	// Convert userID to uint
	// Assuming they are valid uints for simplicity
	userIDUint, err := strconv.ParseUint(userID, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid userId"})
		return
	}

	conversations, err := h.msgService.LoadConversations(c.Request.Context(), userIDUint)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, conversations)
}

func (h *handler) LoadMessages(c *gin.Context) {
	conversationID := c.Param("conversationId")
	if conversationID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "conversationId is required"})
		return
	}
	// Convert conversationID to uint
	// Assuming they are valid uints for simplicity
	conversationIDUint, err := strconv.ParseUint(conversationID, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid conversationId"})
		return
	}

	messages, err := h.msgService.LoadMessages(c.Request.Context(), conversationIDUint)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, messages)
}

func (h *handler) UploadAvatar(c *gin.Context) {
	userID := c.Param("userId")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "userId is required"})
		return
	}
	// Convert userID to uint
	// Assuming they are valid uints for simplicity
	userIDUint, err := strconv.ParseUint(userID, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid userId"})
		return
	}

	fileHeader, err := c.FormFile("avatar")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	filename := fileHeader.Filename
	file, err := fileHeader.Open()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	defer file.Close()

	resp, err := h.userService.UploadAvatar(c.Request.Context(), userIDUint, filename, &file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *handler) FindUsers(c *gin.Context) {
	search := c.Query("search")
	if search == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "search is required"})
		return
	}

	user, err := h.userService.FindUsers(c.Request.Context(), search)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

func (h *handler) GetUser(c *gin.Context) {
	userID := c.Param("userId")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "userId is required"})
		return
	}
	// Convert userID to uint
	// Assuming they are valid uints for simplicity
	userIDUint, err := strconv.ParseUint(userID, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid userId"})
		return
	}

	user, err := h.userService.GetUser(c.Request.Context(), userIDUint)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}
