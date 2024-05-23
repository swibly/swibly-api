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

	u := h.Group("/user")
	{
		u.GET("/name/:name", SearchByUsernameHandler)
	}
}

func SearchByUsernameHandler(ctx *gin.Context) {
	users, err := usecase.UserInstance.GetBySimilarName(ctx.Param("name"))
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
