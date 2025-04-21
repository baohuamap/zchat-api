package repository

import (
	"context"

	"github.com/baohuamap/zchat-api/models"
	"gorm.io/gorm"
)

type ParticipantRepository interface {
	Create(ctx context.Context, participant *models.Participant) error
	BulkCreate(ctx context.Context, participants []models.Participant) error
	Get(ctx context.Context, id uint64) (models.Participant, error)
	GetConversationByParticipants(ctx context.Context, userID uint64) ([]models.Conversation, error)
	GetByUserID(ctx context.Context, userID uint64) ([]models.Participant, error)
	GetByConversationID(ctx context.Context, conversationID uint64) ([]models.Participant, error)
	GetByUserIDAndConversationID(ctx context.Context, userID, conversationID uint64) (models.Participant, error)
	Update(ctx context.Context, participant models.Participant) error
	Delete(ctx context.Context, id uint64) error
}

type participant struct {
	DB *gorm.DB
}

func NewParticipantRepository(DB *gorm.DB) ParticipantRepository {
	return &participant{DB: DB}
}

func (r participant) Create(ctx context.Context, participant *models.Participant) error {
	return r.DB.Create(&participant).Error
}

func (r participant) BulkCreate(ctx context.Context, participants []models.Participant) error {
	return r.DB.Create(&participants).Error
}

func (r participant) Get(ctx context.Context, id uint64) (models.Participant, error) {
	var p models.Participant
	err := r.DB.First(&p, id).Error
	return p, err
}

func (r participant) GetByUserID(ctx context.Context, userID uint64) ([]models.Participant, error) {
	var participants []models.Participant
	err := r.DB.Where("user_id = ?", userID).Find(&participants).Error
	return participants, err
}

func (r participant) GetByConversationID(ctx context.Context, conversationID uint64) ([]models.Participant, error) {
	var participants []models.Participant
	err := r.DB.Where("conversation_id = ?", conversationID).Preload("User").Find(&participants).Error
	return participants, err
}

func (r participant) GetByUserIDAndConversationID(ctx context.Context, userID, conversationID uint64) (models.Participant, error) {
	var p models.Participant
	err := r.DB.Where("user_id = ? AND conversation_id = ?", userID, conversationID).First(&p).Error
	return p, err
}

func (r participant) Update(ctx context.Context, participant models.Participant) error {
	return r.DB.Save(&participant).Error
}

func (r participant) Delete(ctx context.Context, id uint64) error {
	return r.DB.Delete(&models.Participant{}, id).Error
}

func (r participant) GetConversationByParticipants(ctx context.Context, userID uint64) ([]models.Conversation, error) {
	var conversations []models.Conversation
	stmt := r.DB.Table("participants").
		Select("conversations.id, conversations.name, conversations.created_at, MAX(messages.created_at) AS last_message_time").
		Joins("JOIN conversations ON participants.conversation_id = conversations.id").
		Joins("LEFT JOIN messages ON messages.conversation_id = conversations.id").
		Where("participants.user_id = ?", userID).
		Group("conversations.id, conversations.name, conversations.created_at").
		Order("last_message_time DESC").
		Find(&conversations)

	if len(conversations) == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	return conversations, stmt.Error
}
