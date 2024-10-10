package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PasswordResetKey struct {
	Key    uuid.UUID `gorm:"primarykey;type:uuid;default:gen_random_uuid()"`
	UserID uint      `gorm:"unique;not null"`

	ExpiresAt time.Time `gorm:"not null"`
}

func (p *PasswordResetKey) BeforeCreate(tx *gorm.DB) (err error) {
	p.ExpiresAt = time.Now().Add(24 * time.Hour)
	return
}
