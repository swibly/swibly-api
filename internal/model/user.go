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

  // TODO: Add Role

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
