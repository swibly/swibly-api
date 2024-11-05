package model

import (
	"time"

	"github.com/swibly/swibly-api/pkg/notification"
	"gorm.io/gorm"
)

type Notification struct {
	ID        uint `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	Title   string
	Message string

	Type notification.NotificationType `gorm:"type:notification_type;default:'information'"`

	Redirect *string
}

type NotificationUser struct {
	ID        uint `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time

	NotificationID uint `gorm:"not null;index"`
	UserID         uint `gorm:"not null;index"`
}

type NotificationUserRead struct {
	ID        uint `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time

	NotificationID uint `gorm:"not null;index"`
	UserID         uint `gorm:"not null;index"`
}
