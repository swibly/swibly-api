package model

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Fullname string
	Username string `gorm:"unique"`
	Email    string `gorm:"unique"`
	Password string
}
