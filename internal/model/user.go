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

	FirstName string
	LastName  string
	Bio       string
	Verified  bool

	Username string `gorm:"unique"`
	Email    string `gorm:"unique"`
	Password string

	XP      uint64 `gorm:"default:500"`
	Arkhoin uint64 `gorm:"default:1000"`

	// TODO: Add followers and following

	Comments []Comment `gorm:"foreignKey:OwnerID"`

	Notification struct {
		InApp bool `gorm:"default:true"`
		Email bool `gorm:"default:false"`
		// SMS bool `gorm:"default:false"` // NOTE: Not sure if we want to send SMS, it can get expensive
	} `gorm:"embedded;embeddedPrefix:notify_"`

	Show struct {
		Profile    bool `gorm:"default:true"`
		Image      bool `gorm:"default:true"`
		Comments   bool `gorm:"default:true"`
		Favorites  bool `gorm:"default:true"`
		Projects   bool `gorm:"default:true"`
		Components bool `gorm:"default:true"`
		Followers  bool `gorm:"default:true"`
		Following  bool `gorm:"default:true"`
		Inventory  bool `gorm:"default:false"`
		Formations bool `gorm:"default:true"`
	} `gorm:"embedded;embeddedPrefix:show_"`

	// TODO: Implement enums Language, Theme and Country (country shouldnt be an enum)
}

type UserRegister struct {
	FirstName string `validate:"required,min=3"                  json:"firstname"`
	LastName  string `validate:"required,min=3"                  json:"lastname"`
	Username  string `validate:"required,username,min=3,max=32"  json:"username"`
	Email     string `validate:"required,email"                  json:"email"`
	Password  string `validate:"required,password,min=12,max=48" json:"password"`
}

type UserLogin struct {
	Username string `validate:"omitempty,username,min=3,max=32" json:"username"`
	Email    string `validate:"omitempty,email"                 json:"email"`
	Password string `validate:"required,password,min=12,max=48" json:"password"`
}
