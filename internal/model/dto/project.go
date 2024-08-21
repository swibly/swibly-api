package dto

type ProjectCreation struct {
	Owner string `json:"-"`

	Name        string `validate:"required"  json:"name"`
	Description string `validate:"omitempty" json:"description"`

	Content   map[string]any `validate:"omitempty" json:"content" gorm:"type:json"`
	Thumbnail string         `validate:"omitempty" json:"thumbnail"`
	Budget    int            `validate:"omitempty" json:"budget"`
}
