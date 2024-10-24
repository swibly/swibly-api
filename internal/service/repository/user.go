package repository

import (
	"fmt"
	"strings"

	"github.com/swibly/swibly-api/internal/model"
	"github.com/swibly/swibly-api/internal/model/dto"
	"github.com/swibly/swibly-api/pkg/aws"
	"github.com/swibly/swibly-api/pkg/db"
	"github.com/swibly/swibly-api/pkg/pagination"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type userRepository struct {
	db *gorm.DB

	followRepo     FollowRepository
	permissionRepo PermissionRepository
}

type UserRepository interface {
	Create(*model.User) error

	Update(uint, *dto.UserUpdate) error

	UnsafeGet(*model.User) (*model.User, error)
	Get(*model.User) (*dto.UserProfile, error)

	SearchByName(name string, page, perpage int) (*dto.Pagination[dto.UserProfile], error)

	Delete(uint) error
}

func NewUserRepository() UserRepository {
	return &userRepository{db: db.Postgres, followRepo: NewFollowRepository(), permissionRepo: NewPermissionRepository()}
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

	return user, nil
}

func (u userRepository) SearchByName(name string, page, perPage int) (*dto.Pagination[dto.UserProfile], error) {
	terms := strings.Fields(name)

	var query = u.db.Model(&model.User{})

	for _, term := range terms {
		alike := fmt.Sprintf("%%%s%%", strings.ToLower(term))
		query = query.Or("LOWER(username) LIKE ? OR LOWER(first_name) LIKE ? OR LOWER(last_name) LIKE ?", alike, alike, alike)
	}

	query = query.Order(clause.OrderBy{
		Expression: clause.Expr{
			SQL:                "CASE WHEN LOWER(username) = LOWER(?) THEN 1 ELSE 2 END",
			Vars:               []any{name},
			WithoutParentheses: true,
		},
	})

	return pagination.Generate[dto.UserProfile](query, page, perPage)
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
		if err := tx.Where("id = ?", projectID).Unscoped().Delete(&model.Project{}).Error; err != nil {
			tx.Rollback()
			return err
		}
		if err := aws.DeleteProjectImage(projectID); err != nil {
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
