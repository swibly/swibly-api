package dto

import "time"

type NewFollower struct {
	FollowerID  uint
	FollowingID uint
}

type Follower struct {
	ID             uint
	FirstName      string
	LastName       string
	Verified       bool
	Username       string
	FollowingSince time.Time
}
