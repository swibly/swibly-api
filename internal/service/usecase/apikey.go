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

func (auc *APIKeyUseCase) Create(ownerID, maxUsage uint) (*model.APIKey, error) {
	key := new(model.APIKey)
	key.Key = uuid.New().String()

	key.OwnerID = ownerID
	key.MaxUsage = maxUsage

	return key, auc.ar.Store(key)
}

func (auc *APIKeyUseCase) Update(key string, updateModel *dto.APIKey) error {
	return auc.ar.Update(key, updateModel)
}

func (auc *APIKeyUseCase) RegisterUse(key string) error {
	return auc.ar.RegisterUse(key)
}

func (auc *APIKeyUseCase) FindAll() ([]*model.APIKey, error) {
	return auc.ar.FindAll()
}

func (auc *APIKeyUseCase) Find(key string) (*model.APIKey, error) {
	return auc.ar.Find(key)
}

func (auc *APIKeyUseCase) FindByOwnerID(id uint) ([]*model.APIKey, error) {
	return auc.ar.FindByOwnerID(id)
}

func (auc *APIKeyUseCase) Delete(key string) error {
	return auc.ar.Delete(key)
}
