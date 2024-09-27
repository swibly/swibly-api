package repository

import (
	"errors"

	"github.com/devkcud/arkhon-foundation/arkhon-api/internal/model"
	"github.com/devkcud/arkhon-foundation/arkhon-api/internal/model/dto"
	"github.com/devkcud/arkhon-foundation/arkhon-api/pkg/db"
	"gorm.io/gorm"
)

type projectRepository struct {
	db *gorm.DB

	userRepo UserRepository
}

type ProjectRepository interface {
	Create(createModel *dto.ProjectCreation) error

	Assign(userID uint, projectID uint, allowList *dto.Allow) error

	Get(*model.Project) (*dto.ProjectInfo, error)
	GetByOwner(userID uint, onlyPublic bool, page, pageSize int) (*dto.Pagination[dto.ProjectInfo], error)
	GetPublic(page, perPage int) (*dto.Pagination[dto.ProjectInfo], error)

	SearchByName(name string, page, perpage int) (*dto.Pagination[dto.ProjectInfo], error)

	GetContent() (any, error)
	SaveContent(any) error

	SafeDelete(uint) error
	UnsafeDelete(uint) error
}

func NewProjectRepository(userRepo UserRepository) ProjectRepository {
	return &projectRepository{db.Postgres, userRepo}
}

func (pr *projectRepository) Create(createModel *dto.ProjectCreation) error {
	tx := pr.db.Begin()

	project := &model.Project{
		Name:        createModel.Name,
		Description: createModel.Description,
		Content:     createModel.Content,
		Budget:      createModel.Budget,
	}

	if err := tx.Create(&project).Error; err != nil {
		tx.Rollback()
		return err
	}

	projectOwner := &model.ProjectOwner{
		UserID:    createModel.OwnerID,
		ProjectID: project.ID,
	}

	if err := tx.Create(&projectOwner).Error; err != nil {
		tx.Rollback()
		return err
	}

	if createModel.Public {
		projectPublication := &model.ProjectPublication{
			ProjectID: project.ID,
		}

		if err := tx.Create(&projectPublication).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit().Error
}

func (pr *projectRepository) Assign(userID uint, projectID uint, allowList *dto.Allow) error {
	if allowList.IsEmpty() {
		return pr.removePermissions(userID, projectID)
	}

	var count int64
	pr.db.Select(&model.ProjectOwner{ID: userID, ProjectID: projectID}).Count(&count)

	if count >= 1 {
		return errors.New("cannot assign owner")
	}

	return pr.upsertPermissions(userID, projectID, allowList)
}

func (pr *projectRepository) removePermissions(userID uint, projectID uint) error {
	return pr.db.Transaction(func(tx *gorm.DB) error {
		var userPermission model.ProjectUserPermission

		if err := tx.Where("user_id = ? AND project_id = ?", userID, projectID).First(&userPermission).Error; err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}

		if err := tx.Delete(&userPermission).Error; err != nil {
			return err
		}

		return nil
	})
}

func (pr *projectRepository) upsertPermissions(userID uint, projectID uint, allowList *dto.Allow) error {
	return pr.db.Transaction(func(tx *gorm.DB) error {
		var userPermission model.ProjectUserPermission

		if err := tx.Where("user_id = ? AND project_id = ?", userID, projectID).First(&userPermission).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				userPermission = model.ProjectUserPermission{
					UserID:    userID,
					ProjectID: projectID,
					Allow:     *allowList,
				}

				if err := tx.Create(&userPermission).Error; err != nil {
					return err
				}

				return nil
			}

			return err
		}

		updates := map[string]bool{
			"allow_view":            allowList.View,
			"allow_edit":            allowList.Edit,
			"allow_delete":          allowList.Delete,
			"allow_publish":         allowList.Publish,
			"allow_share":           allowList.Share,
			"allow_manage_users":    allowList.Manage.Users,
			"allow_manage_metadata": allowList.Manage.Metadata,
		}

		if err := tx.Model(&userPermission).Updates(updates).Error; err != nil {
			return err
		}

		return nil
	})
}

func (pr *projectRepository) Get(projectModel *model.Project) (*dto.ProjectInfo, error) {
	var project model.Project
	var projectOwner model.ProjectOwner
	var projectPublication model.ProjectPublication
	var allowedUsers []model.ProjectUserPermission

	if err := pr.db.Where("id = ?", projectModel.ID).First(&project).Error; err != nil {
		return nil, err
	}

	if err := pr.db.Where("project_id = ?", project.ID).First(&projectOwner).Error; err != nil {
		return nil, err
	}

	isPublic := false
	if err := pr.db.Where("project_id = ?", project.ID).First(&projectPublication).Error; err == nil {
		isPublic = true
	}

	if err := pr.db.Where("project_id = ?", project.ID).Find(&allowedUsers).Error; err != nil {
		return nil, err
	}

	ownerModel := &model.User{ID: projectOwner.UserID}
	ownerProfile, err := pr.userRepo.Get(ownerModel)
	if err != nil {
		return nil, err
	}

	owner := dto.UserInfoLite{
		Username:       ownerProfile.Username,
		ProfilePicture: ownerProfile.ProfilePicture,
	}

	var allowedUserDTOs []dto.ProjectUserPermissions
	for _, userPerm := range allowedUsers {
		userModel := &model.User{ID: userPerm.UserID}
		userProfile, err := pr.userRepo.Get(userModel)
		if err != nil {
			return nil, err
		}

		allowedUserDTOs = append(allowedUserDTOs, dto.ProjectUserPermissions{
			UserInfoLite: dto.UserInfoLite{
				Username:       userProfile.Username,
				ProfilePicture: userProfile.ProfilePicture,
			},
			Allow: dto.Allow{
				View:    userPerm.Allow.View,
				Edit:    userPerm.Allow.Edit,
				Delete:  userPerm.Allow.Delete,
				Publish: userPerm.Allow.Publish,
				Share:   userPerm.Allow.Share,
				Manage: dto.AllowManage{
					Users:    userPerm.Allow.Manage.Users,
					Metadata: userPerm.Allow.Manage.Metadata,
				},
			},
		})
	}

	projectInfo := &dto.ProjectInfo{
		Name:         project.Name,
		Description:  project.Description,
		Content:      project.Content,
		Budget:       project.Budget,
		Public:       isPublic,
		Owner:        owner,
		AllowedUsers: allowedUserDTOs,
	}

	return projectInfo, nil
}

func (pr *projectRepository) GetByOwner(userID uint, onlyPublic bool, page, pageSize int) (*dto.Pagination[dto.ProjectInfo], error) {
	panic("TODO: implement!!")
}

func (pr *projectRepository) GetPublic(page int, perPage int) (*dto.Pagination[dto.ProjectInfo], error) {
	panic("TODO: implement!!")
}

func (pr *projectRepository) SearchByName(name string, page, perpage int) (*dto.Pagination[dto.ProjectInfo], error) {
	panic("TODO: implement!!")
}

func (pr *projectRepository) GetContent() (any, error) {
	var content any

	result := pr.db.Model(&model.Project{}).Pluck("content", &content)

	if result.Error != nil {
		return nil, result.Error
	}

	return content, nil
}

func (pr *projectRepository) SaveContent(content any) error {
	return pr.db.Updates(&model.Project{
		Content: content,
	}).Error
}

func (pr *projectRepository) SafeDelete(id uint) error {
	return pr.db.Delete(&model.Project{ID: id}).Error
}

func (pr *projectRepository) UnsafeDelete(id uint) error {
	var project model.Project

	err := pr.db.Unscoped().Where("id = ?", id).First(&project).Error
	if err != nil {
		return err
	}

	if project.DeletedAt.Valid {
		return pr.db.Unscoped().Delete(&project).Error
	}

	return errors.New("project is not trashed")
}
