package model

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	// NOTE: Not using gorm.Model since it's properties cannot be accessed directly
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	FirstName string `validate:"required,min=3"`
	LastName  string `validate:"required,min=3"`
	Username  string `validate:"required,username,min=3,max=32" gorm:"unique"`
	Email     string `validate:"required,email" gorm:"unique"`
	Password  string `validate:"required,password,min=12,max=48"`

	// TODO: Implement rest of the fields
}

type UserRegister struct {
	FirstName string `validate:"required,min=3" json:"firstname"`
	LastName  string `validate:"required,min=3" json:"lastname"`
	Username  string `validate:"required,username,min=3,max=32" json:"username" gorm:"unique"`
	Email     string `validate:"required,email" json:"email" gorm:"unique"`
	Password  string `validate:"required,password,min=12,max=48" json:"password"`
}
