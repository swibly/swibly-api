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
	return puc.pr.Create(p)
}

func (puc *ProjectUseCase) GetPublicAll(page, perPage int) (*dto.Pagination[dto.ProjectInformation], error) {
	return puc.pr.Get(&model.Project{Published: true}, page, perPage)
}

func (puc *ProjectUseCase) GetByOwner(owner string, amIOwner bool, page, perPage int) (*dto.Pagination[dto.ProjectInformation], error) {
	var search *model.Project

	search.Owner = owner

	if !amIOwner {
		search.Published = true
	}

	return puc.pr.Get(search, page, perPage)
}

func (puc *ProjectUseCase) GetByID(id uint) (*dto.ProjectInformation, error) {
	project, err := puc.pr.Get(&model.Project{ID: id}, 1, 1)

	if err != nil {
		return nil, err
	}

	return project.Data[0], nil
}

func (puc *ProjectUseCase) GetContent(id uint) (any, error) {
	return puc.pr.GetContent(id)
}

func (puc *ProjectUseCase) SearchByName(name string, page, perpage int) (*dto.Pagination[dto.ProjectInformation], error) {
	return puc.pr.SearchByName(name, page, perpage)
}

func (puc *ProjectUseCase) SaveContent(id uint, content any) error {
	return puc.pr.Update(id, &model.Project{Content: content})
}

func (puc *ProjectUseCase) Publish(id uint) error {
	return puc.pr.Update(id, &model.Project{Published: true})
}

func (puc *ProjectUseCase) Unpublish(id uint) error {
	return puc.pr.Update(id, &model.Project{Published: false})
}

func (puc *ProjectUseCase) Favorite(userId, projectId uint) error {
	return puc.pr.Unfavorite(userId, projectId)
}

func (puc *ProjectUseCase) Unfavorite(userId, projectId uint) error {
	return puc.pr.Unfavorite(userId, projectId)
}

func (puc *ProjectUseCase) Delete(id uint) error {
	return puc.pr.Delete(id)
}
