package repository

import (
	"fmt"

	"github.com/devkcud/arkhon-foundation/arkhon-api/internal/model"
	"github.com/devkcud/arkhon-foundation/arkhon-api/internal/model/dto"
	"github.com/devkcud/arkhon-foundation/arkhon-api/pkg/db"
	"gorm.io/gorm"
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

	return user, nil
}

func (u userRepository) SearchLikeName(username string) ([]*dto.ProfileSearch, error) {
	var users []*dto.ProfileSearch
	alike := fmt.Sprintf("%%%s%%", username)
	err := u.db.
		Model(&model.User{}).
		Where("(username LIKE ? OR first_name LIKE ? OR last_name LIKE ?) AND show_profile <> -1", alike, alike, alike).
		Find(&users).Error

	if err != nil {
		return nil, err
	}

	return users, nil
}

func (u userRepository) Delete(id uint) error {
	return u.db.Where("id = ?", id).Unscoped().Delete(&model.User{}).Error
}
