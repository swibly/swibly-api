package usecase

import (
	"github.com/devkcud/arkhon-foundation/arkhon-api/internal/model/dto"
	"github.com/devkcud/arkhon-foundation/arkhon-api/internal/service/repository"
)

type FollowUseCase struct {
	fr repository.FollowRepository
}

func NewFollowUseCase() FollowUseCase {
	return FollowUseCase{fr: repository.NewFollowRepository()}
}

func (f FollowUseCase) FollowUser(followingID, followerID uint) error {
	if err := f.fr.Follow(followingID, followerID); err != nil {
		return err
	}

	return nil
}

func (f FollowUseCase) UnfollowUser(followingID, followerID uint) error {
	if err := f.fr.Unfollow(followingID, followerID); err != nil {
		return err
	}

	return nil
}

func (f FollowUseCase) GetFollowers(userID uint, page, perpage int) (*dto.Pagination[dto.Follower], error) {
	return f.fr.GetFollowers(userID, page, perpage)
}

func (f FollowUseCase) GetFollowing(userID uint, page, perpage int) (*dto.Pagination[dto.Follower], error) {
	return f.fr.GetFollowing(userID, page, perpage)
}

func (f FollowUseCase) GetFollowersCount(userID uint, page int, perpage int) (*dto.Pagination[dto.Follower], error) {
	return f.fr.GetFollowers(userID, page, perpage)
}

func (f FollowUseCase) GetFollowingCount(userID uint) (int64, error) {
	return f.fr.GetFollowingCount(userID)
}

func (f FollowUseCase) Exists(followingID, followerID uint) (bool, error) {
	return f.fr.Exists(followingID, followerID)
}
