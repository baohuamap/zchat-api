package repository

import (
	"context"

	"github.com/baohuamap/zchat-api/models"
	"gorm.io/gorm"
)

type FriendshipRepository interface {
	Create(ctx context.Context, friendship models.Friendship) error
	Get(ctx context.Context, id uint) (models.Friendship, error)
	GetByUserID(ctx context.Context, userID uint64) ([]models.Friendship, error)
	GetByFriendID(ctx context.Context, friendID uint64) ([]models.Friendship, error)
	GetByUserIDAndFriendID(ctx context.Context, userID, friendID uint64) (models.Friendship, error)
	Update(ctx context.Context, friendship models.Friendship) error
	Delete(ctx context.Context, id uint) error
}

type friendship struct {
	DB *gorm.DB
}

func NewFriendshipRepository(DB *gorm.DB) FriendshipRepository {
	return &friendship{DB: DB}
}

func (r friendship) Create(ctx context.Context, friendship models.Friendship) error {
	return r.DB.Create(&friendship).Error
}

func (r friendship) Get(ctx context.Context, id uint) (models.Friendship, error) {
	var f models.Friendship
	err := r.DB.First(&f, id).Error
	return f, err
}

func (r friendship) GetByUserID(ctx context.Context, userID uint64) ([]models.Friendship, error) {
	var friendships []models.Friendship
	err := r.DB.Where("user_id = ?", userID).Find(&friendships).Error
	return friendships, err
}

func (r friendship) GetByFriendID(ctx context.Context, friendID uint64) ([]models.Friendship, error) {
	var friendships []models.Friendship
	err := r.DB.Where("friend_id = ?", friendID).Find(&friendships).Error
	return friendships, err
}

func (r friendship) GetByUserIDAndFriendID(ctx context.Context, userID, friendID uint64) (models.Friendship, error) {
	var f models.Friendship
	err := r.DB.Where("user_id = ? AND friend_id = ?", userID, friendID).First(&f).Error
	return f, err
}

func (r friendship) Update(ctx context.Context, friendship models.Friendship) error {
	return r.DB.Save(&friendship).Error
}

func (r friendship) Delete(ctx context.Context, id uint) error {
	return r.DB.Delete(&models.Friendship{}, id).Error
}
