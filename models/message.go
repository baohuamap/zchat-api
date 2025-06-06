package models

import "gorm.io/gorm"

type Message struct {
	gorm.Model
	ID             uint64       `gorm:"primaryKey" autoIncrement:"true" json:"id"`
	Content        string       `gorm:"not null" json:"content"`
	ConversationID uint64       `gorm:"not null" json:"conversation_id"`
	Conversation   Conversation `gorm:"foreignKey:ConversationID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"conversation"`
	SenderID       uint64       `gorm:"not null" json:"sender_id"`
	Sender         User         `gorm:"foreignKey:SenderID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"sender"`
}
