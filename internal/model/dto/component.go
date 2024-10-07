package dto

import (
	"time"

	"github.com/swibly/swibly-api/pkg/utils"
	"gorm.io/gorm"
)

type ComponentCreation struct {
	Name        string `validate:"required"  json:"name"`
	Description string `validate:"omitempty" json:"description"`

	Content any `validate:"required" json:"content"`

	Price int `validate:"omitempty" json:"price"`

	OwnerID uint `json:"-"` // Set using JWT

	Public bool `json:"-"` // Set in an URL query param
}

type ComponentUpdate struct {
	Name        *string `validate:"omitempty,min=3,max=32" json:"name"`
	Description *string `validate:"omitempty,min=3,max=5000" json:"description"`

	Content *any `validate:"omitempty" json:"content"`

	Price *int `validate:"omitempty,min=0,max=1000000" json:"price"`

	Public *bool `validate:"omitempty" json:"public"`
}

type ComponentInfoJSON struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`

	Name        string `json:"name"`
	Description string `json:"description"`

	Content utils.JSON `json:"content"`

	OwnerID             uint   `json:"owner_id"`
	OwnerUsername       string `json:"owner_username"`
	OwnerProfilePicture string `json:"owner_profile_picture"`

	Price     int  `json:"price"`
	PaidPrice *int `json:"paid_price"`
	SellPrice *int `json:"sell_price"`

	IsPublic bool `json:"is_public"`

	Holders    int64 `json:"holders"`
	Bought     bool  `json:"bought"`
	TotalSells int64 `json:"total_sells"`
}

type ComponentInfo struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`

	Name        string `json:"name"`
	Description string `json:"description"`

	Content any `json:"content"`

	OwnerID             uint   `json:"owner_id"`
	OwnerUsername       string `json:"owner_username"`
	OwnerProfilePicture string `json:"owner_profile_picture"`

	Price     int  `json:"price"`
	PaidPrice *int `json:"paid_price"`
	SellPrice *int `json:"sell_price"`

	IsPublic bool `json:"is_public"`

	Holders int64 `json:"holders"`
	Bought  bool  `json:"bought"`
	TotalSells int64 `json:"total_sells"`
}
