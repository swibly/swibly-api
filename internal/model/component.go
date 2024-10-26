package model

import (
	"time"

	"gorm.io/gorm"
)

type Component struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	Name        string `gorm:"not null"`
	Description string `gorm:"default:''"`

	Content any `gorm:"type:jsonb;not null"`

	Price  int `gorm:"default:0"`
	Budget int `gorm:"default:0"`
}

type ComponentOwner struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time

	ComponentID uint  `gorm:"index;unique;not null;constraint:OnDelete:CASCADE;"`
	UserID      *uint `gorm:"index;constraint:OnDelete:CASCADE;"`
}

type ComponentPublication struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time

	ComponentID uint `gorm:"index;unique;not null;constraint:OnDelete:CASCADE;"`
}

type ComponentHolder struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time

	ComponentID uint `gorm:"index;not null;constraint:OnDelete:CASCADE;"`
	UserID      uint `gorm:"index;not null;constraint:OnDelete:CASCADE;"`

	PricePaid int `gorm:"default:0"`
}
