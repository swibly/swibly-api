package repository

import (
	"strings"

	"github.com/devkcud/arkhon-foundation/arkhon-api/internal/model"
	"github.com/devkcud/arkhon-foundation/arkhon-api/pkg/db"
	"gorm.io/gorm"
)

type UserRepository interface {
	Store(newUser model.User) error
	Update(id uint, newUser model.User) error
	Find(query model.UserQuery) (*model.User, error)
}

type userRepository struct {
	db *gorm.DB
}

var UserRepositoryInstance = NewRepository()

func NewRepository() UserRepository {
	return userRepository{db: db.Postgres.Model(&model.User{})}
}

func (u userRepository) Store(newUser model.User) error {
	if err := u.db.Create(&newUser).Error; err != nil {
		return err
	}
	return nil
}

func (u userRepository) Update(id uint, newUser model.User) error {
	if err := u.db.Where("id = ?", id).Updates(newUser).Error; err != nil {
		return err
	}
	return nil
}

func (u userRepository) Find(query model.UserQuery) (*model.User, error) {
	var fields []string
	var args []any

	if strings.TrimSpace(query.ID) != "" {
		fields = append(fields, "id = ?")
		args = append(args, query.ID)
	}

	if strings.TrimSpace(query.Username) != "" {
		fields = append(fields, "username = ?")
		args = append(args, query.Username)
	}

	if strings.TrimSpace(query.Email) != "" {
		fields = append(fields, "email = ?")
		args = append(args, query.Email)
	}

	if len(fields) == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	var user model.User

	if err := u.db.Where(strings.Join(fields, " OR "), args...).First(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}
