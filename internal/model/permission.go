package model

import (
	"time"
)

type Permission struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time

	Name string `gorm:"unique"`
}

type UserPermission struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time

	UserID       uint `gorm:"Index"`
	PermissionID uint `gorm:"Index"`
}
