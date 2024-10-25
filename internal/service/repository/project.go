package repository

import (
	"encoding/json"
	"errors"

	"github.com/swibly/swibly-api/internal/model"
	"github.com/swibly/swibly-api/internal/model/dto"
	"github.com/swibly/swibly-api/pkg/aws"
	"github.com/swibly/swibly-api/pkg/db"
	"github.com/swibly/swibly-api/pkg/pagination"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type projectRepository struct {
	db *gorm.DB

	userRepo UserRepository
}

type ProjectRepository interface {
	Create(createModel *dto.ProjectCreation) (uint, error)
	Update(projectID uint, updateModel *dto.ProjectUpdate) error
	Unlink(projectID uint) error

	Assign(userID uint, projectID uint, allowList *dto.ProjectAssign) error

	Get(issuerID uint, projectModel *model.Project) (*dto.ProjectInfo, error)
	GetByOwner(issuerID, userID uint, onlyPublic bool, page, perPage int) (*dto.Pagination[dto.ProjectInfo], error)
	GetPublic(issuerID uint, page, perPage int) (*dto.Pagination[dto.ProjectInfo], error)
	GetFavorited(issuerID, userID uint, onlyPublic bool, page, perPage int) (*dto.Pagination[dto.ProjectInfo], error)
	GetTrashed(ownerID uint, page, perPage int) (*dto.Pagination[dto.ProjectInfo], error)

	SearchByName(issuerID uint, name string, page, perPage int) (*dto.Pagination[dto.ProjectInfo], error)

	GetContent(projectID uint) (any, error)
	SaveContent(projectID uint, content any) error

	Favorite(projectID, userID uint) error
	Unfavorite(projectID, userID uint) error

	SafeDelete(projectID uint) error
	Restore(projectID uint) error
	UnsafeDelete(projectID uint) error
	ClearTrash(userID uint) error
}

var (
	ErrProjectTrashed          = errors.New("project is trashed")
	ErrProjectNotTrashed       = errors.New("project is not trashed")
	ErrProjectAlreadyTrashed   = errors.New("project is already trashed")
	ErrProjectAlreadyFavorited = errors.New("project is already favorited by the user")
	ErrProjectNotFavorited     = errors.New("cannot unfavorite a project that is not favorited")
	ErrProjectIsNotAFork       = errors.New("project is not a fork")
	ErrUpstreamNotPublic       = errors.New("cannot publish this project because the upstream project is not public")
	ErrCannotAssignOwner       = errors.New("cannot assign owner")
	ErrUserNotAssigned         = errors.New("user is not assigned to the project")
)

func NewProjectRepository(userRepo UserRepository) ProjectRepository {
	return &projectRepository{db.Postgres, userRepo}
}

func (pr *projectRepository) baseProjectQuery(issuerID uint) *gorm.DB {
	return pr.db.Table("projects p").
		Select(`
			p.id as id,
			p.created_at as created_at,
			p.updated_at as updated_at,
			p.deleted_at as deleted_at,
			p.name as name,
			p.description as description,
			p.budget as budget,
      p.width as width,
      p.height as height,
      p.banner_url as banner_url,
			p.fork as fork,
			u.id AS owner_id,
			u.first_name AS owner_first_name,
			u.last_name AS owner_last_name,
			u.username AS owner_username,
			u.profile_picture AS owner_profile_picture,
			u.verified AS owner_verified,
			EXISTS (
				SELECT 1 
				FROM project_publications pp 
				WHERE pp.project_id = p.id
			) AS is_public,
			COALESCE((
				SELECT json_agg(
					json_build_object(
						'id', pu.user_id,
						'firstname', puu.first_name,
						'lastname', puu.last_name,
						'username', puu.username,
						'pfp', puu.profile_picture,
						'verified', puu.verified,
						'allow_view', pu.allow_view,
						'allow_edit', pu.allow_edit,
						'allow_delete', pu.allow_delete,
						'allow_publish', pu.allow_publish,
						'allow_share', pu.allow_share,
						'allow_manage_users', pu.allow_manage_users,
						'allow_manage_metadata', pu.allow_manage_metadata
					)
				)
				FROM project_user_permissions pu
				JOIN users puu ON pu.user_id = puu.id
				WHERE pu.project_id = p.id
			), '[]') AS allowed_users,
			COALESCE((
				SELECT true
				FROM project_user_favorites f
				WHERE f.project_id = p.id
				AND f.user_id = ?
				LIMIT 1
			), false) AS is_favorited,
			(
				SELECT COUNT(*)
				FROM project_user_favorites f
				WHERE f.project_id = p.id
			) AS total_favorites
		`, issuerID).
		Joins("JOIN project_owners po ON po.project_id = p.id").
		Joins("JOIN users u ON po.user_id = u.id")
}

func (pr *projectRepository) paginateProjects(query *gorm.DB, page, perPage int) (*dto.Pagination[dto.ProjectInfo], error) {
	paginationResult, err := pagination.Generate[dto.ProjectInfoJSON](query, page, perPage)
	if err != nil {
		return nil, err
	}

	projectInfoList := make([]*dto.ProjectInfo, 0, len(paginationResult.Data))
	for _, projectInfoJSON := range paginationResult.Data {
		projectInfo, err := convertToProjectInfo(projectInfoJSON)
		if err != nil {
			return nil, err
		}

		projectInfoList = append(projectInfoList, &projectInfo)
	}

	return &dto.Pagination[dto.ProjectInfo]{
		Data:         projectInfoList,
		TotalRecords: paginationResult.TotalRecords,
		TotalPages:   paginationResult.TotalPages,
		CurrentPage:  paginationResult.CurrentPage,
		NextPage:     paginationResult.NextPage,
		PreviousPage: paginationResult.PreviousPage,
	}, nil
}

func convertToProjectInfo(jsonInfo *dto.ProjectInfoJSON) (dto.ProjectInfo, error) {
	var allowedUsers []dto.ProjectUserPermissions
	err := json.Unmarshal(jsonInfo.AllowedUsers, &allowedUsers)
	if err != nil {
		return dto.ProjectInfo{}, err
	}

	return dto.ProjectInfo{
		ID:                  jsonInfo.ID,
		CreatedAt:           jsonInfo.CreatedAt,
		UpdatedAt:           jsonInfo.UpdatedAt,
		DeletedAt:           jsonInfo.DeletedAt,
		Name:                jsonInfo.Name,
		Description:         jsonInfo.Description,
		Budget:              jsonInfo.Budget,
		Width:               jsonInfo.Width,
		Height:              jsonInfo.Height,
		BannerURL:           jsonInfo.BannerURL,
		IsPublic:            jsonInfo.IsPublic,
		Fork:                jsonInfo.Fork,
		OwnerID:             jsonInfo.OwnerID,
		OwnerFirstName:      jsonInfo.OwnerFirstName,
		OwnerLastName:       jsonInfo.OwnerLastName,
		OwnerUsername:       jsonInfo.OwnerUsername,
		OwnerProfilePicture: jsonInfo.OwnerProfilePicture,
		OwnerVerified:       jsonInfo.OwnerVerified,
		IsFavorited:         jsonInfo.IsFavorited,
		TotalFavorites:      jsonInfo.TotalFavorites,
		AllowedUsers:        allowedUsers,
	}, nil
}

func (pr *projectRepository) Create(createModel *dto.ProjectCreation) (uint, error) {
	tx := pr.db.Begin()

	out, err := json.MarshalIndent(createModel.Content, "", "")
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	project := &model.Project{
		Name:        createModel.Name,
		Description: createModel.Description,
		Width:       createModel.Width,
		Height:      createModel.Height,
		Budget:      createModel.Budget,
		Content:     string(out),
		Fork:        createModel.Fork,
	}

	if err := tx.Create(&project).Error; err != nil {
		tx.Rollback()
		return 0, err
	}

	if createModel.BannerImage != nil {
		url, err := aws.UploadProjectImage(project.ID, createModel.BannerImage)
		if err != nil {
			tx.Rollback()
			return 0, err
		}

		if err := tx.Model(&model.Project{}).Where("id = ?", project.ID).Update("banner_url", url).Error; err != nil {
			tx.Rollback()
			return 0, err
		}
	}

	projectOwner := &model.ProjectOwner{
		UserID:    createModel.OwnerID,
		ProjectID: project.ID,
	}

	if err := tx.Create(&projectOwner).Error; err != nil {
		tx.Rollback()
		return 0, err
	}

	if createModel.Public && createModel.Fork == nil {
		projectPublication := &model.ProjectPublication{
			ProjectID: project.ID,
		}

		if err := tx.Create(&projectPublication).Error; err != nil {
			tx.Rollback()
			return 0, err
		}
	}

	if createModel.Fork != nil {
		fork, err := pr.Get(createModel.OwnerID, &model.Project{ID: *createModel.Fork})
		if err != nil {
			tx.Rollback()
			return 0, err
		}

		if err := tx.Create(&model.ProjectUserPermission{ProjectID: fork.ID, UserID: fork.OwnerID, Allow: dto.Allow{View: true}}).Error; err != nil {
			tx.Rollback()
			return 0, err
		}
	}

	if err := tx.Commit().Error; err != nil {
		return 0, err
	}

	return project.ID, nil
}

func (pr *projectRepository) Update(projectID uint, updateModel *dto.ProjectUpdate) error {
	if pr.db.Where("id = ?", projectID).First(&model.Project{}).Error == gorm.ErrRecordNotFound {
		return ErrProjectTrashed
	}

	tx := pr.db.Begin()

	if updateModel.Published != nil {
		switch *updateModel.Published {
		case true:
			var project model.Project
			if err := tx.Where("id = ?", projectID).First(&project).Error; err != nil {
				tx.Rollback()
				return err
			}

			if project.Fork != nil {
				var upstreamProject model.Project
				if err := tx.Where("id = ?", *project.Fork).First(&upstreamProject).Error; err != nil {
					tx.Rollback()
					return err
				}

				var upstreamPublication model.ProjectPublication
				if err := tx.Where("project_id = ?", upstreamProject.ID).First(&upstreamPublication).Error; err != nil {
					if errors.Is(err, gorm.ErrRecordNotFound) {
						tx.Rollback()
						return ErrUpstreamNotPublic
					}
					return err
				}
			}

			if err := tx.Clauses(clause.OnConflict{DoNothing: true}).Create(&model.ProjectPublication{ProjectID: projectID}).Error; err != nil {
				tx.Rollback()
				return err
			}
		case false:
			if err := tx.Where("project_id = ?", projectID).Unscoped().Delete(&model.ProjectPublication{}).Error; err != nil {
				tx.Rollback()
				return err
			}

			err := tx.
				Where("project_id = ? AND user_id NOT IN (?) AND user_id NOT IN (?)",
					projectID,
					tx.Table("project_user_permissions").
						Select("user_id").
						Where("project_id = ? AND allow_view = true", projectID),
					tx.Table("project_owners").
						Select("user_id").
						Where("project_id = ?", projectID)).
				Unscoped().
				Delete(&model.ProjectUserFavorite{}).Error
			if err != nil {
				tx.Rollback()
				return err
			}
		}
	}

	updates := make(map[string]interface{})

	if updateModel.Name != nil {
		updates["name"] = *updateModel.Name
	}
	if updateModel.Description != nil {
		updates["description"] = *updateModel.Description
	}
	if updateModel.Content != nil {
		contentJSON, err := json.Marshal(updateModel.Content)
		if err != nil {
			tx.Rollback()
			return err
		}
		updates["content"] = string(contentJSON)
	}
	if updateModel.Budget != nil {
		updates["budget"] = *updateModel.Budget
	}

	if updateModel.BannerImage != nil {
		url, err := aws.UploadProjectImage(projectID, updateModel.BannerImage)
		if err != nil {
			tx.Rollback()
			return err
		}

		if err := tx.Model(&model.Project{}).Where("id = ?", projectID).Update("banner_url", url).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	if updateModel.Width != nil {
		updates["width"] = *updateModel.Width
	}

	if updateModel.Height != nil {
		updates["height"] = *updateModel.Height
	}

	if err := tx.Model(&model.Project{}).Where("id = ?", projectID).Updates(updates).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (pr *projectRepository) Unlink(projectID uint) error {
	if pr.db.Where("id = ? AND fork IS NOT NULL", projectID).First(&model.Project{}).Error == gorm.ErrRecordNotFound {
		return ErrProjectIsNotAFork
	}

	return pr.db.Model(&model.Project{}).Where("id = ?", projectID).Update("fork", nil).Error
}

func (pr *projectRepository) Assign(userID uint, projectID uint, allowList *dto.ProjectAssign) error {
	if allowList.IsEmpty() {
		return pr.removePermissions(userID, projectID)
	}

	var count int64
	pr.db.Model(&model.ProjectOwner{}).Where("user_id = ? AND project_id = ?", userID, projectID).Count(&count)

	if count >= 1 {
		return ErrCannotAssignOwner
	}

	return pr.upsertPermissions(userID, projectID, allowList)
}

func (pr *projectRepository) removePermissions(userID uint, projectID uint) error {
	return pr.db.Transaction(func(tx *gorm.DB) error {
		var userPermission model.ProjectUserPermission

		if err := tx.Where("user_id = ? AND project_id = ?", userID, projectID).First(&userPermission).Error; err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		} else if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrUserNotAssigned
		}

		if err := tx.Delete(&userPermission).Error; err != nil {
			return err
		}

		if err := tx.Exec(`DELETE FROM project_user_favorites
			WHERE project_id = ? AND user_id = ? AND NOT EXISTS (
				SELECT 1 FROM project_publications 
				WHERE project_publications.project_id = project_user_favorites.project_id
			)`, projectID, userID).Error; err != nil {
			return err
		}

		return nil
	})
}

func (pr *projectRepository) upsertPermissions(userID uint, projectID uint, allowList *dto.ProjectAssign) error {
	return pr.db.Transaction(func(tx *gorm.DB) error {
		var userPermission model.ProjectUserPermission

		err := tx.Where("user_id = ? AND project_id = ?", userID, projectID).First(&userPermission).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			userPermission = model.ProjectUserPermission{
				UserID:    userID,
				ProjectID: projectID,
				Allow:     dto.Allow{},
			}

			if allowList.View != nil {
				userPermission.Allow.View = *allowList.View
			}
			if allowList.Edit != nil {
				userPermission.Allow.Edit = *allowList.Edit
			}
			if allowList.Delete != nil {
				userPermission.Allow.Delete = *allowList.Delete
			}
			if allowList.Publish != nil {
				userPermission.Allow.Publish = *allowList.Publish
			}
			if allowList.Share != nil {
				userPermission.Allow.Share = *allowList.Share
			}
			if allowList.ManageUsers != nil {
				userPermission.Allow.Manage.Users = *allowList.ManageUsers
			}
			if allowList.ManageMetadata != nil {
				userPermission.Allow.Manage.Metadata = *allowList.ManageMetadata
			}

			if err := tx.Create(&userPermission).Error; err != nil {
				return err
			}

			return nil
		} else if err != nil {
			return err
		}

		updates := make(map[string]interface{})

		if allowList.View != nil {
			updates["allow_view"] = *allowList.View
		}
		if allowList.Edit != nil {
			updates["allow_edit"] = *allowList.Edit
		}
		if allowList.Delete != nil {
			updates["allow_delete"] = *allowList.Delete
		}
		if allowList.Publish != nil {
			updates["allow_publish"] = *allowList.Publish
		}
		if allowList.Share != nil {
			updates["allow_share"] = *allowList.Share
		}
		if allowList.ManageUsers != nil {
			updates["allow_manage_users"] = *allowList.ManageUsers
		}
		if allowList.ManageMetadata != nil {
			updates["allow_manage_metadata"] = *allowList.ManageMetadata
		}

		if err := tx.Model(&userPermission).Updates(updates).Error; err != nil {
			return err
		}

		return nil
	})
}

func (pr *projectRepository) Get(userID uint, projectModel *model.Project) (*dto.ProjectInfo, error) {
	var project model.Project
	var projectOwner model.ProjectOwner
	var projectPublication model.ProjectPublication
	var allowedUsers []model.ProjectUserPermission

	if err := pr.db.Unscoped().Where("id = ?", projectModel.ID).First(&project).Error; err != nil {
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

	ownerProfile, err := pr.userRepo.Get(&model.User{ID: projectOwner.UserID})
	if err != nil {
		return nil, err
	}

	owner := dto.UserInfoLite{
		ID:             ownerProfile.ID,
		FirstName:      ownerProfile.FirstName,
		LastName:       ownerProfile.LastName,
		Username:       ownerProfile.Username,
		ProfilePicture: ownerProfile.ProfilePicture,
		Verified:       ownerProfile.Verified,
	}

	allowedUserDTOs := []dto.ProjectUserPermissions{}
	for _, userPerm := range allowedUsers {
		userModel := &model.User{ID: userPerm.UserID}
		userProfile, err := pr.userRepo.Get(userModel)
		if err != nil {
			return nil, err
		}

		allowedUserDTOs = append(allowedUserDTOs, dto.ProjectUserPermissions{
			ID:             userProfile.ID,
			FirstName:      userProfile.FirstName,
			LastName:       userProfile.LastName,
			Username:       userProfile.Username,
			ProfilePicture: userProfile.ProfilePicture,
			Verified:       userProfile.Verified,
			View:           userPerm.Allow.View,
			Edit:           userPerm.Allow.Edit,
			Delete:         userPerm.Allow.Delete,
			Share:          userPerm.Allow.Share,
			Publish:        userPerm.Allow.Publish,
			ManageUsers:    userPerm.Allow.Manage.Users,
			ManageMetadata: userPerm.Allow.Manage.Metadata,
		})
	}

	isFavorited := false
	if err := pr.db.Where("project_id = ? AND user_id = ?", project.ID, userID).First(&model.ProjectUserFavorite{}).Error; err == nil {
		isFavorited = true
	}

	var totalFavorites int64
	if err := pr.db.Model(&model.ProjectUserFavorite{}).Where("project_id = ?", project.ID).Count(&totalFavorites).Error; err != nil {
		return nil, err
	}

	projectInfo := &dto.ProjectInfo{
		ID:                  project.ID,
		CreatedAt:           project.CreatedAt,
		UpdatedAt:           project.UpdatedAt,
		DeletedAt:           project.DeletedAt,
		OwnerID:             owner.ID,
		OwnerFirstName:      owner.FirstName,
		OwnerLastName:       owner.LastName,
		OwnerUsername:       owner.Username,
		OwnerProfilePicture: owner.ProfilePicture,
		OwnerVerified:       owner.Verified,
		Name:                project.Name,
		Description:         project.Description,
		Budget:              project.Budget,
		Width:               project.Width,
		Height:              project.Height,
		BannerURL:           project.BannerURL,
		IsPublic:            isPublic,
		Fork:                project.Fork,
		IsFavorited:         isFavorited,
		TotalFavorites:      int(totalFavorites),
		AllowedUsers:        allowedUserDTOs,
	}

	return projectInfo, nil
}

func (pr *projectRepository) GetByOwner(issuerID, userID uint, onlyPublic bool, page, perPage int) (*dto.Pagination[dto.ProjectInfo], error) {
	query := pr.baseProjectQuery(issuerID).
		Where("deleted_at IS NULL").
		Where(`
			po.user_id = ? OR 
			EXISTS (
				SELECT 1 
				FROM project_user_permissions pu 
				WHERE pu.project_id = p.id 
				AND pu.user_id = ? 
				AND pu.allow_view = true
			)`, userID, userID)

	if onlyPublic {
		query = query.Joins("JOIN project_publications pp ON pp.project_id = p.id")
	}

	return pr.paginateProjects(query, page, perPage)
}

func (pr *projectRepository) GetPublic(issuerID uint, page, perPage int) (*dto.Pagination[dto.ProjectInfo], error) {
	query := pr.baseProjectQuery(issuerID).Where("deleted_at IS NULL").Where("EXISTS (SELECT 1 FROM project_publications pp WHERE pp.project_id = p.id)")

	return pr.paginateProjects(query, page, perPage)
}

func (pr *projectRepository) GetFavorited(issuerID, userID uint, onlyPublic bool, page, perPage int) (*dto.Pagination[dto.ProjectInfo], error) {
	query := pr.baseProjectQuery(issuerID).Where("deleted_at IS NULL").Joins("JOIN project_user_favorites f ON f.project_id = p.id AND f.user_id = ?", userID)

	if onlyPublic {
		query = query.Joins("JOIN project_publications pp ON pp.project_id = p.id")
	}

	return pr.paginateProjects(query, page, perPage)
}

func (pr *projectRepository) GetTrashed(issuerID uint, page, perPage int) (*dto.Pagination[dto.ProjectInfo], error) {
	query := pr.baseProjectQuery(issuerID).
		Where("deleted_at IS NOT NULL").
		Where(`
			u.id = ? OR 
			EXISTS (
				SELECT 1
				FROM project_user_permissions pu
				WHERE pu.project_id = p.id
				AND pu.user_id = ?
				AND pu.allow_delete = true
			)
		`, issuerID, issuerID)

	return pr.paginateProjects(query, page, perPage)
}

func (pr *projectRepository) SearchByName(issuerID uint, name string, page, perPage int) (*dto.Pagination[dto.ProjectInfo], error) {
	query := pr.baseProjectQuery(issuerID).
		Where(`to_tsvector('simple', p.name || ' ' || p.description || ' ' || u.username) @@ plainto_tsquery('simple', ?)`, name)

	return pr.paginateProjects(query, page, perPage)
}

func (pr *projectRepository) GetContent(projectID uint) (any, error) {
	var content string

	result := pr.db.Model(&model.Project{}).Select("content").Where("id = ?", projectID).Scan(&content)

	if result.Error != nil {
		return nil, result.Error
	}

	var contentData any
	if err := json.Unmarshal([]byte(content), &contentData); err != nil {
		return nil, err
	}

	return contentData, nil
}

func (pr *projectRepository) SaveContent(projectID uint, content any) error {
	contentJSON, err := json.Marshal(content)
	if err != nil {
		return err
	}

	contentString := string(contentJSON)

	return pr.db.Model(&model.Project{}).
		Where("id = ?", projectID).
		Update("content", contentString).
		Error
}

func (pr *projectRepository) Favorite(projectID, userID uint) error {
	var favorite model.ProjectUserFavorite
	err := pr.db.Where("project_id = ? AND user_id = ?", projectID, userID).First(&favorite).Error

	if err == nil {
		return ErrProjectAlreadyFavorited
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	return pr.db.Create(&model.ProjectUserFavorite{
		ProjectID: projectID,
		UserID:    userID,
	}).Error
}

func (pr *projectRepository) Unfavorite(projectID, userID uint) error {
	var favorite model.ProjectUserFavorite
	err := pr.db.Where("project_id = ? AND user_id = ?", projectID, userID).First(&favorite).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return ErrProjectNotFavorited
	} else if err != nil {
		return err
	}

	return pr.db.Delete(&favorite).Error
}

func (pr *projectRepository) SafeDelete(id uint) error {
	var project model.Project

	err := pr.db.Unscoped().Where("id = ?", id).First(&project).Error
	if err != nil {
		return err
	}

	if !project.DeletedAt.Valid {
		return pr.db.Delete(&project).Error
	}

	return ErrProjectAlreadyTrashed
}

func (pr *projectRepository) Restore(id uint) error {
	var project model.Project

	err := pr.db.Unscoped().Where("id = ?", id).First(&project).Error
	if err != nil {
		return err
	}

	if project.DeletedAt.Valid {
		return pr.db.Unscoped().Model(&project).Update("deleted_at", nil).Error
	}

	return ErrProjectNotTrashed
}

func (pr *projectRepository) UnsafeDelete(id uint) error {
	var project model.Project

	err := pr.db.Unscoped().Where("id = ?", id).First(&project).Error
	if err != nil {
		return err
	}

	if !project.DeletedAt.Valid {
		return ErrProjectNotTrashed
	}

	tx := pr.db.Begin()

	if err := tx.Unscoped().Where("project_id = ?", id).Delete(&model.ProjectUserFavorite{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Unscoped().Where("project_id = ?", id).Delete(&model.ProjectUserPermission{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Unscoped().Where("project_id = ?", id).Delete(&model.ProjectPublication{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Unscoped().Where("project_id = ?", id).Delete(&model.ProjectOwner{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := aws.DeleteProjectImage(id); err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Unscoped().Delete(&project).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

func (pr *projectRepository) ClearTrash(userID uint) error {
	tx := pr.db.Begin()

	if err := tx.Unscoped().
		Where(`
			EXISTS (
				SELECT 1
				FROM project_owners po
				WHERE po.project_id = project_user_favorites.project_id
				AND po.user_id = ?
			)
			OR EXISTS (
				SELECT 1
				FROM project_user_permissions pu
				WHERE pu.project_id = project_user_favorites.project_id
				AND pu.user_id = ?
				AND pu.allow_delete = true
			)
		`, userID, userID).
		Delete(&model.ProjectUserFavorite{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Unscoped().
		Where(`
			EXISTS (
				SELECT 1
				FROM project_owners po
				WHERE po.project_id = project_user_permissions.project_id
				AND po.user_id = ?
			)
			OR EXISTS (
				SELECT 1
				FROM project_user_permissions pu
				WHERE pu.project_id = project_user_permissions.project_id
				AND pu.user_id = ?
				AND pu.allow_delete = true
			)
		`, userID, userID).
		Delete(&model.ProjectUserPermission{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Unscoped().
		Where(`
			EXISTS (
				SELECT 1
				FROM project_owners po
				WHERE po.project_id = project_publications.project_id
				AND po.user_id = ?
			)
			OR EXISTS (
				SELECT 1
				FROM project_user_permissions pu
				WHERE pu.project_id = project_publications.project_id
				AND pu.user_id = ?
				AND pu.allow_delete = true
			)
		`, userID, userID).
		Delete(&model.ProjectPublication{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Unscoped().
		Where(`
			EXISTS (
				SELECT 1
				FROM project_owners po
				WHERE po.project_id = project_owners.project_id
				AND po.user_id = ?
			)
			OR EXISTS (
				SELECT 1
				FROM project_user_permissions pu
				WHERE pu.project_id = project_owners.project_id
				AND pu.user_id = ?
				AND pu.allow_delete = true
			)
		`, userID, userID).
		Delete(&model.ProjectOwner{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	projects := []model.Project{}
	if err := tx.Unscoped().
		Where("deleted_at IS NOT NULL").
		Where(`
			EXISTS (
				SELECT 1
				FROM project_owners po
				WHERE po.project_id = projects.id
				AND po.user_id = ?
			)
			OR EXISTS (
				SELECT 1
				FROM project_user_permissions pu
				WHERE pu.project_id = projects.id
				AND pu.user_id = ?
				AND pu.allow_delete = true
			)
		`, userID, userID).
		Find(&projects).Error; err != nil {
		tx.Rollback()
		return err
	}

	for _, project := range projects {
		if err := aws.DeleteProjectImage(project.ID); err != nil {
			tx.Rollback()
			return err
		}
	}

	if err := tx.Unscoped().
		Where("deleted_at IS NOT NULL").
		Where(`
			EXISTS (
				SELECT 1
				FROM project_owners po
				WHERE po.project_id = projects.id
				AND po.user_id = ?
			)
			OR EXISTS (
				SELECT 1
				FROM project_user_permissions pu
				WHERE pu.project_id = projects.id
				AND pu.user_id = ?
				AND pu.allow_delete = true
			)
		`, userID, userID).
		Delete(&model.Project{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return err
	}

	return nil
}
