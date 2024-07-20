package usecase

import (
	"github.com/devkcud/arkhon-foundation/arkhon-api/internal/model"
	"github.com/devkcud/arkhon-foundation/arkhon-api/internal/model/dto"
	"github.com/devkcud/arkhon-foundation/arkhon-api/internal/service/repository"

	"github.com/google/uuid"
)

type APIKeyUseCase struct {
	ar repository.APIKeyRepository
}

func NewAPIKeyUseCase() APIKeyUseCase {
	return APIKeyUseCase{ar: repository.NewAPIKeyRepository()}
}

func (auc *APIKeyUseCase) Create(ownerUsername string, maxUsage uint) (*model.APIKey, error) {
	key := new(model.APIKey)
	key.Key = uuid.New().String()

	if ownerUsername != "" {
		if _, err := NewUserUseCase().GetByUsername(ownerUsername); err != nil {
			return nil, err
		}

		key.OwnerUsername = ownerUsername
	}

	key.MaxUsage = maxUsage

	return key, auc.ar.Store(key)
}

func (auc *APIKeyUseCase) Update(key string, updateModel *dto.UpdateAPIKey) error {
	return auc.ar.Update(key, updateModel)
}

func (auc *APIKeyUseCase) RegisterUse(key string) error {
	return auc.ar.RegisterUse(key)
}

func (auc *APIKeyUseCase) FindAll(page, perPage int) (*dto.Pagination[dto.ReadAPIKey], error) {
	return auc.ar.FindAll(page, perPage)
}

func (auc *APIKeyUseCase) Find(key string) (*dto.ReadAPIKey, error) {
	return auc.ar.Find(key)
}

func (auc *APIKeyUseCase) FindByOwnerUsername(username string, page, perPage int) (*dto.Pagination[dto.ReadAPIKey], error) {
	return auc.ar.FindByOwnerUsername(username, page, perPage)
}

func (auc *APIKeyUseCase) Delete(key string) error {
	return auc.ar.Delete(key)
}
