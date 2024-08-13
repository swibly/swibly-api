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
	Update(string, *dto.UpdateAPIKey) error
	RegisterUse(string) error
	FindAll(page, perPage int) (*dto.Pagination[dto.ReadAPIKey], error)
	Find(key string) (*dto.ReadAPIKey, error)
	FindByOwnerUsername(username string, page, perPage int) (*dto.Pagination[dto.ReadAPIKey], error)
	Delete(string) error
	Regenerate(oldKey, newKey string) (*model.APIKey, error)
}

func NewAPIKeyRepository() APIKeyRepository {
	return apiKeyRepository{db: db.Postgres}
}

func (a apiKeyRepository) Store(createModel *model.APIKey) error {
	return a.db.Create(&createModel).Error
}

func (a apiKeyRepository) Update(key string, updateModel *dto.UpdateAPIKey) error {
	return a.db.Model(&model.APIKey{}).Where("key = ?", key).Updates(&updateModel).Error
}

func (a apiKeyRepository) RegisterUse(key string) error {
	return a.db.Exec("UPDATE api_keys SET times_used = times_used + 1 WHERE key = ?", key).Error
}

func (a apiKeyRepository) FindAll(page, perPage int) (*dto.Pagination[dto.ReadAPIKey], error) {
	return pagination.Generate[dto.ReadAPIKey](a.db.Model(&model.APIKey{}).Exec("SELECT * FROM api_keys"), page, perPage)
}

func (a apiKeyRepository) Find(key string) (*dto.ReadAPIKey, error) {
	var apikey *dto.ReadAPIKey

	if err := a.db.Model(&model.APIKey{}).First(&apikey, "key = ?", key).Error; err != nil {
		return nil, err
	}

	return apikey, nil
}

func (a apiKeyRepository) FindByOwnerUsername(ownerUsername string, page, perPage int) (*dto.Pagination[dto.ReadAPIKey], error) {
	return pagination.Generate[dto.ReadAPIKey](a.db.Model(&model.APIKey{}).Where("owner_username = ?", ownerUsername), page, perPage)
}

func (a apiKeyRepository) Delete(key string) error {
	return a.db.Exec("DELETE FROM api_keys WHERE key = ?", key).Error
}

func (a apiKeyRepository) Regenerate(oldKey, newKey string) (*model.APIKey, error) {
	var apiKey model.APIKey

	err := a.db.Where("key = ?", oldKey).First(&apiKey).Error
	if err != nil {
		return nil, err
	}

	err = a.db.Model(&apiKey).Update("key", newKey).Error
	if err != nil {
		return nil, err
	}

	return &apiKey, nil
}
