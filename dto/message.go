package dto

import "time"

type Message struct {
	Content        string `json:"content"`
	ConversationID string `json:"conversationId"`
	SenderID       uint64 `json:"senderId"`
}

type MessageRes struct {
	Content        string    `json:"content"`
	CreateAt       time.Time `json:"createAt"`
	ConversationID uint64    `json:"conversationId"`
	SenderID       uint64    `json:"senderId"`
}

type MessageListRes struct {
	Messages []MessageRes `json:"messages"`
}

type SeenMessagesReq struct {
	UserID uint64 `json:"user_id"`
}
