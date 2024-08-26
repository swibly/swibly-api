package usecase

import (
	"github.com/devkcud/arkhon-foundation/arkhon-api/internal/model/dto"
	"github.com/devkcud/arkhon-foundation/arkhon-api/internal/service/repository"
)

type ProjectUseCase struct {
	pr repository.ProjectRepository
}

func NewProjectUseCase() ProjectUseCase {
	return ProjectUseCase{pr: repository.NewProjectRepository()}
}

func (puc *ProjectUseCase) Create(p *dto.ProjectCreation) error {
	return puc.pr.Store(p)
}

func (puc *ProjectUseCase) GetPublicAll(page, perPage int) (*dto.Pagination[dto.ProjectInformation], error) {
	return puc.pr.GetPublicAll(page, perPage)
}

func (puc *ProjectUseCase) GetByOwnerUsername(ownerUsername string, amIOwner bool, page, perPage int) (*dto.Pagination[dto.ProjectInformation], error) {
	return puc.pr.GetByOwnerUsername(ownerUsername, amIOwner, page, perPage)
}

func (puc *ProjectUseCase) GetByID(id uint) (*dto.ProjectInformation, error) {
	return puc.pr.GetByID(id)
}

func (puc *ProjectUseCase) GetContent(id uint) any {
	return puc.pr.GetContent(id)
}

func (puc *ProjectUseCase) GetBySimilarName(name string, page, perpage int) (*dto.Pagination[dto.ProjectInformation], error) {
	return puc.pr.SearchLikeName(name, page, perpage)
}

func (puc *ProjectUseCase) SaveContent(id uint, content any) error {
	return puc.pr.SaveContent(id, content)
}

func (puc *ProjectUseCase) Publish(id uint) error {
	return puc.pr.Publish(id)
}

func (puc *ProjectUseCase) Unpublish(id uint) error {
	return puc.pr.Unpublish(id)
}

func (puc *ProjectUseCase) Delete(id uint) error {
	return puc.pr.Delete(id)
}
