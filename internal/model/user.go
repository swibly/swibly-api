package model

import (
	"time"

	"github.com/swibly/swibly-api/pkg/language"
)

type User struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time

	FirstName string
	LastName  string
	Bio       string
	Verified  bool

	Username string `gorm:"unique"`
	Email    string `gorm:"unique"`
	Password string

	XP      uint64 `gorm:"default:500"`
	Arkhoin uint64 `gorm:"default:1000"`

	Notification struct {
		InApp bool `gorm:"default:true"`
		Email bool `gorm:"default:false"`
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

	Country string

	Language language.Language `gorm:"type:enum_language;default:pt"`

	ProfilePicture string
}
