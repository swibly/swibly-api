package repository

import (
	"encoding/json"
	"errors"

	"github.com/devkcud/arkhon-foundation/arkhon-api/internal/model"
	"github.com/devkcud/arkhon-foundation/arkhon-api/internal/model/dto"
	"github.com/devkcud/arkhon-foundation/arkhon-api/pkg/db"
	"github.com/devkcud/arkhon-foundation/arkhon-api/pkg/pagination"
	"github.com/devkcud/arkhon-foundation/arkhon-api/pkg/utils"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type projectRepository struct {
	db *gorm.DB

	userRepo UserRepository
}

type ProjectRepository interface {
	Create(createModel *dto.ProjectCreation) error
	Update(uint, *dto.ProjectUpdate) error

	Assign(userID uint, projectID uint, allowList *dto.Allow) error

	Get(uint, *model.Project) (*dto.ProjectInfo, error)
	GetByOwner(issuerID, userID uint, onlyPublic bool, page, perPage int) (*dto.Pagination[dto.ProjectInfo], error)
	GetByOwnerLikes(issuerID, userID uint, onlyPublic bool, page, perPage int) (*dto.Pagination[dto.ProjectInfo], error)
	GetByOwnerDislikes(issuerID, userID uint, onlyPublic bool, page, perPage int) (*dto.Pagination[dto.ProjectInfo], error)
	GetPublic(issuerID uint, page, perPage int) (*dto.Pagination[dto.ProjectInfo], error)
	GetTrashed(ownerID uint, page, perPage int) (*dto.Pagination[dto.ProjectInfo], error)

	SearchByName(issuerID uint, name string, page, perPage int) (*dto.Pagination[dto.ProjectInfo], error)

	GetContent(uint) (any, error)
	SaveContent(uint, any) error

	Like(uint, uint) error
	Dislike(uint, uint) error

	SafeDelete(uint) error
	Restore(uint) error
	UnsafeDelete(uint) error
	ClearTrash(uint) error
}

var (
	ErrProjectTrashed        = errors.New("project is trashed")
	ErrProjectNotTrashed     = errors.New("project is not trashed")
	ErrProjectAlreadyTrashed = errors.New("project is already trashed")
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
			u.id AS owner_id,
			u.username AS owner_username,
			u.profile_picture AS owner_profile_picture,
			COALESCE(pl.project_id IS NOT NULL, false) AS is_liked,
			COALESCE(pd.project_id IS NOT NULL, false) AS is_disliked,
			COALESCE(like_counts.total_likes, 0) AS total_likes,
			COALESCE(dislike_counts.total_dislikes, 0) AS total_dislikes,
			CASE
				WHEN COALESCE(like_counts.total_likes, 0) + COALESCE(dislike_counts.total_dislikes, 0) = 0
				THEN 0
				ELSE COALESCE(like_counts.total_likes, 0) * 1.0 / (COALESCE(like_counts.total_likes, 0) + COALESCE(dislike_counts.total_dislikes, 0))
			END AS like_dislike_ratio,
			EXISTS (
				SELECT 1 
				FROM project_publications pp 
				WHERE pp.project_id = p.id
			) AS is_public,
			COALESCE((
				SELECT json_agg(
					json_build_object(
						'id', pu.user_id,
						'username', puu.username,
						'profile_picture', puu.profile_picture,
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
			), '[]') AS allowed_users
		`).
		Joins("JOIN project_owners po ON po.project_id = p.id").
		Joins("JOIN users u ON po.user_id = u.id").
		Joins("LEFT JOIN project_likes pl ON pl.project_id = p.id AND pl.user_id = ?", issuerID).
		Joins("LEFT JOIN project_dislikes pd ON pd.project_id = p.id AND pd.user_id = ?", issuerID).
		Joins("LEFT JOIN (SELECT project_id, COUNT(*) AS total_likes FROM project_likes GROUP BY project_id) AS like_counts ON like_counts.project_id = p.id").
		Joins("LEFT JOIN (SELECT project_id, COUNT(*) AS total_dislikes FROM project_dislikes GROUP BY project_id) AS dislike_counts ON dislike_counts.project_id = p.id")
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
		IsPublic:            jsonInfo.IsPublic,
		OwnerID:             jsonInfo.OwnerID,
		OwnerUsername:       jsonInfo.OwnerUsername,
		OwnerProfilePicture: jsonInfo.OwnerProfilePicture,
		IsLiked:             jsonInfo.IsLiked,
		IsDisliked:          jsonInfo.IsDisliked,
		TotalLikes:          jsonInfo.TotalLikes,
		TotalDislikes:       jsonInfo.TotalDislikes,
		LikeDislikeRatio:    jsonInfo.LikeDislikeRatio,
		AllowedUsers:        allowedUsers,
	}, nil
}

func (pr *projectRepository) Create(createModel *dto.ProjectCreation) error {
	tx := pr.db.Begin()

	var contentJSON utils.JSON
	var err error
	if createModel.Content != nil {
		contentJSON, err = json.Marshal(createModel.Content)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	project := &model.Project{
		Name:        createModel.Name,
		Description: createModel.Description,
		Budget:      createModel.Budget,
		Content:     string(contentJSON),
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

func (pr *projectRepository) Update(projectID uint, updateModel *dto.ProjectUpdate) error {
	if pr.db.Where("id = ?", projectID).First(&model.Project{}).Error == gorm.ErrRecordNotFound {
		return ErrProjectTrashed
	}

	tx := pr.db.Begin()

	if updateModel.Published != nil {
		switch *updateModel.Published {
		case true:
			if err := tx.Clauses(clause.OnConflict{DoNothing: true}).Create(&model.ProjectPublication{ProjectID: projectID}).Error; err != nil {
				tx.Rollback()
				return err
			}
		case false:
			if err := tx.Where("project_id = ?", projectID).Unscoped().Delete(&model.ProjectPublication{}).Error; err != nil {
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

	if err := tx.Model(&model.Project{}).Where("id = ?", projectID).Updates(updates).Error; err != nil {
		tx.Rollback()
		return err
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

func (pr *projectRepository) Get(userID uint, projectModel *model.Project) (*dto.ProjectInfo, error) {
	var project model.Project
	var projectOwner model.ProjectOwner
	var projectPublication model.ProjectPublication
	var allowedUsers []model.ProjectUserPermission
	var projectLike model.ProjectLikes
	var projectDislike model.ProjectDislikes

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
		Username:       ownerProfile.Username,
		ProfilePicture: ownerProfile.ProfilePicture,
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
			Username:       userProfile.Username,
			ProfilePicture: userProfile.ProfilePicture,
			View:           userPerm.Allow.View,
			Edit:           userPerm.Allow.Edit,
			Delete:         userPerm.Allow.Delete,
			Share:          userPerm.Allow.Share,
			Publish:        userPerm.Allow.Publish,
			ManageUsers:    userPerm.Allow.Manage.Users,
			ManageMetadata: userPerm.Allow.Manage.Metadata,
		})
	}

	isLiked := false
	if err := pr.db.Where("project_id = ? AND user_id = ?", project.ID, userID).First(&projectLike).Error; err == nil {
		isLiked = true
	}

	isDisliked := false
	if err := pr.db.Where("project_id = ? AND user_id = ?", project.ID, userID).First(&projectDislike).Error; err == nil {
		isDisliked = true
	}

	var totalLikes int64
	if err := pr.db.Model(&model.ProjectLikes{}).Where("project_id = ?", project.ID).Count(&totalLikes).Error; err != nil {
		return nil, err
	}

	var totalDislikes int64
	if err := pr.db.Model(&model.ProjectDislikes{}).Where("project_id = ?", project.ID).Count(&totalDislikes).Error; err != nil {
		return nil, err
	}

	likeDislikeRatio := 0.0
	if totalLikes+totalDislikes > 0 {
		likeDislikeRatio = float64(totalLikes) / float64(totalLikes+totalDislikes)
	}

	projectInfo := &dto.ProjectInfo{
		ID:                  project.ID,
		CreatedAt:           project.CreatedAt,
		UpdatedAt:           project.UpdatedAt,
		DeletedAt:           project.DeletedAt,
		OwnerID:             owner.ID,
		OwnerUsername:       owner.Username,
		OwnerProfilePicture: owner.ProfilePicture,
		Name:                project.Name,
		Description:         project.Description,
		Budget:              project.Budget,
		IsPublic:            isPublic,
		IsLiked:             isLiked,
		IsDisliked:          isDisliked,
		TotalLikes:          int(totalLikes),
		TotalDislikes:       int(totalDislikes),
		LikeDislikeRatio:    likeDislikeRatio,
		AllowedUsers:        allowedUserDTOs,
	}

	return projectInfo, nil
}

func (pr *projectRepository) GetByOwner(issuerID, userID uint, onlyPublic bool, page, perPage int) (*dto.Pagination[dto.ProjectInfo], error) {
	query := pr.baseProjectQuery(issuerID).Where("po.user_id = ?", userID)

	if onlyPublic {
		query = query.Joins("JOIN project_publications pp ON pp.project_id = p.id")
	}

	return pr.paginateProjects(query, page, perPage)
}

func (pr *projectRepository) GetByOwnerLikes(issuerID, userID uint, onlyPublic bool, page, perPage int) (*dto.Pagination[dto.ProjectInfo], error) {
	query := pr.baseProjectQuery(issuerID).Joins("JOIN project_likes pl ON pl.project_id = p.id AND pl.user_id = ?", userID)

	if onlyPublic {
		query = query.Where("EXISTS (SELECT 1 FROM project_publications pp WHERE pp.project_id = p.id)")
	}

	return pr.paginateProjects(query, page, perPage)
}

func (pr *projectRepository) GetByOwnerDislikes(issuerID, userID uint, onlyPublic bool, page, perPage int) (*dto.Pagination[dto.ProjectInfo], error) {
	query := pr.baseProjectQuery(issuerID).Joins("JOIN project_dislikes pd ON pd.project_id = p.id AND pd.user_id = ?", userID)

	if onlyPublic {
		query = query.Where("EXISTS (SELECT 1 FROM project_publications pp WHERE pp.project_id = p.id)")
	}

	return pr.paginateProjects(query, page, perPage)
}

func (pr *projectRepository) GetPublic(issuerID uint, page, perPage int) (*dto.Pagination[dto.ProjectInfo], error) {
	query := pr.baseProjectQuery(issuerID).Where("deleted_at IS NULL").Where("EXISTS (SELECT 1 FROM project_publications pp WHERE pp.project_id = p.id)")

	return pr.paginateProjects(query, page, perPage)
}

func (pr *projectRepository) GetTrashed(issuerID uint, page, perPage int) (*dto.Pagination[dto.ProjectInfo], error) {
	query := pr.baseProjectQuery(issuerID).
		Unscoped().
		Where("p.deleted_at IS NOT NULL").
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
		Where(`to_tsvector('simple', p.name || ' ' || p.description || ' ' || u.username) @@ plainto_tsquery('simple', ?)`, name).
		Order("like_dislike_ratio DESC, total_likes DESC")

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

func (pr *projectRepository) Like(projectID, userID uint) error {
	if err := pr.db.Where("project_id = ? AND user_id = ?", projectID, userID).Delete(&model.ProjectDislikes{}).Error; err != nil {
		return err
	}

	var existingLike model.ProjectLikes
	if err := pr.db.Where("project_id = ? AND user_id = ?", projectID, userID).First(&existingLike).Error; err == nil {
		return nil
	}

	newLike := model.ProjectLikes{
		ProjectID: projectID,
		UserID:    userID,
	}

	if err := pr.db.Create(&newLike).Error; err != nil {
		return err
	}

	return nil
}

func (pr *projectRepository) Dislike(projectID, userID uint) error {
	if err := pr.db.Where("project_id = ? AND user_id = ?", projectID, userID).Delete(&model.ProjectLikes{}).Error; err != nil {
		return err
	}

	var existingDislike model.ProjectDislikes
	if err := pr.db.Where("project_id = ? AND user_id = ?", projectID, userID).
		First(&existingDislike).Error; err == nil {
		return nil
	}

	newDislike := model.ProjectDislikes{
		ProjectID: projectID,
		UserID:    userID,
	}

	if err := pr.db.Create(&newDislike).Error; err != nil {
		return err
	}

	return nil
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

	if project.DeletedAt.Valid {
		return pr.db.Unscoped().Delete(&project).Error
	}

	return ErrProjectNotTrashed
}

func (pr *projectRepository) ClearTrash(userID uint) error {
	err := pr.db.Unscoped().
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
		Delete(&model.Project{}).Error

	if err != nil {
		return err
	}

	return nil
}
