package usecase

import (
	"github.com/swibly/swibly-api/internal/model"
	"github.com/swibly/swibly-api/internal/model/dto"
	"github.com/swibly/swibly-api/internal/service/repository"

	"github.com/google/uuid"
)

type APIKeyUseCase struct {
	akr repository.APIKeyRepository
}

func NewAPIKeyUseCase() APIKeyUseCase {
	return APIKeyUseCase{akr: repository.NewAPIKeyRepository()}
}

func (akuc *APIKeyUseCase) Create(owner string, maxUsage uint) (*dto.ReadAPIKey, error) {
	key := new(model.APIKey)
	key.Key = uuid.New().String()

	if owner != "" {
		if _, err := NewUserUseCase().GetByUsername(owner); err != nil {
			return nil, err
		}

		key.Owner = owner
	}

	key.MaxUsage = maxUsage

	return &dto.ReadAPIKey{
		Key:                key.Key,
		Owner:              key.Owner,
		EnabledKeyManage:   key.EnabledKeyManage,
		EnabledAuth:        key.EnabledAuth,
		EnabledSearch:      key.EnabledSearch,
		EnabledUserFetch:   key.EnabledUserFetch,
		EnabledUserActions: key.EnabledUserActions,
		TimesUsed:          key.TimesUsed,
		MaxUsage:           key.MaxUsage,
	}, akuc.akr.Create(key)
}

func (akuc *APIKeyUseCase) Update(key string, updateModel *dto.UpdateAPIKey) error {
	return akuc.akr.Update(key, &model.APIKey{
		Owner:              updateModel.Owner,
		EnabledKeyManage:   updateModel.EnabledKeyManage,
		EnabledAuth:        updateModel.EnabledAuth,
		EnabledSearch:      updateModel.EnabledSearch,
		EnabledUserFetch:   updateModel.EnabledUserFetch,
		EnabledUserActions: updateModel.EnabledUserActions,
		TimesUsed:          updateModel.TimesUsed,
		MaxUsage:           updateModel.MaxUsage,
	})
}

func (akuc *APIKeyUseCase) Delete(key string) error {
	return akuc.akr.Delete(key)
}

func (akuc *APIKeyUseCase) GetAll(page, perPage int) (*dto.Pagination[dto.ReadAPIKey], error) {
	return akuc.akr.GetAll(page, perPage)
}

func (akuc *APIKeyUseCase) GetByKey(key string) (*dto.ReadAPIKey, error) {
	return akuc.akr.GetByKey(key)
}

func (akuc *APIKeyUseCase) GetByOwner(owner string, page, perPage int) (*dto.Pagination[dto.ReadAPIKey], error) {
	return akuc.akr.GetByOwner(owner, page, perPage)
}

func (akuc *APIKeyUseCase) RegisterUse(key string) error {
	return akuc.akr.RegisterUse(key)
}

func (akuc *APIKeyUseCase) Regenerate(key string) error {
	newKey := uuid.New().String()

	if existingKey, _ := akuc.GetByKey(newKey); existingKey != nil {
		return akuc.Regenerate(key)
	}

	return akuc.akr.Update(key, &model.APIKey{Key: newKey})
}
