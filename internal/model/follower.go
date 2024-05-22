package model

type Follower struct {
	FollowingID uint `gorm:"index"`
	FollowerID  uint `gorm:"index"`

	Since int64 `gorm:"autoCreateTime"`
}
