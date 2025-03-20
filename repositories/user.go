package repositories

import (
	"context"

	"github.com/baohuamap/zchat-api/models"
	"gorm.io/gorm"
)

type UserRepository interface {
	Create(user models.User) error
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
}

type user struct {
	DB *gorm.DB
}

func ProvideUserRepository(DB *gorm.DB) UserRepository {
	return &user{DB: DB}
}

func (r user) Create(user models.User) error {
	return r.DB.Create(&user).Error
}

func (r user) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	err := r.DB.WithContext(ctx).Where("email = ?", email).First(&user).Error
	return &user, err

}
