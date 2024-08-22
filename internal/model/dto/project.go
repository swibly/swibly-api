package dto

import "time"

type ProjectCreation struct {
	Owner string `json:"-"`

	Name        string `validate:"required"  json:"name"`
	Description string `validate:"omitempty" json:"description"`

	Content   map[string]any `validate:"omitempty" json:"content" gorm:"type:json"`
	Thumbnail string         `validate:"omitempty" json:"thumbnail"`
	Budget    int            `validate:"omitempty" json:"budget"`
}

type ProjectInformation struct {
	ID        uint      `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	Owner string `json:"owner"`

	Name        string `json:"name"`
	Description string `json:"description"`

	Thumbnail string `json:"thumbnail"`
	Budget    int    `json:"budget"`

	Published bool `json:"published"`
	Upstream  uint `json:"upstream"`
}
