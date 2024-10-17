package usecase

import (
	"github.com/swibly/swibly-api/internal/model"
	"github.com/swibly/swibly-api/internal/model/dto"
	"github.com/swibly/swibly-api/internal/service/repository"
	"github.com/swibly/swibly-api/pkg/utils"
)

type ProjectUseCase struct {
	pr repository.ProjectRepository
}

func NewProjectUseCase() ProjectUseCase {
	return ProjectUseCase{pr: repository.NewProjectRepository(repository.NewUserRepository())}
}

func (puc ProjectUseCase) Create(createModel *dto.ProjectCreation) (uint, error) {
	return puc.pr.Create(createModel)
}

func (puc ProjectUseCase) Fork(projectID, issuerID uint) (uint, error) {
	content, err := puc.pr.GetContent(projectID)
	if err != nil {
		return 0, err
	}

	project, err := puc.pr.Get(issuerID, &model.Project{ID: projectID})
	if err != nil {
		return 0, err
	}

	return puc.pr.Create(&dto.ProjectCreation{
		Name:        project.Name,
		Description: project.Description,
		Budget:      project.Budget,
		OwnerID:     issuerID,
		Public:      false,
		Content:     content,
		Fork:        &project.ID,
	})
}

func (puc ProjectUseCase) Unlink(id uint) error {
	return puc.pr.Unlink(id)
}

func (puc ProjectUseCase) Update(projectID uint, updateModel *dto.ProjectUpdate) error {
	return puc.pr.Update(projectID, updateModel)
}

func (puc ProjectUseCase) Publish(projectID uint) error {
	return puc.pr.Update(projectID, &dto.ProjectUpdate{Published: utils.ToPtr(true)})
}

func (puc ProjectUseCase) Unpublish(projectID uint) error {
	return puc.pr.Update(projectID, &dto.ProjectUpdate{Published: utils.ToPtr(false)})
}

func (puc ProjectUseCase) Assign(userID uint, projectID uint, allowList *dto.ProjectAssign) error {
	return puc.pr.Assign(userID, projectID, allowList)
}

func (puc ProjectUseCase) GetByID(issuerID, id uint) (*dto.ProjectInfo, error) {
	return puc.pr.Get(issuerID, &model.Project{ID: id})
}

func (puc ProjectUseCase) GetByOwner(issuerID, userID uint, onlyPublic bool, page, perPage int) (*dto.Pagination[dto.ProjectInfo], error) {
	return puc.pr.GetByOwner(issuerID, userID, onlyPublic, page, perPage)
}

func (puc ProjectUseCase) GetPublic(issuerID uint, page, perPage int) (*dto.Pagination[dto.ProjectInfo], error) {
	return puc.pr.GetPublic(issuerID, page, perPage)
}

func (puc ProjectUseCase) GetFavorited(issuerID, userID uint, onlyPublic bool, page, perPage int) (*dto.Pagination[dto.ProjectInfo], error) {
	return puc.pr.GetFavorited(issuerID, userID, onlyPublic, page, perPage)
}

func (puc ProjectUseCase) GetTrashed(ownerID uint, page, perPage int) (*dto.Pagination[dto.ProjectInfo], error) {
	return puc.pr.GetTrashed(ownerID, page, perPage)
}

func (puc ProjectUseCase) SearchByName(issuerID uint, name string, page, perPage int) (*dto.Pagination[dto.ProjectInfo], error) {
	return puc.pr.SearchByName(issuerID, name, page, perPage)
}

func (puc ProjectUseCase) GetContent(projectID uint) (any, error) {
	return puc.pr.GetContent(projectID)
}

func (puc ProjectUseCase) SaveContent(projectID uint, content any) error {
	return puc.pr.SaveContent(projectID, content)
}

func (puc ProjectUseCase) Favorite(projectID, userID uint) error {
	return puc.pr.Favorite(projectID, userID)
}

func (puc ProjectUseCase) Unfavorite(projectID, userID uint) error {
	return puc.pr.Unfavorite(projectID, userID)
}

func (puc ProjectUseCase) ClearContent(projectID uint) error {
	return puc.pr.SaveContent(projectID, nil)
}

func (puc ProjectUseCase) Trash(id uint) error {
	return puc.pr.SafeDelete(id)
}

func (puc ProjectUseCase) Restore(id uint) error {
	return puc.pr.Restore(id)
}

func (puc ProjectUseCase) RemoveFromTrash(id uint) error {
	return puc.pr.UnsafeDelete(id)
}

func (puc ProjectUseCase) ClearTrash(issuerID uint) error {
	return puc.pr.ClearTrash(issuerID)
}
