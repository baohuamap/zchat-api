package ws

import (
	"log"
	"strconv"

	"github.com/baohuamap/zchat-api/models"
	"github.com/gorilla/websocket"
)

type Client struct {
	Conn           *websocket.Conn
	Message        chan *Message
	ID             string `json:"id"`
	ConversationID string `json:"conversationId"`
	Username       string `json:"username"`
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

		msg := &models.Message{
			Content:        string(m),
			ConversationID: convID,
			SenderID:       c.ID,
		}

		msg := &Message{
			Content:        string(m),
			ConversationID: c.ConversationID,
			Username:       c.Username,
		}

		hub.Broadcast <- msg
	}
}
