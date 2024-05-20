package model

import "time"

type Follower struct {
	FollowingID uint `gorm:"index"`
	FollowerID  uint `gorm:"index"`

	Since time.Time `gorm:"autoCreateTime"`
}
