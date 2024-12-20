package model

import (
	"time"

	"github.com/swibly/swibly-api/internal/model/dto"
	"gorm.io/gorm"
)

type Project struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	Name        string `gorm:"not null"`
	Description string `gorm:"default:''"`

	BannerURL string `gorm:"default:''"`

	Width  int `gorm:"not null;default:30"`
	Height int `gorm:"not null;default:30"`

	Content any `gorm:"type:jsonb;not null;default:'{}'"`
	Budget  int `gorm:"default:0"`

	Fork *uint `gorm:"index"`
}

type ProjectOwner struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time

	ProjectID uint `gorm:"unique;index;not null;constraint:OnDelete:CASCADE;"`
	UserID    uint `gorm:"index;not null;constraint:OnDelete:CASCADE;"`
}

type ProjectPublication struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time

	ProjectID uint `gorm:"unique;index;not null;constraint:OnDelete:CASCADE;"`
}

type ProjectUserFavorite struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time

	ProjectID uint `gorm:"index;not null;constraint:OnDelete:CASCADE;"`
	UserID    uint `gorm:"index;not null;constraint:OnDelete:CASCADE;"`
}

type ProjectUserPermission struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time

	ProjectID uint `gorm:"index;not null;constraint:OnDelete:CASCADE;"`
	UserID    uint `gorm:"index;not null;constraint:OnDelete:CASCADE;"`

	Allow dto.Allow `gorm:"embedded;embeddedPrefix:allow_"`
}
