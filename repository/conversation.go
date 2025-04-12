package repository

import (
	"context"

	"github.com/baohuamap/zchat-api/models"
	"gorm.io/gorm"
)

type ConversationRepository interface {
	Create(ctx context.Context, conversation models.Conversation) error
	Get(ctx context.Context, id uint64) (models.Conversation, error)
	Update(ctx context.Context, conversation models.Conversation) error
	Delete(ctx context.Context, id uint64) error
}

type conversation struct {
	DB *gorm.DB
}

func NewConversationRepository(DB *gorm.DB) ConversationRepository {
	return &conversation{DB: DB}
}

func (r conversation) Create(ctx context.Context, conversation models.Conversation) error {
	return r.DB.Create(&conversation).Error
}

func (r conversation) Get(ctx context.Context, id uint64) (models.Conversation, error) {
	var c models.Conversation
	err := r.DB.First(&c, id).Error
	return c, err
}

func (r conversation) Update(ctx context.Context, conversation models.Conversation) error {
	return r.DB.Save(&conversation).Error
}

func (r conversation) Delete(ctx context.Context, id uint64) error {
	return r.DB.Delete(&models.Conversation{}, id).Error
}
