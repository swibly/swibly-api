package repository

import (
	"fmt"
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
	SearchLikeName(string) ([]*dto.ProfileSearch, error)
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

func (u userRepository) SearchLikeName(name string) ([]*dto.ProfileSearch, error) {
	var users []*dto.ProfileSearch

	// Overcomplicated query to search users in all the params (separated by spaces)
	// wasted 4 hrs in this 💀

	var whereConditions []string
	var whereArgs []any

	for _, term := range strings.Fields(name) {
		alike := fmt.Sprintf("%%%s%%", strings.ToLower(term))
		whereConditions = append(whereConditions, "LOWER(username) LIKE ? OR LOWER(first_name) LIKE ? OR LOWER(last_name) LIKE ?")
		whereArgs = append(whereArgs, alike, alike, alike)
	}

	query := u.db.
		Model(&model.User{}).
		Where(strings.Join(whereConditions, " OR "), whereArgs...).
		Order(clause.OrderBy{
			Expression: clause.Expr{
				SQL:                "CASE WHEN LOWER(username) = LOWER(?) THEN 1 ELSE 2 END",
				Vars:               []any{name},
				WithoutParentheses: true,
			},
		}).
		Find(&users)

	if query.Error != nil {
		return nil, query.Error
	}

	return users, nil
}

func (u userRepository) Delete(id uint) error {
	return u.db.Where("id = ?", id).Unscoped().Delete(&model.User{}).Error
}
