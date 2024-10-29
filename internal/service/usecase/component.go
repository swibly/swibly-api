package usecase

import (
	"github.com/swibly/swibly-api/internal/model"
	"github.com/swibly/swibly-api/internal/model/dto"
	"github.com/swibly/swibly-api/internal/service/repository"
	"github.com/swibly/swibly-api/pkg/utils"
)

type ComponentUseCase struct {
	cr repository.ComponentRepository
}

func NewComponentUseCase() ComponentUseCase {
	return ComponentUseCase{cr: repository.NewComponentRepository(repository.NewUserRepository())}
}

func (cuc *ComponentUseCase) Create(createModel *dto.ComponentCreation) error {
	return cuc.cr.Create(createModel)
}

func (cuc *ComponentUseCase) Update(componentID uint, updateModel *dto.ComponentUpdate) error {
	return cuc.cr.Update(componentID, updateModel)
}

func (cuc *ComponentUseCase) Publish(componentID uint) error {
	return cuc.cr.Update(componentID, &dto.ComponentUpdate{Public: utils.ToPtr(true)})
}

func (cuc *ComponentUseCase) Unpublish(componentID uint) error {
	return cuc.cr.Update(componentID, &dto.ComponentUpdate{Public: utils.ToPtr(false)})
}

func (cuc *ComponentUseCase) GetByID(issuerID, componentID uint) (*dto.ComponentInfo, error) {
	return cuc.cr.Get(issuerID, &model.Component{ID: componentID})
}

func (cuc *ComponentUseCase) GetPublic(issuerID uint, page, perPage int, freeOnly bool) (*dto.Pagination[dto.ComponentInfo], error) {
	return cuc.cr.GetPublic(issuerID, page, perPage, freeOnly)
}

func (cuc *ComponentUseCase) GetTrashed(ownerID uint, page, perPage int) (*dto.Pagination[dto.ComponentInfo], error) {
	return cuc.cr.GetTrashed(ownerID, page, perPage)
}

func (cuc *ComponentUseCase) GetByOwnerID(issuerID, ownerID uint, onlyPublic bool, page, perPage int) (*dto.Pagination[dto.ComponentInfo], error) {
	return cuc.cr.GetByOwnerID(issuerID, ownerID, onlyPublic, page, perPage)
}

func (cuc *ComponentUseCase) GetOwned(issuerID, userID uint, onlyPublic bool, page, perPage int) (*dto.Pagination[dto.ComponentInfo], error) {
	return cuc.cr.GetOwned(issuerID, userID, onlyPublic, page, perPage)
}

func (cuc *ComponentUseCase) GetHoldersByID(issuerID, componentID uint, page, perPage int) (*dto.Pagination[dto.ComponentInfo], error) {
	return cuc.cr.GetHoldersByID(issuerID, componentID, page, perPage)
}

func (cuc *ComponentUseCase) Search(issuerID uint, search *dto.SearchComponent, page, perPage int) (*dto.Pagination[dto.ComponentInfo], error) {
	return cuc.cr.Search(issuerID, search, page, perPage)
}

func (cuc *ComponentUseCase) Buy(issuerID, componentID uint) error {
	return cuc.cr.Buy(issuerID, componentID)
}

func (cuc *ComponentUseCase) Sell(issuerID, componentID uint) error {
	return cuc.cr.Sell(issuerID, componentID)
}

func (cuc *ComponentUseCase) SafeDelete(componentID uint) error {
	return cuc.cr.SafeDelete(componentID)
}

func (cuc *ComponentUseCase) Restore(componentID uint) error {
	return cuc.cr.Restore(componentID)
}

func (cuc *ComponentUseCase) UnsafeDelete(componentID uint) error {
	return cuc.cr.UnsafeDelete(componentID)
}

func (cuc *ComponentUseCase) ClearTrash(userID uint) error {
	return cuc.cr.ClearTrash(userID)
}
