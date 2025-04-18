package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	ID        uint64 `gorm:"primaryKey autoIncrement:true" json:"id"`
	Username  string `json:"username"`
	Password  string `json:"password"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Avatar    string `json:"avatar"`
	Phone     string `json:"phone" gorm:"unique"`
}
