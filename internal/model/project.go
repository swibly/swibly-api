package model

import "time"

type Project struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time

	Name        string `gorm:"not null"`
	Description string `gorm:"default:''"`

	Content any `gorm:"type:jsonb;not null;default:'{}'"`
	Budget  int `gorm:"default:0"`
}

type ProjectPublication struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time

	ProjectID uint `gorm:"unique;index;not null"`
}

type ProjectUserFavorite struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time

	ProjectID uint `gorm:"index;not null"`
	UserID    uint `gorm:"index;not null"`
}

type ProjectUserPermission struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time

	ProjectID uint `gorm:"index;not null"`
	UserID    uint `gorm:"index;not null"`

	Allow struct { // If user is admin it will ignore all fields
		View    bool `gorm:"not null;default:false"` // Will be ignored if project is public
		Edit    bool `gorm:"not null;default:false"`
		Delete  bool `gorm:"not null;default:false"`
		Publish bool `gorm:"not null;default:false"`
		Share   bool `gorm:"not null;default:false"` // Will be ignored if project is public
		Manage  struct {
			Users    bool `gorm:"not null;default:false"`
			Metadata bool `gorm:"not null;default:false"`
		} `gorm:"embedded;embeddedPrefix:manage_"`
	} `gorm:"embedded;embeddedPrefix:allow_"`
}
