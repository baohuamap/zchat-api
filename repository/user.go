package repository

import (
	"context"

	"github.com/baohuamap/zchat-api/models"

	"gorm.io/gorm"
)

type UserRepository interface {
	Create(ctx context.Context, user models.User) error
	Get(ctx context.Context, id uint) (models.User, error)
	GetByEmail(ctx context.Context, email string) (models.User, error)
	GetByPhone(ctx context.Context, phone string) (models.User, error)
}

type user struct {
	DB *gorm.DB
}

func NewUserRepository(DB *gorm.DB) UserRepository {
	return &user{DB: DB}
}

func (r user) Create(ctx context.Context, user models.User) error {
	return r.DB.Create(&user).Error
}

func (r user) Get(ctx context.Context, id uint) (models.User, error) {
	var u models.User
	err := r.DB.First(&u, id).Error
	return u, err
}

func (r user) GetByEmail(ctx context.Context, email string) (models.User, error) {
	var u models.User
	err := r.DB.Where("email = ?", email).First(&u).Error
	return u, err
}

func (r user) GetByPhone(ctx context.Context, phone string) (models.User, error) {
	var u models.User
	err := r.DB.Where("phone = ?", phone).First(&u).Error
	return u, err
}
