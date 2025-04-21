package models

import "gorm.io/gorm"

type Conversation struct {
	gorm.Model
	ID        uint64           `gorm:"primaryKey autoIncrement:true" json:"id"`
	Name      string           `gorm:"null" json:"name"`                            // Name of the conversation
	Type      ConversationType `gorm:"type:conversation_type;not null" json:"type"` // Enum: 'private', 'group'
	CreatorID uint64           `gorm:"null" json:"creator_id"`
	Creator   User             `gorm:"foreignKey:CreatorID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"creator"`
	Seen      bool             `gorm:"default:false" json:"seen"`
}

type ConversationType string

const (
	ConversationTypePrivate ConversationType = "private"
	ConversationTypeGroup   ConversationType = "group"
)
