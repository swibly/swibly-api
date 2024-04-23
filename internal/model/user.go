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

	FirstName string `validator:"required,min=3"`
	LastName  string `validator:"required,min=3"`
	Username  string `validator:"required,username,min=3,max=32" gorm:"unique"`
	Email     string `validator:"required,email" gorm:"unique"`
	Password  string `validator:"required,password,min=12,max=48"`
	// TODO: Implement rest of the fields
}
