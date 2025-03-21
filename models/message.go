package models

import "gorm.io/gorm"

type Message struct {
	gorm.Model
	Message string `json:"message"`
	ChatID  uint   `json:"chat_id"`
	Chat    Chat   `json:"chat"`
	IsRead  bool   `json:"is_read"`
}
