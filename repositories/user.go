package repositories

import (
	"github.com/baohuamap/zchat-api/models"

	"gorm.io/gorm"
)

type UserRepository interface {
	Create(user models.User) error
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
