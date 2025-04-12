package dto

type CreateConversationReq struct {
	Type         string   `json:"type"` // 1: private, 2: group
	CreatorID    uint64   `json:"creatorId"`
	Participants []uint64 `json:"participants"`
}

type ConversationRes struct {
	ID string `json:"id"`
}
