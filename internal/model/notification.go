package model

import (
	"time"

	"github.com/swibly/swibly-api/pkg/notification"
)

type Notification struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time

	Title string
	Body  string

	Type notification.NotificationType `gorm:"index;type:notification_type"`

	Redirect *string
}

type NotificationUser struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time

	UserID         uint `gorm:"index"`
	NotificationID uint `gorm:"index"`
}

type NotificationUserRead struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time

	UserID         uint `gorm:"index"`
	NotificationID uint `gorm:"index"`

	ReadAt *time.Time `gorm:"index"`
}
