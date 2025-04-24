package dto

import (
	"time"

	"github.com/baohuamap/zchat-api/models"
)

type CreateConversationReq struct {
	Type         models.ConversationType `json:"type"` // 1: private, 2: group
	CreatorID    uint64                  `json:"creator_id"`
	Participants []uint64                `json:"participants"`
	Name         string                  `json:"name"`
}

type CreateConversationRes struct {
	ID        string `json:"id"`
	Type      string `json:"type"` // 1: private, 2: group
	CreatorID uint64 `json:"creator_id"`
	Name      string `json:"name"`
}

type ConversationRes struct {
	ID                        uint64            `json:"id"`
	Name                      string            `json:"name"`
	Type                      string            `json:"type"` // 1: private, 2: group
	CreatorID                 uint64            `json:"creator_id"`
	Participants              []ParticipantInfo `json:"participants"`
	Seen                      bool              `json:"seen"`
	LatestMessageID           uint64            `json:"latest_message_id"`
	LatestMessageSenderID     uint64            `json:"latest_message_sender_id"`
	LatestMessageSenderName   string            `json:"latest_message_sender_name"`
	LatestMessageSenderAvatar string            `json:"latest_message_sender_avatar"`
	LatestMessageContent      string            `json:"latest_message_content"`
	LatestMessageCreatedAt    time.Time         `json:"latest_message_created_at"`
}

type ConversationListRes struct {
	Conversations []ConversationRes `json:"conversations"`
}

type ParticipantInfo struct {
	ID        uint64 `json:"id"`
	Phone     string `json:"phone"`
	Username  string `json:"username"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Avatar    string `json:"avatar"`
}

type AddParticipantsReq struct {
	Participants []uint64 `json:"participants"`
}
