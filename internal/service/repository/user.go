package repository

import (
	"github.com/devkcud/arkhon-foundation/arkhon-api/internal/model"
	"github.com/devkcud/arkhon-foundation/arkhon-api/pkg/db"
	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository() Repository[model.User] {
	return userRepository{db: db.Postgres.Model(&model.User{})}
}

func (u userRepository) Store(createModel *model.User) error {
	return u.db.Create(&createModel).Error
}

func (u userRepository) Update(id uint, updateModel *model.User) error {
	return u.db.Where("id = ?", id).Updates(&updateModel).Error
}

func (u userRepository) Find(searchModel *model.User) (*model.User, error) {
	var user model.User

	if err := u.db.First(&user, &searchModel).Error; err != nil {
		return nil, err
	}

	return &user, nil
}
