package ws

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"

	"github.com/baohuamap/zchat-api/dto"
	"github.com/baohuamap/zchat-api/models"
	"github.com/baohuamap/zchat-api/repository"
)

type Handler interface {
	CreateConversation(c *gin.Context)
	JoinConversation(c *gin.Context)
	// GetConversations(c *gin.Context)
	GetClients(c *gin.Context)
}

type handler struct {
	hub         *Hub
	msg         repository.MessageRepository
	conv        repository.ConversationRepository
	participant repository.ParticipantRepository
}

func NewHandler(
	h *Hub, conv repository.ConversationRepository, participant repository.ParticipantRepository,
	msg repository.MessageRepository,
) Handler {
	return &handler{
		hub:         h,
		conv:        conv,
		participant: participant,
		msg:         msg,
	}
}

func (h *handler) CreateConversation(c *gin.Context) {
	var req dto.CreateConversationReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	conv := &models.Conversation{
		Type:      req.Type,
		CreatorID: req.CreatorID,
	}

	if err := h.conv.Create(c, conv); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	participants := make([]models.Participant, 0)
	for _, userID := range req.Participants {
		participants = append(participants, models.Participant{
			UserID:         userID,
			ConversationID: conv.ID,
		})
	}
	if err := h.participant.BulkCreate(c, participants); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	convID := strconv.FormatUint(conv.ID, 10)
	h.hub.Conversations[convID] = &Conversation{
		ID:      convID,
		Type:    conv.Type,
		Creator: conv.CreatorID,
		Clients: make(map[string]*Client),
	}

	res := &dto.CreateConversationRes{
		ID: convID,
	}

	c.JSON(http.StatusOK, res)
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (h *handler) JoinConversation(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	conversationID := c.Param("conversationId")
	clientID := c.Query("userId")
	username := c.Query("username")

	// clientIDUint, err := strconv.ParseUint(clientID, 10, 64)
	// if err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
	// 	return
	// }

	// convIDUint, err := strconv.ParseUint(conversationID, 10, 64)
	// if err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid conversation ID"})
	// 	return
	// }

	// participant := &models.Participant{
	// 	UserID:         clientIDUint,
	// 	ConversationID: convIDUint,
	// }
	// if err := h.participant.Create(c, participant); err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	// 	return
	// }

	cl := &Client{
		Conn:            conn,
		Message:         make(chan *Message, 10),
		ID:              clientID,
		ConversationID:  conversationID,
		Username:        username,
		msgRepo:         h.msg,
		convRepo:        h.conv,
		participantRepo: h.participant,
	}

	m := &Message{
		Content:        "A new user has joined the conversation",
		ConversationID: conversationID,
		Username:       username,
	}

	h.hub.Register <- cl
	h.hub.Broadcast <- m

	go cl.writeMessage()
	cl.readMessage(h.hub)
}

// func (h *handler) GetConversations(c *gin.Context) {
// 	conversations := make([]dto.ConversationRes, 0)

// 	for _, r := range h.hub.Conversations {
// 		conversations = append(conversations, dto.ConversationRes{
// 			ID: r.ID,
// 		})
// 	}

// 	c.JSON(http.StatusOK, conversations)
// }

func (h *handler) GetClients(c *gin.Context) {
	var clients []dto.ClientRes
	conversationId := c.Param("conversationId")

	if _, ok := h.hub.Conversations[conversationId]; !ok {
		clients = make([]dto.ClientRes, 0)
		c.JSON(http.StatusOK, clients)
	}

	for _, c := range h.hub.Conversations[conversationId].Clients {
		clients = append(clients, dto.ClientRes{
			ID:       c.ID,
			Username: c.Username,
		})
	}

	c.JSON(http.StatusOK, clients)
}
