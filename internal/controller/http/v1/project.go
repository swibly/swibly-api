package v1

import (
	"log"
	"net/http"
	"strconv"

	"github.com/devkcud/arkhon-foundation/arkhon-api/internal/model/dto"
	"github.com/devkcud/arkhon-foundation/arkhon-api/internal/service"
	"github.com/devkcud/arkhon-foundation/arkhon-api/pkg/middleware"
	"github.com/devkcud/arkhon-foundation/arkhon-api/pkg/utils"
	"github.com/devkcud/arkhon-foundation/arkhon-api/translations"
	"github.com/gin-gonic/gin"
)

func newProjectRoutes(handler *gin.RouterGroup) {
	g := handler.Group("/projects")
	{
		g.GET("", ListPublicProjects)

		g.POST("/create", middleware.AuthMiddleware, CreateProject)
	}

	g.GET("/mine", middleware.AuthMiddleware, ListMyProjects)

	g.GET("/owner/:owner", middleware.OptionalAuthMiddleware, ListProjectsByOwner)

	specific := g.Group("/:project")
	specific.Use(middleware.AuthMiddleware, middleware.ProjectLookup)
	{
		specific.GET("", GetProjectInfo)
		specific.GET("/content", GetProjectContent)

		specific.PATCH("/content", middleware.ProjectOwnership, UpdateProject)

		specific.POST("/publish", middleware.ProjectOwnership, PublishProject)
		specific.POST("/unpublish", middleware.ProjectOwnership, UnpublishProject)

		specific.DELETE("", middleware.ProjectOwnership, DeleteProject)
	}
}

func CreateProject(ctx *gin.Context) {
	dict := translations.GetTranslation(ctx)

	issuer := ctx.Keys["auth_user"].(*dto.UserProfile)

	var body dto.ProjectCreation

	if err := ctx.BindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": dict.InvalidBody})
		return
	}

	if errs := utils.ValidateStruct(&body); errs != nil {
		err := utils.ValidateErrorMessage(ctx, errs[0])

		log.Print(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": gin.H{err.Param: err.Message}})
		return
	}

	body.Owner = issuer.Username

	err := service.Project.Create(&body)
	if err != nil {
		log.Print(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": dict.InternalServerError})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": dict.ProjectCreated})
}

func ListMyProjects(ctx *gin.Context) {
	dict := translations.GetTranslation(ctx)

	issuer := ctx.Keys["auth_user"].(*dto.UserProfile)

	projects, err := service.Project.GetByOwnerUsername(issuer.Username, true, 1, 10)
	if err != nil {
		log.Print(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": dict.InternalServerError})
		return
	}

	ctx.JSON(http.StatusOK, projects)
}

func ListPublicProjects(ctx *gin.Context) {
	dict := translations.GetTranslation(ctx)

	var (
		page    int = 1
		perpage int = 10
	)

	if ctx.Query("page") != "" {
		page, _ = strconv.Atoi(ctx.Query("page"))
	}

	if ctx.Query("perpage") != "" {
		perpage, _ = strconv.Atoi(ctx.Query("perpage"))
	}

	projects, err := service.Project.GetPublicAll(page, perpage)
	if err != nil {
		log.Print(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": dict.InternalServerError})
		return
	}

	ctx.JSON(http.StatusOK, projects)
}

func ListProjectsByOwner(ctx *gin.Context) {
	dict := translations.GetTranslation(ctx)

	var issuerUsername string = ""
	if u, exists := ctx.Get("auth_user"); exists {
		issuerUsername = u.(*dto.UserProfile).Username
	}

	ownerUsername := ctx.Param("owner")

	page := 1
	perPage := 10

	if i, e := strconv.Atoi(ctx.Query("page")); e == nil && ctx.Query("page") != "" {
		page = i
	}

	if i, e := strconv.Atoi(ctx.Query("perpage")); e == nil && ctx.Query("perpage") != "" {
		perPage = i
	}

	projects, err := service.Project.GetByOwnerUsername(ownerUsername, issuerUsername == ownerUsername, page, perPage)
	if err != nil {
		log.Print(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": dict.InternalServerError})
		return
	}

	ctx.JSON(http.StatusOK, projects)
}

func GetProjectInfo(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, ctx.Keys["project_lookup"])
}

func GetProjectContent(ctx *gin.Context) {
	content := service.Project.GetContent(ctx.Keys["project_lookup"].(*dto.ProjectInformation).ID)

	ctx.JSON(http.StatusOK, content)
}

func UpdateProject(ctx *gin.Context) {
	dict := translations.GetTranslation(ctx)

	var content any

	if err := ctx.BindJSON(&content); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": dict.InvalidBody})
		return
	}

	service.Project.SaveContent(ctx.Keys["project_lookup"].(*dto.ProjectInformation).ID, content)

	ctx.JSON(http.StatusOK, gin.H{"message": dict.ProjectUpdated})
}

func PublishProject(ctx *gin.Context) {
	dict := translations.GetTranslation(ctx)

	err := service.Project.Publish(ctx.Keys["project_lookup"].(*dto.ProjectInformation).ID)
	if err != nil {
		log.Print(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": dict.InternalServerError})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": dict.ProjectPublished})
}

func UnpublishProject(ctx *gin.Context) {
	dict := translations.GetTranslation(ctx)

	err := service.Project.Unpublish(ctx.Keys["project_lookup"].(*dto.ProjectInformation).ID)
	if err != nil {
		log.Print(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": dict.InternalServerError})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": dict.ProjectUnpublished})
}

func DeleteProject(ctx *gin.Context) {
	dict := translations.GetTranslation(ctx)

	err := service.Project.Delete(ctx.Keys["project_lookup"].(*dto.ProjectInformation).ID)
	if err != nil {
		log.Print(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": dict.InternalServerError})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": dict.ProjectDeleted})
}
