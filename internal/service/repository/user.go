package repository

import (
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
	Store(*model.User) error
	Update(uint, *dto.UserUpdate) error
	UnsafeFind(*model.User) (*model.User, error)
	Find(*model.User) (*dto.ProfileSearch, error)
	SearchLikeName(name string, page, perpage int) (*dto.Pagination[dto.ProfileSearch], error)
	Delete(uint) error
}

func NewUserRepository() UserRepository {
	return userRepository{db: db.Postgres}
}

func (u userRepository) Store(createModel *model.User) error {
	return u.db.Create(&createModel).Error
}

func (u userRepository) Update(id uint, updateModel *dto.UserUpdate) error {
	return u.db.Model(&model.User{}).Where("id = ?", id).Updates(&updateModel).Error
}

func (u userRepository) UnsafeFind(searchModel *model.User) (*model.User, error) {
	var user *model.User

	if err := u.db.First(&user, searchModel).Error; err != nil {
		return nil, err
	}

	return user, nil
}

func (u userRepository) Find(searchModel *model.User) (*dto.ProfileSearch, error) {
	var user *dto.ProfileSearch

	if err := u.db.Model(&model.User{}).First(&user, searchModel).Error; err != nil {
		return nil, err
	}

	// Temp follow repo
	// Tricky, not recommended, not performant
	// But I don't care
	tempFollowRepo := NewFollowRepository()

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

	return user, nil
}

func (u userRepository) SearchLikeName(name string, page, perPage int) (*dto.Pagination[dto.ProfileSearch], error) {
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

	return pagination.Generate[dto.ProfileSearch](query, page, perPage)
}

func (u userRepository) Delete(id uint) error {
	return u.db.Where("id = ?", id).Unscoped().Delete(&model.User{}).Error
}
