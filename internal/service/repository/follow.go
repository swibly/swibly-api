package repository

import (
	"math"

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
	Unfollow(followingID, followerID uint) error
	Exists(followingID, followerID uint) (bool, error)
	GetFollowers(userID uint, page int, pageSize int) (*dto.Pagination[dto.Follower], error)
	GetFollowing(userID uint, page int, pageSize int) (*dto.Pagination[dto.Follower], error)
}

func NewFollowRepository() FollowRepository {
	return followRepository{db: db.Postgres}
}

func (f followRepository) Follow(followingID, followerID uint) error {
	return f.db.Create(&model.Follower{FollowingID: followingID, FollowerID: followerID}).Error
}

func (f followRepository) Unfollow(followingID, followerID uint) error {
	return f.db.Where("following_id = ? AND follower_ID = ?", followingID, followerID).Delete(&model.Follower{}).Error
}

func (f followRepository) Exists(followingID, followerID uint) (bool, error) {
	var count int64
	err := f.db.Model(&model.Follower{}).Where("following_id = ? AND follower_id = ?", followingID, followerID).Count(&count).Error
	return count > 0, err
}

// TODO: Make sure show_profile is enabled

func (f followRepository) GetFollowers(userID uint, page int, pageSize int) (*dto.Pagination[dto.Follower], error) {
	var followers []*dto.Follower

	var totalRecords int64

	err := f.db.Model(&dto.Follower{}).
		Joins("JOIN users ON users.id = followers.follower_id").
		Where("users.id = ?", userID).
		Count(&totalRecords).Error
	if err != nil {
		return nil, err
	}

	totalPages := int(math.Ceil(float64(totalRecords) / float64(pageSize)))

	if page < 1 {
		page = 1
	} else if page > totalPages {
		page = totalPages
	}

	offset := (page - 1) * pageSize

	err = f.db.Table("users").
		Select("users.*, followers.since").
		Joins("JOIN followers ON followers.follower_id = users.id").
		Where("users.id = ?", userID).
		Offset(offset).
		Limit(pageSize).
		Scan(&followers).Error

	if err != nil {
		return nil, err
	}

	pagination := &dto.Pagination[dto.Follower]{
		Data:         followers,
		TotalRecords: int(totalRecords),
		TotalPages:   totalPages,
		CurrentPage:  page,
		NextPage:     page + 1,
		PreviousPage: page - 1,
	}

	if len(followers) == 0 {
		pagination.NextPage = -1
		pagination.PreviousPage = -1
	} else {
		pagination.NextPage = page + 1
		if pagination.NextPage > totalPages {
			pagination.NextPage = -1
		}

		pagination.PreviousPage = page - 1
		if pagination.PreviousPage < 1 {
			pagination.PreviousPage = -1
		}
	}

	return pagination, nil
}

func (f followRepository) GetFollowing(userID uint, page int, pageSize int) (*dto.Pagination[dto.Follower], error) {
	var followers []*dto.Follower

	var totalRecords int64

	err := f.db.Model(&dto.Follower{}).
		Joins("JOIN users ON users.id = followers.following_id").
		Where("users.id = ?", userID).
		Count(&totalRecords).Error
	if err != nil {
		return nil, err
	}

	totalPages := int(math.Ceil(float64(totalRecords) / float64(pageSize)))

	if page < 1 {
		page = 1
	} else if page > totalPages {
		page = totalPages
	}

	offset := (page - 1) * pageSize

	err = f.db.Table("users").
		Select("users.*, followers.since").
		Joins("JOIN followers ON followers.following_id = users.id").
		Where("users.id = ?", userID).
		Offset(offset).
		Limit(pageSize).
		Scan(&followers).Error

	if err != nil {
		return nil, err
	}

	pagination := &dto.Pagination[dto.Follower]{
		Data:         followers,
		TotalRecords: int(totalRecords),
		TotalPages:   totalPages,
		CurrentPage:  page,
		NextPage:     page + 1,
		PreviousPage: page - 1,
	}

	if len(followers) == 0 {
		pagination.NextPage = -1
		pagination.PreviousPage = -1
	} else {
		pagination.NextPage = page + 1
		if pagination.NextPage > totalPages {
			pagination.NextPage = -1
		}

		pagination.PreviousPage = page - 1
		if pagination.PreviousPage < 1 {
			pagination.PreviousPage = -1
		}
	}

	return pagination, nil
}
