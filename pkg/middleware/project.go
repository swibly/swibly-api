package middleware

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/swibly/swibly-api/config"
	"github.com/swibly/swibly-api/internal/model/dto"
	"github.com/swibly/swibly-api/internal/service"
	"github.com/swibly/swibly-api/translations"
	"gorm.io/gorm"
)

func ProjectLookup(ctx *gin.Context) {
	dict := translations.GetTranslation(ctx)

	issuer := ctx.Keys["auth_user"].(*dto.UserProfile)

	projectID, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		log.Print(err)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": dict.ProjectInvalid})
		return
	}

	project, err := service.Project.GetByID(issuer.ID, uint(projectID))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": dict.ProjectNotFound})
			return
		}

		log.Print(err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": dict.InternalServerError})
		return
	}

	ctx.Set("project_lookup", project)
	ctx.Next()
}

// middleware.ProjectLookup must be called before middleware.ProjectOwnership
func ProjectIsAllowed(requiredPermissions dto.Allow) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		dict := translations.GetTranslation(ctx)

		project := ctx.Keys["project_lookup"].(*dto.ProjectInfo)
		issuer := ctx.Keys["auth_user"].(*dto.UserProfile)

		if project.OwnerUsername != issuer.Username && !issuer.HasPermissions(config.Permissions.ManageProjects) {
			isAllowed := false

			fmt.Println(requiredPermissions.View && project.IsPublic)

			if requiredPermissions.View && project.IsPublic {
				isAllowed = true
			}

			for _, allowedUser := range project.AllowedUsers {
				if allowedUser.Username == issuer.Username {
					if (!requiredPermissions.View || allowedUser.View || project.IsPublic) &&
						(!requiredPermissions.Edit || allowedUser.Edit) &&
						(!requiredPermissions.Delete || allowedUser.Delete) &&
						(!requiredPermissions.Publish || allowedUser.Publish) &&
						(!requiredPermissions.Share || allowedUser.Share) &&
						(!requiredPermissions.Manage.Users || allowedUser.ManageUsers || project.IsPublic) &&
						(!requiredPermissions.Manage.Metadata || allowedUser.ManageMetadata) {
						isAllowed = true
						continue
					}
				}
			}

			if !isAllowed {
				ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": dict.ProjectMissingPermissions})
				return
			}
		}

		ctx.Next()
	}
}

// middleware.ProjectLookup must be called before middleware.ProjectOwnership
func ProjectOwnership(ctx *gin.Context) {
	dict := translations.GetTranslation(ctx)

	issuer := ctx.Keys["auth_user"].(*dto.UserProfile)

	project := ctx.Keys["project_lookup"].(*dto.ProjectInfo)
	if issuer.HasPermissions(config.Permissions.ManageProjects) || project.OwnerUsername != issuer.Username {
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": dict.ProjectNotFound})
		return
	}

	ctx.Next()
}
