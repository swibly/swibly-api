package model

import (
	"time"

	"gorm.io/gorm"
)

var (
	ROLE_ADMIN = "admin"
	ROLE_USER  = "user"
)

type User struct {
	// NOTE: Not using gorm.Model since it's properties cannot be accessed directly
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	Role string `gorm:"default:user"`

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
	// TODO: Add comments, the last implementation was wacky to say the least.

	Notification struct {
		InApp int `gorm:"default:1"  json:"inapp"`
		Email int `gorm:"default:-1" json:"email"`
		// SMS int `gorm:"default:-1"` // NOTE: Not sure if we want to send SMS, it can get expensive
	} `gorm:"embedded;embeddedPrefix:notify_"`

	Show struct {
		Profile    int `gorm:"default:1" json:"profile"`
		Image      int `gorm:"default:1" json:"image"`
		Comments   int `gorm:"default:1" json:"comments"`
		Favorites  int `gorm:"default:1" json:"favorites"`
		Projects   int `gorm:"default:1" json:"projects"`
		Components int `gorm:"default:1" json:"components"`
		Followers  int `gorm:"default:1" json:"followers"`
		Following  int `gorm:"default:1" json:"following"`
		Inventory  int `gorm:"default:-1" json:"inventory"`
		Formations int `gorm:"default:1" json:"formations"`
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

type ProfileUpdate struct {
	FirstName string `validate:"omitempty,min=3"                 json:"firstname"`
	LastName  string `validate:"omitempty,min=3"                 json:"lastname"`
	Username  string `validate:"omitempty,username,min=3,max=32" json:"username"`
	Email     string `validate:"omitempty,email"                 json:"email"`
	Bio       string `validate:"omitempty"                       json:"bio"`
	Show      struct {
		Profile    int `json:"profile"`
		Image      int `json:"image"`
		Comments   int `json:"comments"`
		Favorites  int `json:"favorites"`
		Projects   int `json:"projects"`
		Components int `json:"components"`
		Followers  int `json:"followers"`
		Following  int `json:"following"`
		Inventory  int `json:"inventory"`
		Formations int `json:"formations"`
	} `json:"show"`
	Notification struct {
		InApp int `json:"inapp"`
		Email int `json:"email"`
	} `json:"notification"`
}
