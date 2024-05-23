package dto

import (
	"time"
)

type ProfileSearch struct {
	ID        uint `json:"id"`
	CreatedAt time.Time
	UpdatedAt time.Time

	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
	Username  string `json:"username"`

	Bio      string `json:"bio"`
	Verified bool   `json:"verified"`

	XP      uint64 `gorm:"default:500"`
	Arkhoin uint64 `gorm:"default:1000"`

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
}
