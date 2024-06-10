package repository

import (
	"fmt"
	"math"
	"strings"

	"github.com/devkcud/arkhon-foundation/arkhon-api/internal/model"
	"github.com/devkcud/arkhon-foundation/arkhon-api/internal/model/dto"
	"github.com/devkcud/arkhon-foundation/arkhon-api/pkg/db"
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

func (u userRepository) SearchLikeName(name string, page, perpage int) (*dto.Pagination[dto.ProfileSearch], error) {
	var users []*dto.ProfileSearch

	// Overcomplicated query to search users in all the params (separated by spaces)
	var whereConditions []string
	var whereArgs []any

	for _, term := range strings.Fields(name) {
		alike := fmt.Sprintf("%%%s%%", strings.ToLower(term))
		whereConditions = append(whereConditions, "LOWER(username) LIKE? OR LOWER(first_name) LIKE? OR LOWER(last_name) LIKE?")
		whereArgs = append(whereArgs, alike, alike, alike)
	}

	query := u.db.Model(&model.User{}).
		Where(strings.Join(whereConditions, " OR "), whereArgs...).
		Order(clause.OrderBy{
			Expression: clause.Expr{
				SQL:                "CASE WHEN LOWER(username) = LOWER(?) THEN 1 ELSE 2 END",
				Vars:               []any{name},
				WithoutParentheses: true,
			},
		})

	var totalRecords int64

	if err := query.Count(&totalRecords).Error; err != nil {
		return nil, err
	}

	totalPages := int(math.Ceil(float64(totalRecords) / float64(perpage)))

	if page < 1 {
		page = 1
	} else if page > totalPages {
		page = totalPages
	}

	offset := (page - 1) * perpage
	limit := perpage

	if err := query.Limit(limit).Offset(offset).Find(&users).Error; err != nil {
		return nil, err
	}

	pagination := &dto.Pagination[dto.ProfileSearch]{
		Data:         users,
		TotalRecords: int(totalRecords),
		CurrentPage:  page,
		TotalPages:   totalPages,
		NextPage:     page + 1,
		PreviousPage: page - 1,
	}

	if pagination.NextPage > totalPages {
		pagination.NextPage = -1
	}

	if pagination.PreviousPage < 1 {
		pagination.PreviousPage = -1
	}

	return pagination, nil
}

func (u userRepository) Delete(id uint) error {
	return u.db.Where("id = ?", id).Unscoped().Delete(&model.User{}).Error
}
