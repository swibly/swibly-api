package usecase

import (
	"github.com/devkcud/arkhon-foundation/arkhon-api/internal/model"
	"github.com/devkcud/arkhon-foundation/arkhon-api/internal/service/repository"
)

type APIKeyUseCase struct {
	ar repository.APIKeyRepository
}

func NewAPIKeyUseCase() APIKeyUseCase {
	return APIKeyUseCase{ar: repository.NewAPIKeyRepository()}
}

func (auc *APIKeyUseCase) Find(key string) (*model.APIKey, error) {
	return auc.ar.Find(key)
}
