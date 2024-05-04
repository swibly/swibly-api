package repository

import (
	"fmt"
	"reflect"
	"strings"
	"sync"

	"github.com/devkcud/arkhon-foundation/arkhon-api/internal/model"
	"github.com/devkcud/arkhon-foundation/arkhon-api/pkg/db"
	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

// NOTE: As we modify the Repo default methods, we need to create a new repository (which imo is better for readability)
// Now, we should continue this pattern. Create a repository for every model :) (models that overlap, e.g. user and comments, can be merged together)
type UserRepository interface {
	Repository[model.User]
	GetComments(uint) ([]model.Comment, error)
	AddComment(uint, model.Comment) error
}

func NewUserRepository() UserRepository {
	return userRepository{db: db.Postgres}
}

func (u userRepository) Store(createModel *model.User) error {
	return u.db.Create(&createModel).Error
}

func (u userRepository) Update(id uint, updateModel *model.User) error {
	return u.db.Where("id = ?", id).Updates(&updateModel).Error
}

func (u userRepository) Find(searchModel *model.User) (*model.User, error) {
	var wg sync.WaitGroup
	var mu sync.Mutex

	fields := reflect.TypeOf(*searchModel)
	values := reflect.ValueOf(*searchModel)

	var conditions []string
	var queryValues []any

	for i := 0; i < fields.NumField(); i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()

			value := values.Field(i)

			if value.IsZero() {
				return
			}

			// FIXME: Hardcoded() "users" table name, not ideal for when the name change in the future
			fieldName := u.db.NamingStrategy.ColumnName("users", fields.Field(i).Name)

			mu.Lock()
			conditions = append(conditions, fmt.Sprintf("%s = ?", fieldName))
			queryValues = append(queryValues, value.Interface())
			mu.Unlock()
		}(i)
	}

	wg.Wait()

	var user model.User

	if err := u.db.Where(strings.Join(conditions, " OR "), queryValues...).First(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (u userRepository) Delete(id uint) error {
	return u.db.Where("id = ?", id).Unscoped().Delete(&model.User{}).Error
}

// COMMENTS

func (u userRepository) GetComments(userID uint) ([]model.Comment, error) {
	var user model.User

	if err := u.db.Preload("Comments").First(&user, userID).Error; err != nil {
		return nil, err
	}

	return user.Comments, nil
}

func (u userRepository) AddComment(userID uint, comment model.Comment) error {
	var user model.User

	if err := u.db.First(&user, userID).Error; err != nil {
		return err
	}

	comment.OwnerID = userID // Ensure the comment is associated with the correct user

	return u.db.Create(&comment).Error
}
