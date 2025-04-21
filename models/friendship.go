package models

import (
	"gorm.io/gorm"
)

type Friendship struct {
	gorm.Model
	ID       uint64 `json:"id" gorm:"primaryKey"`
	UserID   uint64 `json:"user_id" gorm:"not null"`
	FriendID uint64 `json:"friend_id" gorm:"not null"`
	Status   string `json:"status" gorm:"not null"` // "pending", "accepted", "blocked"
}
