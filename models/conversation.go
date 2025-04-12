package models

import "gorm.io/gorm"

type Conversation struct {
	gorm.Model
	ID        uint64 `gorm:"primaryKey autoIncrement:true" json:"id"`
	Type      string `gorm:"not null" json:"type"` // 1: private, 2: group
	CreatorID uint64 `gorm:"null" json:"creator_id"`
	Creator   User   `gorm:"foreignKey:CreatorID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"creator"`
}
