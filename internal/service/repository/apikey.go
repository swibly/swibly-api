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
	Create(createModel *model.APIKey) error
	Update(key string, updateModel *model.APIKey) error
	Delete(key string) error

	GetAll(page, perPage int) (*dto.Pagination[dto.ReadAPIKey], error)
	GetByKey(key string) (*dto.ReadAPIKey, error)
	GetByOwner(owner string, page, perPage int) (*dto.Pagination[dto.ReadAPIKey], error)

	RegisterUse(key string) error
	Regenerate(oldKey, newKey string) error
}

func NewAPIKeyRepository() APIKeyRepository {
	return &apiKeyRepository{db: db.Postgres}
}

func (akr *apiKeyRepository) Create(createModel *model.APIKey) error {
	return akr.db.Create(&createModel).Error
}

func (akr *apiKeyRepository) Update(key string, updateModel *model.APIKey) error {
	return akr.db.Model(&model.APIKey{}).Where("key = ?", key).Updates(updateModel).Error
}

func (akr *apiKeyRepository) Delete(key string) error {
	return akr.db.Unscoped().Delete(&model.APIKey{}, &model.APIKey{Key: key}).Error
}

func (akr *apiKeyRepository) GetAll(page, perPage int) (*dto.Pagination[dto.ReadAPIKey], error) {
	return pagination.Generate[dto.ReadAPIKey](akr.db.Model(&model.APIKey{}).Select("*"), page, perPage)
}

func (akr *apiKeyRepository) GetByKey(key string) (*dto.ReadAPIKey, error) {
	var apikey *dto.ReadAPIKey

	if err := akr.db.Model(&model.APIKey{}).First(&apikey, &model.APIKey{Key: key}).Error; err != nil {
		return nil, err
	}

	return apikey, nil
}

func (akr *apiKeyRepository) GetByOwner(owner string, page, perPage int) (*dto.Pagination[dto.ReadAPIKey], error) {
	return pagination.Generate[dto.ReadAPIKey](akr.db.Model(&model.APIKey{}).Where(&model.APIKey{Owner: owner}), page, perPage)
}

func (akr *apiKeyRepository) RegisterUse(key string) error {
	return akr.db.Model(&model.APIKey{}).Where(&model.APIKey{Key: key}).Update("times_used", gorm.Expr("times_used + ?", 1)).Error
}

func (akr *apiKeyRepository) Regenerate(oldKey, newKey string) error {
	if err := akr.db.Where(&model.APIKey{Key: oldKey}).First(&model.APIKey{}).Error; err != nil {
		return err
	}

	if err := akr.db.Model(&model.APIKey{}).Updates(&model.APIKey{Key: newKey}).Error; err != nil {
		return err
	}

	return nil
}
