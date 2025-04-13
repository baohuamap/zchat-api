package ws

import "github.com/baohuamap/zchat-api/models"

type Conversation struct {
	ID      string                  `json:"id"`
	Type    models.ConversationType `json:"type"` // 1: private, 2: group
	Creator uint64                  `json:"creator"`
	Clients map[string]*Client      `json:"clients"`
}

type Hub struct {
	Conversations map[string]*Conversation
	Register      chan *Client
	Unregister    chan *Client
	Broadcast     chan *Message
}

func NewHub() *Hub {
	return &Hub{
		Conversations: make(map[string]*Conversation),
		Register:      make(chan *Client),
		Unregister:    make(chan *Client),
		Broadcast:     make(chan *Message, 5),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case cl := <-h.Register:
			if _, ok := h.Conversations[cl.ConversationID]; ok {
				r := h.Conversations[cl.ConversationID]

				if _, ok := r.Clients[cl.ID]; !ok {
					r.Clients[cl.ID] = cl
				}
			}
		case cl := <-h.Unregister:
			if _, ok := h.Conversations[cl.ConversationID]; ok {
				if _, ok := h.Conversations[cl.ConversationID].Clients[cl.ID]; ok {
					if len(h.Conversations[cl.ConversationID].Clients) != 0 {
						h.Broadcast <- &Message{
							Content:        "user left the chat",
							ConversationID: cl.ConversationID,
							Username:       cl.Username,
						}
					}

					delete(h.Conversations[cl.ConversationID].Clients, cl.ID)
					close(cl.Message)
				}
			}

		case m := <-h.Broadcast:
			if _, ok := h.Conversations[m.ConversationID]; ok {

				for _, cl := range h.Conversations[m.ConversationID].Clients {
					cl.Message <- m
				}
			}
		}
	}
}
