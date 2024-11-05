package repository

import (
	"github.com/swibly/swibly-api/internal/model"
	"github.com/swibly/swibly-api/internal/model/dto"
	"github.com/swibly/swibly-api/pkg/aws"
	"github.com/swibly/swibly-api/pkg/db"
	"github.com/swibly/swibly-api/pkg/pagination"
	"github.com/swibly/swibly-api/pkg/utils"
	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB

	followRepo       FollowRepository
	permissionRepo   PermissionRepository
	notificationRepo NotificationRepository
}

type UserRepository interface {
	Create(*model.User) error

	Update(uint, *dto.UserUpdate) error

	UnsafeGet(*model.User) (*model.User, error)
	Get(*model.User) (*dto.UserProfile, error)

	Search(issuerID uint, search *dto.SearchUser, page, perpage int) (*dto.Pagination[dto.UserProfile], error)

	Delete(uint) error
}

func NewUserRepository() UserRepository {
	return &userRepository{db: db.Postgres, followRepo: NewFollowRepository(), permissionRepo: NewPermissionRepository(), notificationRepo: NewNotificationRepository()}
}

func (u userRepository) Create(createModel *model.User) error {
	return u.db.Create(&createModel).Error
}

func (u userRepository) Update(id uint, updateModel *dto.UserUpdate) error {
	return u.db.Model(&model.User{}).Where("id = ?", id).Updates(&updateModel).Error
}

func (u userRepository) UnsafeGet(searchModel *model.User) (*model.User, error) {
	var user *model.User

	if err := u.db.First(&user, searchModel).Error; err != nil {
		return nil, err
	}

	return user, nil
}

func (u userRepository) Get(searchModel *model.User) (*dto.UserProfile, error) {
	var user *dto.UserProfile

	if err := u.db.Model(&model.User{}).First(&user, searchModel).Error; err != nil {
		return nil, err
	}

	user.Permissions = []string{}

	if count, err := u.followRepo.GetFollowersCount(user.ID); err != nil {
		return nil, err
	} else {
		user.Followers = count
	}

	if count, err := u.followRepo.GetFollowingCount(user.ID); err != nil {
		return nil, err
	} else {
		user.Following = count
	}

	if permissions, err := u.permissionRepo.GetByUser(user.ID); err != nil {
		return nil, err
	} else {
		for _, permission := range permissions {
			user.Permissions = append(user.Permissions, permission.Name)
		}
	}

	if count, err := u.notificationRepo.GetUnreadCount(user.ID); err != nil {
		return nil, err
	} else {
		user.UnreadNotifications = count
	}

	return user, nil
}

func (u userRepository) Search(issuerID uint, search *dto.SearchUser, page, perpage int) (*dto.Pagination[dto.UserProfile], error) {
	query := u.db.Model(&model.User{}).
		Where("show_profile = TRUE")

	if search.Name != nil {
		query = query.Where(
			"regexp_like(first_name, ?, 'i') OR regexp_like(last_name, ?, 'i') OR regexp_like(username, ?, 'i')",
			utils.RegexPrepareName(*search.Name),
			utils.RegexPrepareName(*search.Name),
			utils.RegexPrepareName(*search.Name),
		)

		// TODO: Create ranking system
	}

	if search.VerifiedOnly {
		query = query.Where("verified = ?", true)
	}

	if search.FollowedUsersOnly {
		query = query.Joins("JOIN followers uf ON uf.following_id = users.id").
			Where("uf.follower_id = ?", issuerID)
	}

	orderDirection := "DESC"
	if search.OrderAscending {
		orderDirection = "ASC"
	}

	if search.OrderAlphabetic {
		query = query.Order("first_name " + orderDirection + ", last_name " + orderDirection)
	} else if search.OrderCreationDate {
		query = query.Order("created_at " + orderDirection)
	} else if search.OrderModifiedDate {
		query = query.Order("updated_at " + orderDirection)
	} else if search.MostFollowers {
		query = query.Joins(`
			LEFT JOIN (
				SELECT following_id, COUNT(*) AS follower_count
				FROM followers
				GROUP BY following_id
			) follower_counts ON follower_counts.following_id = users.id`).
			Order("follower_count" + orderDirection + "NULLS LAST")
	} else {
		query = query.Order("created_at " + orderDirection)
	}

	return pagination.Generate[dto.UserProfile](query, page, perpage)
}

func (u userRepository) Delete(id uint) error {
	tx := u.db.Begin()

	user, err := u.Get(&model.User{ID: id})
	if err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Where("owner = ?", user.Username).Unscoped().Delete(&model.APIKey{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Where("user_id = ?", id).Unscoped().Delete(&model.UserPermission{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Where("project_id IN (?)", tx.Model(&model.ProjectOwner{}).Select("project_id").Where("user_id = ?", id)).Unscoped().Delete(&model.ProjectPublication{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	var projectIDs []uint
	if err := tx.Model(&model.ProjectOwner{}).Where("user_id = ?", id).Pluck("project_id", &projectIDs).Error; err != nil {
		tx.Rollback()
		return err
	}

	for _, projectID := range projectIDs {
		var project model.Project
		if err := tx.Where("id = ?", projectID).First(&project).Error; err != nil {
			tx.Rollback()
			return err
		}

		if project.BannerURL != "" {
			if err := aws.DeleteProjectImage(project.BannerURL); err != nil {
				tx.Rollback()
				return err
			}
		}

		if err := tx.Where("id = ?", projectID).Unscoped().Delete(&model.Project{}).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	if err := tx.Where("user_id = ?", id).Unscoped().Delete(&model.ProjectUserPermission{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Where("user_id = ?", id).Unscoped().Delete(&model.ProjectUserFavorite{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Model(&model.ComponentOwner{}).Where("user_id = ?", id).Update("user_id", nil).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Where("user_id = ?", id).Delete(&model.ComponentHolder{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Where("follower_id = ? OR following_id = ?", id, id).Unscoped().Delete(&model.Follower{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Where("user_id = ?", id).Unscoped().Delete(&model.PasswordResetKey{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Where("id = ?", id).Unscoped().Delete(&model.User{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}
