package dto

import (
	"time"

	"github.com/swibly/swibly-api/pkg/notification"
)

type CreateNotification struct {
	Title    string                        `json:"title"    validate:"required,max=255"`
	Message  string                        `json:"message"  validate:"required"`
	Type     notification.NotificationType `json:"type"     validate:"required,mustbenotificationtype"`
	Redirect *string                       `json:"redirect" validate:"omitempty,url"`
}

type NotificationInfo struct {
	ID        uint      `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	Title   string `json:"title"`
	Message string `json:"message"`

	Type     notification.NotificationType `json:"type"`
	Redirect *string                       `json:"redirect"`

	ReadAt *time.Time `json:"read_at"`
	IsRead bool       `json:"is_read"`
}
