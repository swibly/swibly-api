package model

import (
	"time"

	"github.com/devkcud/arkhon-foundation/arkhon-api/pkg/language"
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
		InApp int `gorm:"default:1"`
		Email int `gorm:"default:-1"`
	} `gorm:"embedded;embeddedPrefix:notify_"`

	Show struct {
		Profile    int `gorm:"default:1"`
		Image      int `gorm:"default:1"`
		Comments   int `gorm:"default:1"`
		Favorites  int `gorm:"default:1"`
		Projects   int `gorm:"default:1"`
		Components int `gorm:"default:1"`
		Followers  int `gorm:"default:1"`
		Following  int `gorm:"default:1"`
		Inventory  int `gorm:"default:-1"`
		Formations int `gorm:"default:1"`
	} `gorm:"embedded;embeddedPrefix:show_"`

	Country string

	Language language.Language `gorm:"type:enum_language;default:pt"`
}
