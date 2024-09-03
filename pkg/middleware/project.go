package middleware

import (
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/devkcud/arkhon-foundation/arkhon-api/internal/model/dto"
	"github.com/devkcud/arkhon-foundation/arkhon-api/internal/service"
	"github.com/devkcud/arkhon-foundation/arkhon-api/translations"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func ProjectLookup(ctx *gin.Context) {
	dict := translations.GetTranslation(ctx)

	projectID, err := strconv.ParseUint(ctx.Param("project"), 10, 64)
	if err != nil {
		log.Print(err)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": dict.ProjectInvalid})
		return
	}

	project, err := service.Project.GetByID(uint(projectID))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": dict.ProjectNotFound})
			return
		}

		log.Print(err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": dict.InternalServerError})
		return
	}

	if !project.Published && project.Owner != ctx.Keys["auth_user"].(*dto.UserProfile).Username {
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": dict.ProjectNotFound})
		return
	}

	ctx.Set("project_lookup", project)
	ctx.Next()
}

func ProjectOwnership(ctx *gin.Context) {
	dict := translations.GetTranslation(ctx)

	project := ctx.Keys["project_lookup"].(*dto.ProjectInformation)
	if project.Owner != ctx.Keys["auth_user"].(*dto.UserProfile).Username {
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": dict.ProjectNotFound})
		return
	}

	ctx.Next()
}