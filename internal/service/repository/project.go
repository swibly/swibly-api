package repository

import (
	"github.com/devkcud/arkhon-foundation/arkhon-api/internal/model"
	"github.com/devkcud/arkhon-foundation/arkhon-api/internal/model/dto"
	"github.com/devkcud/arkhon-foundation/arkhon-api/pkg/db"
	"github.com/devkcud/arkhon-foundation/arkhon-api/pkg/pagination"
	"gorm.io/gorm"
)

type projectRepository struct {
	db *gorm.DB
}

type ProjectRepository interface {
	Store(*dto.ProjectCreation) error
	GetPublicAll(page, perPage int) (*dto.Pagination[model.Project], error)
	GetPublicOwner(ownerID string, page, perPage int) (*dto.Pagination[model.Project], error)
}

func NewProjectRepository() ProjectRepository {
	return projectRepository{db.Postgres}
}

func (p projectRepository) Store(project *dto.ProjectCreation) error {
	newProject := &model.Project{
		Owner:       project.Owner,
		Name:        project.Name,
		Description: project.Description,
		Content:     project.Content,
		Thumbnail:   project.Thumbnail,
		Budget:      project.Budget,
	}

	return p.db.Create(newProject).Error
}

func (p projectRepository) GetPublicAll(page, perPage int) (*dto.Pagination[model.Project], error) {
	return pagination.Generate[model.Project](p.db.Exec("SELECT * FROM projects WHERE published = true"), page, perPage)
}

func (p projectRepository) GetPublicOwner(owner string, page, perPage int) (*dto.Pagination[model.Project], error) {
	return pagination.Generate[model.Project](p.db.Exec("SELECT * FROM projects WHERE published = true AND owner = ?", owner), page, perPage)
}
