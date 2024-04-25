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

func NewUserRepository() Repository[model.User] {
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

			// FIXME: Hardcoded "users" table name, not ideal for when the name change in the future
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
