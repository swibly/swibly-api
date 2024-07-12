package repository

import (
	"github.com/devkcud/arkhon-foundation/arkhon-api/internal/model"
	"github.com/devkcud/arkhon-foundation/arkhon-api/internal/model/dto"
	"github.com/devkcud/arkhon-foundation/arkhon-api/pkg/db"
	"github.com/devkcud/arkhon-foundation/arkhon-api/pkg/pagination"
	"gorm.io/gorm"
)

type apiKeyRepository struct {
	db *gorm.DB
}

type APIKeyRepository interface {
	Store(*model.APIKey) error
	Update(string, *dto.APIKey) error
	RegisterUse(string) error
	FindAll(page, perPage int) (*dto.Pagination[model.APIKey], error)
	Find(key string) (*model.APIKey, error)
	FindByOwnerID(id uint, page, perPage int) (*dto.Pagination[model.APIKey], error)
	Delete(string) error
}

func NewAPIKeyRepository() APIKeyRepository {
	return apiKeyRepository{db: db.Postgres}
}

func (a apiKeyRepository) Store(createModel *model.APIKey) error {
	return a.db.Create(&createModel).Error
}

func (a apiKeyRepository) Update(key string, updateModel *dto.APIKey) error {
	return a.db.Model(&model.APIKey{}).Where("key = ?", key).Updates(&updateModel).Error
}

func (a apiKeyRepository) RegisterUse(key string) error {
	return a.db.Exec("UPDATE api_keys SET times_used = times_used + 1 WHERE key = ?", key).Error
}

func (a apiKeyRepository) FindAll(page, perPage int) (*dto.Pagination[model.APIKey], error) {
	return pagination.Generate[model.APIKey](a.db.Model(&model.APIKey{}).Exec("SELECT * FROM api_keys"), page, perPage)
}

func (a apiKeyRepository) Find(key string) (*model.APIKey, error) {
	var apikey *model.APIKey

	if err := a.db.First(&apikey, "key = ?", key).Error; err != nil {
		return nil, err
	}

	return apikey, nil
}

func (a apiKeyRepository) FindByOwnerID(ownerID uint, page, perPage int) (*dto.Pagination[model.APIKey], error) {
	return pagination.Generate[model.APIKey](a.db.Model(&model.APIKey{}).Where("owner_id = ?", ownerID), page, perPage)
}

func (a apiKeyRepository) Delete(key string) error {
	return a.db.Exec("DELETE FROM api_keys WHERE key = ?", key).Error
}
