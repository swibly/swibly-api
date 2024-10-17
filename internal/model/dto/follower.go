package dto

import "time"

type NewFollower struct {
	FollowerID  uint `json:"follower_id"`
	FollowingID uint `json:"following_id"`
}

type Follower struct {
	ID             uint      `json:"id"`
	FirstName      string    `json:"firstname"`
	LastName       string    `json:"lastname"`
	ProfilePicture string    `json:"pfp"`
	Verified       bool      `json:"verified"`
	Username       string    `json:"username"`
	Since          time.Time `json:"following_since"`
	Show           struct {
		Profile   bool `json:"profile"`
		Image     bool `json:"image"`
		Followers bool `json:"followers"`
		Following bool `json:"following"`
	} `gorm:"embedded;embeddedPrefix:show_" json:"show"`
}
