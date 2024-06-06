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

func (f FollowUseCase) GetFollowers(userID uint) ([]*dto.Follower, error) {
	following, err := f.fr.GetFollowers(userID)
	return following, err
}

func (f FollowUseCase) GetFollowing(userID uint) ([]*dto.Follower, error) {
	followers, err := f.fr.GetFollowing(userID)
	return followers, err
}

func (f FollowUseCase) Exists(followingID, followerID uint) (bool, error) {
	return f.fr.Exists(followingID, followerID)
}
