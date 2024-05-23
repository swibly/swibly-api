package v1

import (
	"log"
	"net/http"

	"github.com/devkcud/arkhon-foundation/arkhon-api/internal/service/usecase"
	"github.com/gin-gonic/gin"
)

func newSearchRoutes(handler *gin.RouterGroup) {
	h := handler.Group("/search")

	u := h.Group("/user")
	{
		u.GET("/name/:name", SearchByUsernameHandler)
	}
}

func SearchByUsernameHandler(ctx *gin.Context) {
	users, err := usecase.UserInstance.GetBySimilarUsername(ctx.Param("name"))
	if err != nil {
		log.Print(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error. Please, try again later."})
		return
	}

	ctx.JSON(http.StatusOK, users)
}
