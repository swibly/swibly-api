package model

import "time"

type Project struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time

	Owner string

	Name        string
	Description string

	Content   any    `gorm:"serializer:json;type:json;default:'{}'"`
	Thumbnail string `gorm:"default:''"`
	Budget    int    `gorm:"default:0"` // In cents

	Published bool `gorm:"default:false"`

	Upstream uint `gorm:"index"`
}

type ProjectFavorite struct {
	ProjectID uint `gorm:"index"`
	UserID    uint `gorm:"index"`
}
