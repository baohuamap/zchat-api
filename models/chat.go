package models

import "gorm.io/gorm"

type Chat struct {
	gorm.Model
	SenderID   uint `json:"sender_id"`
	Sender     User `json:"sender"`
	ReceiverID uint `json:"receiver_id"`
	Receiver   User `json:"receiver"`
}
