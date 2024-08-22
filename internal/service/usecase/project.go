package usecase

import (
	"github.com/devkcud/arkhon-foundation/arkhon-api/internal/model"
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

func (puc *ProjectUseCase) GetPublicAll(page, perPage int) (*dto.Pagination[model.Project], error) {
	return puc.pr.GetPublicAll(page, perPage)
}
