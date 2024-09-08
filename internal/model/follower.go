package model

import (
	"time"
)

type Follower struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time

	FollowingID uint `gorm:"index"`
	FollowerID  uint `gorm:"index"`

	Since time.Time `gorm:"autoCreateTime"`
}
