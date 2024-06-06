package v1

import (
	"errors"
	"log"
	"net/http"

	"github.com/devkcud/arkhon-foundation/arkhon-api/internal/service/usecase"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func newSearchRoutes(handler *gin.RouterGroup) {
	h := handler.Group("/search")

	h.GET("/user", SearchByNameHandler)
}

func SearchByNameHandler(ctx *gin.Context) {
	name := ctx.Query("name")

	if name == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Cannot find by empty name"})
		return
	}

	users, err := usecase.UserInstance.GetBySimilarName(name)
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
