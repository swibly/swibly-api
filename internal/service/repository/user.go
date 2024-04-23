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
	return userRepository{db: db.Postgres.Model(&model.User{})}
}

func (u userRepository) Store(createModel *model.User) error {
	return u.db.Create(&createModel).Error
}

func (u userRepository) Update(id uint, updateModel *model.User) error {
	return u.db.Where("id = ?", id).Updates(&updateModel).Error
}

func (u userRepository) Find(searchModel *model.User) (*model.User, error) {
	var wg sync.WaitGroup

	fields := reflect.TypeOf(*searchModel)
	values := reflect.ValueOf(*searchModel)

	var conditions []string
	var queryValues []any

	query := u.db.Model(&model.User{})

	for i := 0; i < fields.NumField(); i++ {
		go func(i int) {
			value := values.Field(i)

			if value.IsZero() {
				wg.Done()
				return
			}

			// FIXME: Hardcoded "users" table name, not ideal for when the name change in the future
			fieldName := u.db.NamingStrategy.ColumnName("users", fields.Field(i).Name)

			conditions = append(conditions, fmt.Sprintf("%s = ?", fieldName))
			queryValues = append(queryValues, value.Interface())

			wg.Done()
		}(i)
	}

	wg.Add(fields.NumField())

	wg.Wait()

	var user model.User

	if len(conditions) > 0 {
		query = query.Where(strings.Join(conditions, " OR "), queryValues...)
	}

	if err := query.First(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}
