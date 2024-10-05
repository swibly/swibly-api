package v1

import (
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/devkcud/arkhon-foundation/arkhon-api/internal/model/dto"
	"github.com/devkcud/arkhon-foundation/arkhon-api/internal/service"
	"github.com/devkcud/arkhon-foundation/arkhon-api/pkg/middleware"
	"github.com/devkcud/arkhon-foundation/arkhon-api/pkg/utils"
	"github.com/devkcud/arkhon-foundation/arkhon-api/translations"
	"github.com/gin-gonic/gin"
)

func newProjectRoutes(handler *gin.RouterGroup) {
	h := handler.Group("/projects")
	h.Use(middleware.APIKeyHasEnabledProjects, middleware.Auth)
	{
		h.GET("", GetPublicProjectsHandler)

		h.POST("", CreateProjectHandler)
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

func CreateProjectHandler(ctx *gin.Context) {
	dict := translations.GetTranslation(ctx)

	issuer := ctx.Keys["auth_user"].(*dto.UserProfile)

	project := &dto.ProjectCreation{}
	if err := ctx.ShouldBindJSON(project); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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

	if err := service.Project.Create(project); err != nil {
		log.Print(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": dict.InternalServerError})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": dict.ProjectCreated})
}
