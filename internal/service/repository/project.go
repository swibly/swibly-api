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
	Store(*dto.ProjectCreation) error
	GetPublicAll(page, perPage int) (*dto.Pagination[dto.ProjectInformation], error)
	GetByOwnerUsername(ownerUsername string, amIOwner bool, page, perPage int) (*dto.Pagination[dto.ProjectInformation], error)
	GetByID(id uint) (*dto.ProjectInformation, error)
	GetContent(id uint) any
	SearchLikeName(name string, page, perpage int) (*dto.Pagination[dto.ProjectInformation], error)
	SaveContent(id uint, content any) error
	Publish(id uint) error
	Unpublish(id uint) error
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

func (p projectRepository) GetPublicAll(page, perPage int) (*dto.Pagination[dto.ProjectInformation], error) {
	return pagination.Generate[dto.ProjectInformation](p.db.Model(&model.Project{}).Select("*").Where("published = ?", true), page, perPage)
}

func (p projectRepository) GetByOwnerUsername(ownerUsername string, amIOwner bool, page, perPage int) (*dto.Pagination[dto.ProjectInformation], error) {
	var query *gorm.DB
	if amIOwner {
		query = p.db.Model(&model.Project{}).Select("*").Where("owner = ?", ownerUsername)
	} else {
		query = p.db.Model(&model.Project{}).Select("*").Where("owner = ?", ownerUsername).Where("published = ?", true)
	}

	return pagination.Generate[dto.ProjectInformation](query, page, perPage)
}

func (p projectRepository) GetByID(id uint) (*dto.ProjectInformation, error) {
	var project dto.ProjectInformation
	return &project, p.db.Model(&model.Project{}).First(&project, id).Error
}

func (p projectRepository) GetContent(id uint) any {
	var project model.Project
	p.db.Model(&model.Project{}).First(&project, id)
	return project.Content
}

func (p projectRepository) SearchLikeName(name string, page, perPage int) (*dto.Pagination[dto.ProjectInformation], error) {
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

func (p projectRepository) SaveContent(id uint, content any) error {
	return p.db.Model(&model.Project{}).Where("id = ?", id).Updates(&model.Project{Content: content}).Error
}

func (p projectRepository) Publish(id uint) error {
	return p.db.Model(&model.Project{}).Where("id = ?", id).Updates(map[string]any{
		"published": true,
	}).Error
}

func (p projectRepository) Unpublish(id uint) error {
	return p.db.Model(&model.Project{}).Where("id = ?", id).Updates(map[string]any{
		"published": false,
	}).Error
}
