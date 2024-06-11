package v1

import (
	"errors"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/devkcud/arkhon-foundation/arkhon-api/internal/service"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func newSearchRoutes(handler *gin.RouterGroup) {
	h := handler.Group("/search")

	h.GET("/user", SearchByNameHandler)
}

func SearchByNameHandler(ctx *gin.Context) {
	name := ctx.Query("name")

	if !regexp.MustCompile(`[a-zA-Z ]`).MatchString(name) || strings.TrimSpace(name) == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Bad search"})
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
			ctx.JSON(http.StatusNotFound, gin.H{"error": "No user found with that name."})
			return
		}

		log.Print(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error. Please, try again later."})
		return
	}

	ctx.JSON(http.StatusOK, users)
}
