package repository

import (
	"github.com/devkcud/arkhon-foundation/arkhon-api/internal/model"
	"github.com/devkcud/arkhon-foundation/arkhon-api/pkg/db"
	"gorm.io/gorm"
)

type apiKeyRepository struct {
	db *gorm.DB
}

type APIKeyRepository interface {
	Store(*model.APIKey) error
	Update(string, *model.APIKey) error
	Find(string) (*model.APIKey, error)
	Delete(string) error
}

func NewAPIKeyRepository() APIKeyRepository {
	return apiKeyRepository{db: db.Postgres}
}

func (a apiKeyRepository) Store(createModel *model.APIKey) error {
	return a.db.Create(&createModel).Error
}

func (a apiKeyRepository) Update(key string, updateModel *model.APIKey) error {
	return a.db.Where("key = ?", key).Updates(&updateModel).Error
}

func (a apiKeyRepository) Find(key string) (*model.APIKey, error) {
	var apikey *model.APIKey

	if err := a.db.First(&apikey, "key = ?", key).Error; err != nil {
		return nil, err
	}

	return apikey, nil
}

func (a apiKeyRepository) Delete(key string) error {
	return a.db.Exec("DELETE FROM api_keys WHERE key = ?", key).Error
}
