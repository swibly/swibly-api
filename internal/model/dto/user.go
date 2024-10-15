package dto

import (
	"time"

	"github.com/swibly/swibly-api/pkg/language"
	"github.com/swibly/swibly-api/pkg/utils"
)

type UserRegister struct {
	FirstName string `validate:"required,min=3,max=32"          json:"firstname"`
	LastName  string `validate:"required,min=3,max=32"          json:"lastname"`
	Username  string `validate:"required,min=3,max=32,username" json:"username"`
	Email     string `validate:"required,email"                 json:"email"`
	Password  string `validate:"required,password"              json:"password"`
}

type UserLogin struct {
	Username string `validate:"omitempty" json:"username"`
	Email    string `validate:"omitempty" json:"email"`
	Password string `validate:"required"  json:"password"`
}

type UserUpdate struct {
	FirstName *string `validate:"omitempty,min=3,max=32"          json:"firstname"`
	LastName  *string `validate:"omitempty,min=3,max=32"          json:"lastname"`
	Username  *string `validate:"omitempty,min=3,max=32,username" json:"username"`

	Bio      *string `validate:"omitempty,max=480" json:"bio"`
	Verified *bool   `validate:"omitempty"         json:"verified"`

	Email    *string `validate:"omitempty,email"    json:"email"`
	Password *string `validate:"omitempty,password" json:"password"`

	Notification struct {
		InApp *bool `validate:"omitempty" json:"inapp"`
		Email *bool `validate:"omitempty" json:"email"`
	} `validate:"omitempty" json:"notify" gorm:"embedded;embeddedPrefix:notify_"`

	Show struct {
		Profile    *bool `validate:"omitempty" json:"profile"`
		Image      *bool `validate:"omitempty" json:"image"`
		Comments   *bool `validate:"omitempty" json:"comments"`
		Favorites  *bool `validate:"omitempty" json:"favorites"`
		Projects   *bool `validate:"omitempty" json:"projects"`
		Components *bool `validate:"omitempty" json:"components"`
		Followers  *bool `validate:"omitempty" json:"followers"`
		Following  *bool `validate:"omitempty" json:"following"`
		Inventory  *bool `validate:"omitempty" json:"inventory"`
		Formations *bool `validate:"omitempty" json:"formations"`
	} `validate:"omitempty" json:"show" gorm:"embedded;embeddedPrefix:show_"`

	Country *string `validate:"omitempty,max=40" json:"country"`

	Language *language.Language `validate:"omitempty,mustbesupportedlanguage" json:"language"`
}

type UserProfile struct {
	ID        uint      `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
	Username  string `json:"username"`
	Email     string `json:"email"`

	Bio      string `json:"bio"`
	Verified bool   `json:"verified"`

	XP      uint64 `json:"xp"`
	Arkhoin uint64 `json:"arkhoins"`

	Followers int64 `gorm:"-" json:"followers"`
	Following int64 `gorm:"-" json:"following"`

	Notification struct {
		InApp bool `json:"inapp"`
		Email bool `json:"email"`
	} `gorm:"embedded;embeddedPrefix:notify_" json:"notification"`

	Show struct {
		Profile    bool `json:"profile"`
		Image      bool `json:"image"`
		Comments   bool `json:"comments"`
		Favorites  bool `json:"favorites"`
		Projects   bool `json:"projects"`
		Components bool `json:"components"`
		Followers  bool `json:"followers"`
		Following  bool `json:"following"`
		Inventory  bool `json:"inventory"`
		Formations bool `json:"formations"`
	} `gorm:"embedded;embeddedPrefix:show_" json:"show"`

	Country  string `json:"country"`
	Language string `json:"language"`

	Permissions []string `gorm:"-" json:"permissions"`

	ProfilePicture string `json:"pfp"`
}

type UserInfoLite struct {
	ID             uint   `json:"id"`
	Username       string `json:"username"`
	ProfilePicture string `json:"pfp"`
}

type UserShow struct {
	Profile    bool
	Image      bool
	Comments   bool
	Favorites  bool
	Projects   bool
	Components bool
	Followers  bool
	Following  bool
	Inventory  bool
	Formations bool
}

func (u *UserProfile) HasPermissions(permissions ...string) bool {
	return utils.HasPermissions(u.Permissions, permissions...)
}
