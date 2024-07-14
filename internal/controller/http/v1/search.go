package v1

import (
	"errors"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/devkcud/arkhon-foundation/arkhon-api/internal/service"
	"github.com/devkcud/arkhon-foundation/arkhon-api/pkg/middleware"
	"github.com/devkcud/arkhon-foundation/arkhon-api/translations"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func newSearchRoutes(handler *gin.RouterGroup) {
	h := handler.Group("/search")
	h.Use(middleware.APIKeyHasEnabledSearch)
	{
		h.GET("/user", SearchByNameHandler)
	}
}

func SearchByNameHandler(ctx *gin.Context) {
	dict := translations.GetLang(ctx)

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

	users, err := service.User.GetBySimilarName(name, page, perpage)
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
