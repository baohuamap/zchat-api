package dto

type Message struct {
	Content        string `json:"content"`
	ConversationID string `json:"conversationId"`
	SenderID       uint64 `json:"senderId"`
}
