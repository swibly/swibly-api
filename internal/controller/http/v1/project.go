package v1

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/swibly/swibly-api/config"
	"github.com/swibly/swibly-api/internal/model/dto"
	"github.com/swibly/swibly-api/internal/service"
	"github.com/swibly/swibly-api/internal/service/repository"
	"github.com/swibly/swibly-api/pkg/aws"
	"github.com/swibly/swibly-api/pkg/middleware"
	"github.com/swibly/swibly-api/pkg/notification"
	"github.com/swibly/swibly-api/pkg/utils"
	"github.com/swibly/swibly-api/translations"
)

func newProjectRoutes(handler *gin.RouterGroup) {
	h := handler.Group("/projects")
	h.Use(middleware.APIKeyHasEnabledProjects, middleware.Auth)
	{
		h.GET("", GetPublicProjectsHandler)
		h.GET("/trash", GetTrashProjectsHandler)

		h.POST("", CreateProjectHandler)

		h.DELETE("/trash", DeleteTrashProjectsHandler)

		byUser := h.Group("/user/:username", middleware.UserLookup)
		{
			byUser.GET("", middleware.UserPrivacy(dto.UserShow{Projects: true}), GetProjectsByUserHandler)
			byUser.GET("/favorite", middleware.UserPrivacy(dto.UserShow{Favorites: true}), GetFavoriteProjectsByUserHandler)
		}
	}

	specific := h.Group("/:id", middleware.ProjectLookup)
	{
		specific.GET("", middleware.ProjectIsAllowed(dto.Allow{View: true}), GetProjectHandler)
		specific.GET("/content", middleware.ProjectIsAllowed(dto.Allow{View: true}), GetProjectContentHandler)

		specific.POST("/fork", middleware.ProjectIsAllowed(dto.Allow{View: true}), ForkProjectHandler)

		specific.PUT("/content", middleware.ProjectIsAllowed(dto.Allow{Edit: true}), UpdateProjectContentHandler)
		specific.PUT("/content/clear", middleware.ProjectIsAllowed(dto.Allow{Edit: true}), ClearProjectContentHandler)

		specific.PATCH("/update", middleware.ProjectIsAllowed(dto.Allow{Manage: dto.AllowManage{Metadata: true}}), UpdateProjectHandler)
		specific.PATCH("/publish", middleware.ProjectIsAllowed(dto.Allow{Publish: true}), PublishProjectHandler)
		specific.PATCH("/favorite", middleware.ProjectIsAllowed(dto.Allow{View: true}), FavoriteProjectHandler)

		specific.DELETE("/unpublish", middleware.ProjectIsAllowed(dto.Allow{Publish: true}), UnpublishProjectHandler)
		specific.DELETE("/unfavorite", middleware.ProjectIsAllowed(dto.Allow{View: true}), UnfavoriteProjectHandler)
		specific.DELETE("/fork", middleware.ProjectIsAllowed(dto.Allow{Manage: dto.AllowManage{Metadata: true}}), UnlinkProjectHandler)
		specific.DELETE("/leave", middleware.ProjectIsMember, LeaveProjectHandler)

		trashActions := specific.Group("/trash")
		{
			trashActions.PATCH("/restore", middleware.ProjectIsAllowed(dto.Allow{Delete: true}), RestoreProjectHandler)

			trashActions.DELETE("", middleware.ProjectIsAllowed(dto.Allow{Delete: true}), DeleteProjectHandler)
			trashActions.DELETE("/force", middleware.ProjectIsAllowed(dto.Allow{Delete: true}), DeleteProjectForceHandler)
		}

		assignActions := specific.Group("/assign/:username", middleware.UserLookup)
		{
			assignActions.PUT("", middleware.ProjectIsAllowed(dto.Allow{Manage: dto.AllowManage{Users: true}}), AssignProjectHandler)
			assignActions.DELETE("", middleware.ProjectIsAllowed(dto.Allow{Manage: dto.AllowManage{Users: true}}), UnassignProjectHandler)
		}
	}
}

func GetPublicProjectsHandler(ctx *gin.Context) {
	dict := translations.GetTranslation(ctx)

	issuerID := ctx.Keys["auth_user"].(*dto.UserProfile).ID

	page := 1
	perPage := 10

	if i, e := strconv.Atoi(ctx.Query("page")); e == nil && ctx.Query("page") != "" {
		page = i
	}

	if i, e := strconv.Atoi(ctx.Query("perpage")); e == nil && ctx.Query("perpage") != "" {
		perPage = i
	}

	projects, err := service.Project.GetPublic(issuerID, page, perPage)
	if err != nil {
		log.Print(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": dict.InternalServerError})
		return
	}

	ctx.JSON(http.StatusOK, projects)
}

func GetTrashProjectsHandler(ctx *gin.Context) {
	dict := translations.GetTranslation(ctx)

	issuerID := ctx.Keys["auth_user"].(*dto.UserProfile).ID

	page := 1
	perPage := 10

	if i, e := strconv.Atoi(ctx.Query("page")); e == nil && ctx.Query("page") != "" {
		page = i
	}

	if i, e := strconv.Atoi(ctx.Query("perpage")); e == nil && ctx.Query("perpage") != "" {
		perPage = i
	}

	projects, err := service.Project.GetTrashed(issuerID, page, perPage)
	if err != nil {
		log.Print(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": dict.InternalServerError})
		return
	}

	ctx.JSON(http.StatusOK, projects)
}

func CreateProjectHandler(ctx *gin.Context) {
	dict := translations.GetTranslation(ctx)

	issuer := ctx.Keys["auth_user"].(*dto.UserProfile)

	project := &dto.ProjectCreation{}
	if err := ctx.Bind(project); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": dict.InvalidBody})
		return
	}

	if errs := utils.ValidateStruct(project); errs != nil {
		err := utils.ValidateErrorMessage(ctx, errs[0])

		log.Print(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": gin.H{err.Param: err.Message}})
		return
	}

	publicStatus := strings.ToLower(ctx.Query("public"))
	if publicStatus == "true" || publicStatus == "t" || publicStatus == "1" {
		project.Public = true
	}

	project.OwnerID = issuer.ID

	if id, err := service.Project.Create(project); err != nil {
		if errors.Is(err, aws.ErrUnsupportedFileType) {
			log.Print(err)
			ctx.JSON(http.StatusBadRequest, gin.H{"error": dict.UnsupportedFileType})
			return
		}

		if errors.Is(err, aws.ErrUnableToDecode) {
			log.Print(err)
			ctx.JSON(http.StatusBadRequest, gin.H{"error": dict.UnableToDecodeFile})
			return
		}

		if errors.Is(err, aws.ErrUnableToEncode) {
			log.Print(err)
			ctx.JSON(http.StatusBadRequest, gin.H{"error": dict.UnableToDecodeFile})
			return
		}

		if errors.Is(err, aws.ErrFileTooLarge) {
			log.Print(err)
			ctx.JSON(http.StatusBadRequest, gin.H{"error": dict.FileTooLarge})
			return
		}

		log.Print(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": dict.InternalServerError})
		return
	} else {
		service.CreateNotification(dto.CreateNotification{
			Title:    dict.CategoryProject,
			Message:  fmt.Sprintf(dict.NotificationNewProjectCreated, project.Name),
			Type:     notification.Information,
			Redirect: utils.ToPtr(fmt.Sprintf(config.Redirects.Project, id)),
		}, issuer.ID)

		ctx.JSON(http.StatusOK, gin.H{"message": dict.ProjectCreated, "project": id})
	}
}

func DeleteTrashProjectsHandler(ctx *gin.Context) {
	dict := translations.GetTranslation(ctx)

	service.Project.ClearTrash(ctx.Keys["auth_user"].(*dto.UserProfile).ID)

	ctx.JSON(http.StatusOK, gin.H{"message": dict.TrashCleared})
}

func GetProjectsByUserHandler(ctx *gin.Context) {
	dict := translations.GetTranslation(ctx)

	issuer := ctx.Keys["auth_user"].(*dto.UserProfile)
	user := ctx.Keys["user_lookup"].(*dto.UserProfile)

	page := 1
	perPage := 10

	if i, e := strconv.Atoi(ctx.Query("page")); e == nil && ctx.Query("page") != "" {
		page = i
	}

	if i, e := strconv.Atoi(ctx.Query("perpage")); e == nil && ctx.Query("perpage") != "" {
		perPage = i
	}

	projects, err := service.Project.GetByOwner(issuer.ID, user.ID, issuer.ID != user.ID, page, perPage)
	if err != nil {
		log.Print(err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": dict.InternalServerError})
		return
	}

	ctx.JSON(http.StatusOK, projects)
}

func GetFavoriteProjectsByUserHandler(ctx *gin.Context) {
	dict := translations.GetTranslation(ctx)

	issuer := ctx.Keys["auth_user"].(*dto.UserProfile)
	user := ctx.Keys["user_lookup"].(*dto.UserProfile)

	page := 1
	perPage := 10

	if i, e := strconv.Atoi(ctx.Query("page")); e == nil && ctx.Query("page") != "" {
		page = i
	}

	if i, e := strconv.Atoi(ctx.Query("perpage")); e == nil && ctx.Query("perpage") != "" {
		perPage = i
	}

	projects, err := service.Project.GetFavorited(issuer.ID, user.ID, issuer.ID != user.ID, page, perPage)
	if err != nil {
		log.Print(err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": dict.InternalServerError})
		return
	}

	ctx.JSON(http.StatusOK, projects)
}

func GetProjectHandler(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, ctx.Keys["project_lookup"].(*dto.ProjectInfo))
}

func GetProjectContentHandler(ctx *gin.Context) {
	project := ctx.Keys["project_lookup"].(*dto.ProjectInfo)

	content, err := service.Project.GetContent(project.ID)
	if err != nil {
		log.Print(err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, content)
}

func ForkProjectHandler(ctx *gin.Context) {
	dict := translations.GetTranslation(ctx)

	issuer := ctx.Keys["auth_user"].(*dto.UserProfile)
	project := ctx.Keys["project_lookup"].(*dto.ProjectInfo)

	if id, err := service.Project.Fork(project.ID, issuer.ID); err != nil {
		if errors.Is(err, repository.ErrUpstreamNotPublic) {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": dict.UpstreamNotPublic})
			return
		}

		log.Print(err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": dict.InternalServerError})
		return
	} else {
		service.CreateNotification(dto.CreateNotification{
			Title:    dict.CategoryProject,
			Message:  fmt.Sprintf(dict.NotificationUserClonedYourProject, issuer.FirstName+" "+issuer.LastName, project.Name),
			Type:     notification.Information,
			Redirect: utils.ToPtr(fmt.Sprintf(config.Redirects.Project, id)),
		}, project.OwnerID)

		service.CreateNotification(dto.CreateNotification{
			Title:    dict.CategoryProject,
			Message:  fmt.Sprintf(dict.NotificationNewProjectCreated, project.Name),
			Type:     notification.Information,
			Redirect: utils.ToPtr(fmt.Sprintf(config.Redirects.Project, id)),
		}, issuer.ID)

		ctx.JSON(http.StatusOK, gin.H{"message": dict.ProjectForked, "project": id})
	}
}

func UnlinkProjectHandler(ctx *gin.Context) {
	dict := translations.GetTranslation(ctx)

	project := ctx.Keys["project_lookup"].(*dto.ProjectInfo)

	if err := service.Project.Unlink(project.ID); err != nil {
		if errors.Is(err, repository.ErrProjectIsNotAFork) {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": dict.ProjectIsNotAFork})
			return
		}

		log.Print(err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": dict.InternalServerError})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": dict.ProjectUnlinked})
}

func LeaveProjectHandler(ctx *gin.Context) {
	dict := translations.GetTranslation(ctx)

	issuer := ctx.Keys["auth_user"].(*dto.UserProfile)
	project := ctx.Keys["project_lookup"].(*dto.ProjectInfo)

	if err := service.Project.LeaveProject(issuer.ID, project.ID); err != nil {
		log.Print(err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": dict.InternalServerError})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": dict.ProjectUnassignedUser})
}

func UpdateProjectContentHandler(ctx *gin.Context) {
	dict := translations.GetTranslation(ctx)

	project := ctx.Keys["project_lookup"].(*dto.ProjectInfo)

	var body any
	if err := ctx.ShouldBindJSON(&body); err != nil {
		log.Print(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": dict.InvalidBody})
		return
	}

	if err := service.Project.SaveContent(project.ID, body); err != nil {
		log.Print(err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": dict.InternalServerError})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": dict.ProjectUpdated})
}

func ClearProjectContentHandler(ctx *gin.Context) {
	dict := translations.GetTranslation(ctx)

	project := ctx.Keys["project_lookup"].(*dto.ProjectInfo)

	if err := service.Project.ClearContent(project.ID); err != nil {
		log.Print(err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": dict.InternalServerError})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": dict.ProjectUpdated})
}

func UpdateProjectHandler(ctx *gin.Context) {
	dict := translations.GetTranslation(ctx)

	project := ctx.Keys["project_lookup"].(*dto.ProjectInfo)

	var body *dto.ProjectUpdate
	if err := ctx.Bind(&body); err != nil {
		log.Print(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": dict.InvalidBody})
		return
	}

	if errs := utils.ValidateStruct(body); errs != nil {
		err := utils.ValidateErrorMessage(ctx, errs[0])

		log.Print(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": gin.H{err.Param: err.Message}})
		return
	}

	if err := service.Project.Update(project.ID, body); err != nil {
		if errors.Is(err, aws.ErrUnsupportedFileType) {
			log.Print(err)
			ctx.JSON(http.StatusBadRequest, gin.H{"error": dict.UnsupportedFileType})
			return
		}

		if errors.Is(err, aws.ErrUnableToDecode) {
			log.Print(err)
			ctx.JSON(http.StatusBadRequest, gin.H{"error": dict.UnableToDecodeFile})
			return
		}

		if errors.Is(err, aws.ErrUnableToEncode) {
			log.Print(err)
			ctx.JSON(http.StatusBadRequest, gin.H{"error": dict.UnableToDecodeFile})
			return
		}

		if errors.Is(err, aws.ErrFileTooLarge) {
			log.Print(err)
			ctx.JSON(http.StatusBadRequest, gin.H{"error": dict.FileTooLarge})
			return
		}

		if errors.Is(err, repository.ErrProjectTrashed) {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": dict.ProjectAlreadyTrashed})
			return
		}

		log.Print(err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": dict.InternalServerError})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": dict.ProjectUpdated})
}

func PublishProjectHandler(ctx *gin.Context) {
	dict := translations.GetTranslation(ctx)

	project := ctx.Keys["project_lookup"].(*dto.ProjectInfo)

	if err := service.Project.Publish(project.ID); err != nil {
		if errors.Is(err, repository.ErrUpstreamNotPublic) {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": dict.UpstreamNotPublic})
			return
		}

		if errors.Is(err, repository.ErrProjectTrashed) {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": dict.ProjectAlreadyTrashed})
			return
		}

		log.Print(err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": dict.InternalServerError})
		return
	}

	ids := []uint{project.OwnerID}

	for _, user := range project.AllowedUsers {
		ids = append(ids, user.ID)
	}

	service.CreateNotification(dto.CreateNotification{
		Title:   dict.CategoryProject,
		Message: fmt.Sprintf(dict.NotificationYourProjectPublished, project.Name),
		Type:    notification.Warning,
	}, ids...)

	ctx.JSON(http.StatusOK, gin.H{"message": dict.ProjectPublished})
}

func UnpublishProjectHandler(ctx *gin.Context) {
	dict := translations.GetTranslation(ctx)

	project := ctx.Keys["project_lookup"].(*dto.ProjectInfo)

	if err := service.Project.Unpublish(project.ID); err != nil {
		if errors.Is(err, repository.ErrProjectTrashed) {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": dict.ProjectAlreadyTrashed})
			return
		}

		log.Print(err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": dict.InternalServerError})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": dict.ProjectUnpublished})
}

func FavoriteProjectHandler(ctx *gin.Context) {
	dict := translations.GetTranslation(ctx)

	project := ctx.Keys["project_lookup"].(*dto.ProjectInfo)
	issuer := ctx.Keys["auth_user"].(*dto.UserProfile)

	if err := service.Project.Favorite(project.ID, issuer.ID); err != nil {
		if errors.Is(err, repository.ErrProjectAlreadyFavorited) {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": dict.ProjectAlreadyFavorited})
			return
		}

		log.Print(err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": dict.InternalServerError})
		return
	}

	service.CreateNotification(dto.CreateNotification{
		Title:   dict.CategoryProject,
		Message: fmt.Sprintf(dict.NotificationYourProjectFavorited, project.Name, issuer.FirstName+" "+issuer.LastName),
		Type:    notification.Information,
	}, project.OwnerID)

	ctx.JSON(http.StatusOK, gin.H{"message": dict.ProjectFavorited})
}

func UnfavoriteProjectHandler(ctx *gin.Context) {
	dict := translations.GetTranslation(ctx)

	project := ctx.Keys["project_lookup"].(*dto.ProjectInfo)
	issuer := ctx.Keys["auth_user"].(*dto.UserProfile)

	if err := service.Project.Unfavorite(project.ID, issuer.ID); err != nil {
		if errors.Is(err, repository.ErrProjectNotFavorited) {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": dict.ProjectNotFavorited})
			return
		}

		log.Print(err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": dict.InternalServerError})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": dict.ProjectUnfavorited})
}

func DeleteProjectHandler(ctx *gin.Context) {
	dict := translations.GetTranslation(ctx)

	project := ctx.Keys["project_lookup"].(*dto.ProjectInfo)

	if err := service.Project.Trash(project.ID); err != nil {
		if errors.Is(err, repository.ErrProjectAlreadyTrashed) {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": dict.ProjectAlreadyTrashed})
			return
		}

		log.Print(err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": dict.InternalServerError})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": dict.ProjectTrashed})
}

func RestoreProjectHandler(ctx *gin.Context) {
	dict := translations.GetTranslation(ctx)

	project := ctx.Keys["project_lookup"].(*dto.ProjectInfo)

	if err := service.Project.Restore(project.ID); err != nil {
		if errors.Is(err, repository.ErrProjectNotTrashed) {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": dict.ProjectNotTrashed})
			return
		}

		log.Print(err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": dict.InternalServerError})
		return
	}

	service.CreateNotification(dto.CreateNotification{
		Title:    dict.CategoryProject,
		Message:  fmt.Sprintf(dict.NotificationRestoredProjectFromTrash, project.Name),
		Type:     notification.Warning,
		Redirect: utils.ToPtr(fmt.Sprintf(config.Redirects.Project, project.ID)),
	}, project.OwnerID)

	ctx.JSON(http.StatusOK, gin.H{"message": dict.ProjectRestored})
}

func DeleteProjectForceHandler(ctx *gin.Context) {
	dict := translations.GetTranslation(ctx)

	project := ctx.Keys["project_lookup"].(*dto.ProjectInfo)

	if err := service.Project.RemoveFromTrash(project.ID); err != nil {
		if errors.Is(err, repository.ErrProjectNotTrashed) {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": dict.ProjectNotTrashed})
			return
		}

		log.Print(err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": dict.InternalServerError})
		return
	}

	service.CreateNotification(dto.CreateNotification{
		Title:   dict.CategoryProject,
		Message: fmt.Sprintf(dict.NotificationDeletedProjectFromTrash, project.Name),
		Type:    notification.Danger,
	}, project.OwnerID)

	ctx.JSON(http.StatusOK, gin.H{"message": dict.ProjectDeleted})
}

func AssignProjectHandler(ctx *gin.Context) {
	dict := translations.GetTranslation(ctx)

	project := ctx.Keys["project_lookup"].(*dto.ProjectInfo)
	user := ctx.Keys["user_lookup"].(*dto.UserProfile)

	var body dto.ProjectAssign

	if err := ctx.BindJSON(&body); err != nil {
		log.Print(err)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": dict.InvalidBody})
		return
	}

	if errs := utils.ValidateStruct(body); errs != nil {
		err := utils.ValidateErrorMessage(ctx, errs[0])

		log.Print(err)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": gin.H{err.Param: err.Message}})
		return
	}

	allowList := &dto.ProjectAssign{
		View:           body.View,
		Edit:           body.Edit,
		Delete:         body.Delete,
		Publish:        body.Publish,
		Share:          body.Share,
		ManageUsers:    body.ManageUsers,
		ManageMetadata: body.ManageMetadata,
	}

	if allowList.IsEmpty() {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": dict.ProjectEmptyAssign})
		return
	}

	if err := service.Project.Assign(user.ID, project.ID, allowList); err != nil {
		if errors.Is(err, repository.ErrCannotAssignOwner) {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": dict.ProjectCannotAssignOwner})
			return
		}

		log.Print(err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": dict.InternalServerError})
		return
	}

	service.CreateNotification(dto.CreateNotification{
		Title:   dict.CategoryProject,
		Message: fmt.Sprintf(dict.NotificationAddedUserToProject, user.FirstName+user.LastName, project.Name),
		Type:    notification.Danger,
	}, project.OwnerID)

	service.CreateNotification(dto.CreateNotification{
		Title:    dict.CategoryProject,
		Message:  fmt.Sprintf(dict.NotificationAddedYouToProject, project.Name),
		Type:     notification.Information,
		Redirect: utils.ToPtr(fmt.Sprintf(config.Redirects.Project, project.ID)),
	}, user.ID)

	ctx.JSON(http.StatusOK, gin.H{"message": dict.ProjectAssignedUser})
}

func UnassignProjectHandler(ctx *gin.Context) {
	dict := translations.GetTranslation(ctx)

	project := ctx.Keys["project_lookup"].(*dto.ProjectInfo)
	user := ctx.Keys["user_lookup"].(*dto.UserProfile)

	if err := service.Project.Assign(user.ID, project.ID, &dto.ProjectAssign{}); err != nil {
		if errors.Is(err, repository.ErrUserNotAssigned) {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": dict.ProjectUserNotAssigned})
			return
		}

		log.Print(err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": dict.InternalServerError})
		return
	}

	service.CreateNotification(dto.CreateNotification{
		Title:   dict.CategoryProject,
		Message: fmt.Sprintf(dict.NotificationRemovedUserFromProject, user.FirstName+user.LastName, project.Name),
		Type:    notification.Danger,
	}, project.OwnerID)

	service.CreateNotification(dto.CreateNotification{
		Title:    dict.CategoryProject,
		Message:  fmt.Sprintf(dict.NotificationRemovedYouFromProject, project.Name),
		Type:     notification.Information,
		Redirect: utils.ToPtr(fmt.Sprintf(config.Redirects.Project, project.ID)),
	}, user.ID)

	ctx.JSON(http.StatusOK, gin.H{"message": dict.ProjectUnassignedUser})
}
