package repository

import (
	"fmt"
	"strings"

	"github.com/devkcud/arkhon-foundation/arkhon-api/internal/model"
	"github.com/devkcud/arkhon-foundation/arkhon-api/internal/model/dto"
	"github.com/devkcud/arkhon-foundation/arkhon-api/pkg/db"
	"github.com/devkcud/arkhon-foundation/arkhon-api/pkg/pagination"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type projectRepository struct {
	db *gorm.DB
}

type ProjectRepository interface {
	Create(*dto.ProjectCreation) error
	Update(id uint, newModel *model.Project) error
	Delete(id uint) error

	Get(searchModel *model.Project, page, perPage int) (*dto.Pagination[dto.ProjectInformation], error)
	GetContent(id uint) (any, error)

	SearchByName(name string, page, perpage int) (*dto.Pagination[dto.ProjectInformation], error)

	Favorite(userId, projectId uint) error
	Unfavorite(userId, projectId uint) error
}

func NewProjectRepository() ProjectRepository {
	return projectRepository{db.Postgres}
}

func (p projectRepository) Create(project *dto.ProjectCreation) error {
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

func (p projectRepository) Update(id uint, newModel *model.Project) error {
	return p.db.Model(&model.Project{}).Where("id = ?", id).Updates(newModel).Error
}

func (p projectRepository) SearchByName(name string, page, perPage int) (*dto.Pagination[dto.ProjectInformation], error) {
	terms := strings.Fields(name)

	query := p.db.Model(&model.Project{}).Where("published = ?", true)

	orConditions := p.db
	for _, term := range terms {
		alike := fmt.Sprintf("%%%s%%", strings.ToLower(term))
		orConditions = orConditions.Or("LOWER(name) LIKE ?", alike)
	}
	query = query.Where(orConditions)

	query = query.Order(clause.OrderBy{
		Expression: clause.Expr{
			SQL:                "CASE WHEN LOWER(name) = LOWER(?) THEN 1 ELSE 2 END",
			Vars:               []any{name},
			WithoutParentheses: true,
		},
	})

	return pagination.Generate[dto.ProjectInformation](query, page, perPage)
}

func (p projectRepository) Get(searchModel *model.Project, page, perPage int) (*dto.Pagination[dto.ProjectInformation], error) {
	return pagination.Generate[dto.ProjectInformation](p.db.Model(&model.Project{}).Where(searchModel), page, perPage)
}

func (p projectRepository) GetContent(id uint) (any, error) {
	var project model.Project
	err := p.db.Where("id = ?", id).First(&project).Error

	return project.Content, err
}

func (p projectRepository) Favorite(userId, projectId uint) error {
	return p.db.Create(&model.ProjectFavorite{ProjectID: projectId, UserID: userId}).Error
}

func (p projectRepository) Unfavorite(userId, projectId uint) error {
	return p.db.Model(&model.ProjectFavorite{}).Where("user_id = ? AND project_id = ?", userId, projectId).Unscoped().Delete(&model.ProjectFavorite{}).Error
}

func (p projectRepository) Delete(id uint) error {
	// TODO: Send to trash instead of deleting
	return p.db.Where("id = ?", id).Unscoped().Delete(&model.Project{}).Error
}
