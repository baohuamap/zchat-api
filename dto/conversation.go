package dto

import "github.com/baohuamap/zchat-api/models"

type CreateConversationReq struct {
	Type         models.ConversationType `json:"type"` // 1: private, 2: group
	CreatorID    uint64                  `json:"creator_id"`
	Participants []uint64                `json:"participants"`
}

type CreateConversationRes struct {
	ID string `json:"id"`
}

type ConversationRes struct {
	ID           uint64                  `json:"id"`
	Type         models.ConversationType `json:"type"` // 1: private, 2: group
	CreatorID    uint64                  `json:"creator_id"`
	Participants []uint64                `json:"participants"`
}

type ConversationListRes struct {
	Conversations []ConversationRes `json:"conversations"`
}
