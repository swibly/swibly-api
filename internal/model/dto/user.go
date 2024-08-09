package dto

import "github.com/devkcud/arkhon-foundation/arkhon-api/pkg/language"

type UserRegister struct {
	FirstName string `validate:"required,min=3"                  json:"firstname"`
	LastName  string `validate:"required,min=3"                  json:"lastname"`
	Username  string `validate:"required,username,min=3,max=32"  json:"username"`
	Email     string `validate:"required,email"                  json:"email"`
	Password  string `validate:"required,password" json:"password"`
}

type UserLogin struct {
	Username string `validate:"omitempty,username,min=3,max=32" json:"username"`
	Email    string `validate:"omitempty,email"                 json:"email"`
	Password string `validate:"required,password" json:"password"`
}

type UserUpdate struct {
	FirstName string `validate:"omitempty,min=3"                 json:"firstname"`
	LastName  string `validate:"omitempty,min=3"                 json:"lastname"`
	Username  string `validate:"omitempty,username,min=3,max=32" json:"username"`

	Bio      string `validate:"omitempty,max=480" json:"bio"`
	Verified bool   `validate:"omitempty"         json:"verified"`

	// NOTE: XP and Arkhoins doesn't make sense to update here

	Email    string `validate:"omitempty,email"                  json:"email"`
	Password string `validate:"omitempty,password" json:"password"`

	Notification struct {
		InApp int `validate:"omitempty,mustbenumericalboolean" json:"inapp"`
		Email int `validate:"omitempty,mustbenumericalboolean" json:"email"`
	} `validate:"omitempty,dive" json:"notify" gorm:"embedded;embeddedPrefix:notify_"`

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
	} `validate:"omitempty,dive" gorm:"embedded;embeddedPrefix:show_"`

	Country string `validate:"omitempty" json:"country"`

	Language language.Language `validate:"omitempty,mustbesupportedlanguage" json:"language"`
}
