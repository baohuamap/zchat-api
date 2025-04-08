package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	ID        uint64 `json:"id" gorm:"unique;primaryKey;autoIncrement"`
	Username  string `json:"username"`
	Password  string `json:"password"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Avatar    string `json:"avatar"`
	Phone     string `json:"phone" gorm:"unique"`
}
