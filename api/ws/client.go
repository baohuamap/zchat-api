package ws

import (
	"context"
	"log"
	"strconv"

	"github.com/baohuamap/zchat-api/models"
	"github.com/baohuamap/zchat-api/repository"
	"github.com/gorilla/websocket"
)

type Client struct {
	Conn            *websocket.Conn
	Message         chan *Message
	ID              string `json:"id"`
	ConversationID  string `json:"conversationId"`
	Username        string `json:"username"`
	msgRepo         repository.MessageRepository
	convRepo        repository.ConversationRepository
	participantRepo repository.ParticipantRepository
}

type Message struct {
	Content        string `json:"content"`
	ConversationID string `json:"conversationId"`
	Username       string `json:"username"`
}

func (c *Client) writeMessage() {
	defer func() {
		c.Conn.Close()
	}()

	for {
		message, ok := <-c.Message
		if !ok {
			return
		}

		c.Conn.WriteJSON(message)
	}
}

func (c *Client) readMessage(hub *Hub) {
	defer func() {
		hub.Unregister <- c
		c.Conn.Close()
	}()

	for {
		_, m, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}

		convID, err := strconv.ParseUint(c.ConversationID, 10, 64)
		if err != nil {
			log.Printf("error: %v", err)
			continue
		}

		senderID, err := strconv.ParseUint(c.ID, 10, 64)
		if err != nil {
			log.Printf("error: %v", err)
			continue
		}

		msgObj := &models.Message{
			Content:        string(m),
			ConversationID: convID,
			SenderID:       senderID,
		}
		if err := c.msgRepo.Create(context.Background(), msgObj); err != nil {
			log.Printf("error: %v", err)
			continue
		}

		msg := &Message{
			Content:        string(m),
			ConversationID: c.ConversationID,
			Username:       c.Username,
		}

		hub.Broadcast <- msg
	}
}
