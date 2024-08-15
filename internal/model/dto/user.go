package dto

import (
	"time"

	"github.com/devkcud/arkhon-foundation/arkhon-api/pkg/language"
)

type UserRegister struct {
	FirstName string `validate:"required,min=3"    json:"firstname"`
	LastName  string `validate:"required,min=3"    json:"lastname"`
	Username  string `validate:"required,username" json:"username"`
	Email     string `validate:"required,email"    json:"email"`
	Password  string `validate:"required,password" json:"password"`
}

type UserLogin struct {
	Username string `validate:"omitempty" json:"username"`
	Email    string `validate:"omitempty" json:"email"`
	Password string `validate:"required"  json:"password"`
}

type UserUpdate struct {
	FirstName string `validate:"omitempty,min=3"    json:"firstname"`
	LastName  string `validate:"omitempty,min=3"    json:"lastname"`
	Username  string `validate:"omitempty,username" json:"username"`

	Bio      string `validate:"omitempty,max=480" json:"bio"`
	Verified bool   `validate:"omitempty"         json:"verified"`

	// NOTE: XP and Arkhoins doesn't make sense to update here

	Email    string `validate:"omitempty,email"    json:"email"`
	Password string `validate:"omitempty,password" json:"password"`

	Notification struct {
		InApp int `validate:"omitempty,mustbenumericalboolean" json:"inapp"`
		Email int `validate:"omitempty,mustbenumericalboolean" json:"email"`
	} `validate:"omitempty" json:"notify" gorm:"embedded;embeddedPrefix:notify_"`

	Show struct {
		Profile    int `validate:"omitempty,mustbenumericalboolean" json:"profile"`
		Image      int `validate:"omitempty,mustbenumericalboolean" json:"image"`
		Comments   int `validate:"omitempty,mustbenumericalboolean" json:"comments"`
		Favorites  int `validate:"omitempty,mustbenumericalboolean" json:"favorites"`
		Projects   int `validate:"omitempty,mustbenumericalboolean" json:"projects"`
		Components int `validate:"omitempty,mustbenumericalboolean" json:"components"`
		Followers  int `validate:"omitempty,mustbenumericalboolean" json:"followers"`
		Following  int `validate:"omitempty,mustbenumericalboolean" json:"following"`
		Inventory  int `validate:"omitempty,mustbenumericalboolean" json:"inventory"`
		Formations int `validate:"omitempty,mustbenumericalboolean" json:"formations"`
  } `validate:"omitempty" json:"show" gorm:"embedded;embeddedPrefix:show_"`

	Country string `validate:"omitempty" json:"country"`

	Language language.Language `validate:"omitempty,mustbesupportedlanguage" json:"language"`
}

type UserProfile struct {
	ID        uint      `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
	Username  string `json:"username"`

	Bio      string `json:"bio"`
	Verified bool   `json:"verified"`

	XP      uint64 `json:"xp"`
	Arkhoin uint64 `json:"arkhoins"`

	Followers int64 `gorm:"-" json:"followers"`
	Following int64 `gorm:"-" json:"following"`

	Notification struct {
		InApp int `json:"inapp"`
		Email int `json:"email"`
	} `gorm:"embedded;embeddedPrefix:notify_" json:"notification"`

	Show struct {
		Profile    int `json:"profile"`
		Image      int `json:"image"`
		Comments   int `json:"comments"`
		Favorites  int `json:"favorites"`
		Projects   int `json:"projects"`
		Components int `json:"components"`
		Followers  int `json:"followers"`
		Following  int `json:"following"`
		Inventory  int `json:"inventory"`
		Formations int `json:"formations"`
	} `gorm:"embedded;embeddedPrefix:show_" json:"show"`

	Country  string `json:"country"`
	Language string `json:"language"`

	Permissions []string `gorm:"-" json:"permissions"`
}
