package models

import "gorm.io/gorm"

type Participant struct {
	gorm.Model
	ID             uint64       `gorm:"primaryKey" autoIncrement:"true" json:"id"`
	UserID         uint64       `gorm:"not null" json:"user_id"`
	User           User         `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"user"`
	ConversationID uint64       `gorm:"not null" json:"conversation_id"`
	Conversation   Conversation `gorm:"foreignKey:ConversationID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"conversation"`
}
