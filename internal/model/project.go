package model

import (
	"time"

	"github.com/devkcud/arkhon-foundation/arkhon-api/internal/model/dto"
	"gorm.io/gorm"
)

type Project struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	Name        string `gorm:"not null"`
	Description string `gorm:"default:''"`

	Content any `gorm:"type:jsonb;not null;default:'{}'"`
	Budget  int `gorm:"default:0"`
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
