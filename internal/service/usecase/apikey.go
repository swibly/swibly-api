package usecase

import (
	"github.com/devkcud/arkhon-foundation/arkhon-api/internal/model"
	"github.com/devkcud/arkhon-foundation/arkhon-api/internal/service/repository"

	"github.com/google/uuid"
)

type APIKeyUseCase struct {
	ar repository.APIKeyRepository
}

func NewAPIKeyUseCase() APIKeyUseCase {
	return APIKeyUseCase{ar: repository.NewAPIKeyRepository()}
}

func (auc *APIKeyUseCase) Create() (*model.APIKey, error) {
	key := new(model.APIKey)
	key.Key = uuid.New().String()

	return key, auc.ar.Store(key)
}

func (auc *APIKeyUseCase) Update(key string, updateModel *model.APIKey) error {
	return auc.ar.Update(key, updateModel)
}

func (auc *APIKeyUseCase) Find(key string) (*model.APIKey, error) {
	return auc.ar.Find(key)
}

func (auc *APIKeyUseCase) Delete(key string) error {
	return auc.ar.Delete(key)
}
