package repository

import (
	"github.com/devkcud/arkhon-foundation/arkhon-api/internal/model"
	"github.com/devkcud/arkhon-foundation/arkhon-api/internal/model/dto"
	"github.com/devkcud/arkhon-foundation/arkhon-api/pkg/db"
	"gorm.io/gorm"
)

type followRepository struct {
	db *gorm.DB
}

type FollowRepository interface {
	Follow(followingID, followerID uint) error
	GetFollowers(userID uint) ([]*dto.Follower, error)
	GetFollowing(userID uint) ([]*dto.Follower, error)
}

func NewFollowRepository() FollowRepository {
	return followRepository{db: db.Postgres}
}

func (f followRepository) Follow(followingID, followerID uint) error {
	if followingID == followerID {
		// Using this so it's easier to debug afterwards
		return gorm.ErrInvalidField
	}
	return f.db.Create(&model.Follower{FollowingID: followingID, FollowerID: followerID}).Error
}

// TODO: Make sure show_profile is enabled

func (f followRepository) GetFollowers(userID uint) ([]*dto.Follower, error) {
	var followers []*dto.Follower
	err := f.db.Table("users").
		Select("users.*, followers.since").
		Joins("JOIN followers ON followers.follower_id = users.id").
		Where("followers.following_id = ?", userID).
		Scan(&followers).Error
	return followers, err
}

func (f followRepository) GetFollowing(userID uint) ([]*dto.Follower, error) {
	var following []*dto.Follower
	err := f.db.Table("users").
		Select("users.*, followers.since").
		Joins("JOIN followers ON followers.following_id = users.id").
		Where("followers.follower_id = ?", userID).
		Scan(&following).Error
	return following, err
}
