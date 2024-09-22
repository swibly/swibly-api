package repository

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/devkcud/arkhon-foundation/arkhon-api/internal/model"
	"github.com/devkcud/arkhon-foundation/arkhon-api/internal/model/dto"
	"github.com/devkcud/arkhon-foundation/arkhon-api/pkg/db"
	"github.com/devkcud/arkhon-foundation/arkhon-api/pkg/pagination"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type userRepository struct {
	db *gorm.DB
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
	return userRepository{db: db.Postgres}
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

	hasher := sha256.Sum256([]byte(user.Email))

	user.Permissions = []string{}
	user.ProfilePicture = fmt.Sprintf("https://www.gravatar.com/avatar/%s?s=512&d=monsterid&r=g", hex.EncodeToString(hasher[:]))

	// Temp repos
	// Tricky, not recommended, not performant
	// But I don't care ;)
	tempFollowRepo := NewFollowRepository()
	tempPermissionRepo := NewPermissionRepository()

	if count, err := tempFollowRepo.GetFollowersCount(user.ID); err != nil {
		return nil, err
	} else {
		user.Followers = count
	}

	if count, err := tempFollowRepo.GetFollowingCount(user.ID); err != nil {
		return nil, err
	} else {
		user.Following = count
	}

	if permissions, err := tempPermissionRepo.GetByUser(user.ID); err != nil {
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
	return u.db.Where("id = ?", id).Unscoped().Delete(&model.User{}).Error
}
