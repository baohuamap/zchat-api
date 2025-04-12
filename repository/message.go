package repository

import (
	"context"

	"github.com/baohuamap/zchat-api/models"

	"gorm.io/gorm"
)

type MessageRepository interface {
	Create(ctx context.Context, user models.Message) error
	Get(ctx context.Context, id uint) (models.Message, error)
	GetByConversationID(ctx context.Context, conversationID uint64) ([]models.Message, error)
	GetBySenderID(ctx context.Context, userID uint64) ([]models.Message, error)
	GetBySenderIDAndConversationID(ctx context.Context, userID, conversationID uint64) ([]models.Message, error)
	Update(ctx context.Context, message models.Message) error
	Delete(ctx context.Context, id uint) error
}

type message struct {
	DB *gorm.DB
}

func NewMessageRepository(DB *gorm.DB) MessageRepository {
	return &message{DB: DB}
}

func (r message) Create(ctx context.Context, msg models.Message) error {
	return r.DB.Create(&msg).Error
}

func (r message) Get(ctx context.Context, id uint) (models.Message, error) {
	var m models.Message
	err := r.DB.First(&m, id).Error
	return m, err
}

func (r message) GetByConversationID(ctx context.Context, conversationID uint64) ([]models.Message, error) {
	var messages []models.Message
	err := r.DB.Where("conversation_id = ?", conversationID).Find(&messages).Error
	return messages, err
}

func (r message) GetBySenderID(ctx context.Context, userID uint64) ([]models.Message, error) {
	var messages []models.Message
	err := r.DB.Where("sender_id = ?", userID).Find(&messages).Error
	return messages, err
}

func (r message) GetBySenderIDAndConversationID(ctx context.Context, userID, conversationID uint64) ([]models.Message, error) {
	var messages []models.Message
	err := r.DB.Where("sender_id = ? AND conversation_id = ?", userID, conversationID).Find(&messages).Error
	return messages, err
}

func (r message) Update(ctx context.Context, message models.Message) error {
	return r.DB.Save(&message).Error
}

func (r message) Delete(ctx context.Context, id uint) error {
	return r.DB.Delete(&models.Message{}, id).Error
}
