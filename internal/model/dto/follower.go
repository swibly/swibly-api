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
}
