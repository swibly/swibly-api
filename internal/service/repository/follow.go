package repository

import (
	"github.com/devkcud/arkhon-foundation/arkhon-api/internal/model"
	"github.com/devkcud/arkhon-foundation/arkhon-api/internal/model/dto"
	"github.com/devkcud/arkhon-foundation/arkhon-api/pkg/db"
	"github.com/devkcud/arkhon-foundation/arkhon-api/pkg/pagination"
	"gorm.io/gorm"
)

type followRepository struct {
	db *gorm.DB
}

type FollowRepository interface {
	Exists(followingID, followerID uint) (bool, error)

	Follow(followingID, followerID uint) error
	Unfollow(followingID, followerID uint) error

	GetFollowers(userID uint, page, perPage int) (*dto.Pagination[dto.Follower], error)
	GetFollowing(userID uint, page, perPage int) (*dto.Pagination[dto.Follower], error)
	GetFollowersCount(userID uint) (int64, error)
	GetFollowingCount(userID uint) (int64, error)
}

func NewFollowRepository() FollowRepository {
	return &followRepository{db: db.Postgres}
}

func (f *followRepository) Follow(followingID, followerID uint) error {
	return f.db.Create(&model.Follower{FollowingID: followingID, FollowerID: followerID}).Error
}

func (f *followRepository) Unfollow(followingID, followerID uint) error {
	return f.db.Where("following_id = ? AND follower_id = ?", followingID, followerID).Delete(&model.Follower{}).Error
}

func (f *followRepository) Exists(followingID, followerID uint) (bool, error) {
	var count int64
	err := f.db.Model(&model.Follower{}).Where("following_id = ? AND follower_id = ?", followingID, followerID).Count(&count).Error
	return count > 0, err
}

func (f *followRepository) GetFollowers(userID uint, page, perPage int) (*dto.Pagination[dto.Follower], error) {
	query := f.db.Table("users").
		Select("users.*, followers.since").
		Joins("JOIN followers ON followers.follower_id = users.id").
		Where("followers.following_id = ?", userID)

	return pagination.Generate[dto.Follower](query, page, perPage)
}

func (f *followRepository) GetFollowing(userID uint, page, perPage int) (*dto.Pagination[dto.Follower], error) {
	query := f.db.Table("users").
		Select("users.*, followers.since").
		Joins("JOIN followers ON followers.following_id = users.id").
		Where("followers.follower_id = ?", userID)

	return pagination.Generate[dto.Follower](query, page, perPage)
}

func (f *followRepository) GetFollowersCount(userID uint) (int64, error) {
	var totalRecords int64

	err := f.db.Model(&model.Follower{}).
		Where("following_id = ?", userID).
		Count(&totalRecords).Error

	return totalRecords, err
}

func (f *followRepository) GetFollowingCount(userID uint) (int64, error) {
	var totalRecords int64

	err := f.db.Model(&model.Follower{}).
		Where("follower_id = ?", userID).
		Count(&totalRecords).Error

	return totalRecords, err
}
