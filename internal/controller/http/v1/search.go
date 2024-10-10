package v1

import (
	"errors"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/swibly/swibly-api/internal/model/dto"
	"github.com/swibly/swibly-api/internal/service"
	"github.com/swibly/swibly-api/pkg/middleware"
	"github.com/swibly/swibly-api/translations"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func newSearchRoutes(handler *gin.RouterGroup) {
	h := handler.Group("/search")
	h.Use(middleware.APIKeyHasEnabledSearch, middleware.Auth)
	{
		h.GET("/user", SearchUserByNameHandler)
		h.GET("/project", SearchProjectByNameHandler)
	}
}

func SearchUserByNameHandler(ctx *gin.Context) {
	dict := translations.GetTranslation(ctx)

	name := ctx.Query("name")

	if !regexp.MustCompile(`[a-zA-Z ]`).MatchString(name) || strings.TrimSpace(name) == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": dict.SearchIncorrect})
		return
	}

	var (
		page    int = 1
		perpage int = 10
	)

	if i, e := strconv.Atoi(ctx.Query("page")); e == nil && ctx.Query("page") != "" {
		page = i
	}

	if i, e := strconv.Atoi(ctx.Query("perpage")); e == nil && ctx.Query("perpage") != "" {
		perpage = i
	}

	users, err := service.User.SearchByName(name, page, perpage)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": dict.SearchNoResults})
			return
		}

		log.Print(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": dict.InternalServerError})
		return
	}

	ctx.JSON(http.StatusOK, users)
}

func SearchProjectByNameHandler(ctx *gin.Context) {
	dict := translations.GetTranslation(ctx)

	name := ctx.Query("name")

	if !regexp.MustCompile(`[a-zA-Z ]`).MatchString(name) || strings.TrimSpace(name) == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": dict.SearchIncorrect})
		return
	}

	issuer := ctx.Keys["auth_user"].(*dto.UserProfile)

	var (
		page    int = 1
		perpage int = 10
	)

	if i, e := strconv.Atoi(ctx.Query("page")); e == nil && ctx.Query("page") != "" {
		page = i
	}

	if i, e := strconv.Atoi(ctx.Query("perpage")); e == nil && ctx.Query("perpage") != "" {
		perpage = i
	}

	projects, err := service.Project.SearchByName(issuer.ID, name, page, perpage)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": dict.SearchNoResults})
			return
		}

		log.Print(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": dict.InternalServerError})
		return
	}

	ctx.JSON(http.StatusOK, projects)
}
