package v1

import (
	"log"
	"net/http"

	"github.com/devkcud/arkhon-foundation/arkhon-api/internal/service/usecase"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func newProfileRoutes(handler *gin.RouterGroup) {
	h := handler.Group("/profile")
	{
		h.GET("/view/:username", GetProfileHandler)
	}
}

func GetProfileHandler(ctx *gin.Context) {
	username := ctx.Param("username")

	user, err := usecase.UserInstance.GetByUsername(username)

	if err == nil {
		ctx.JSON(http.StatusOK, gin.H{
			"id":        user.ID,
			"createdat": user.CreatedAt,
			"updatedat": user.UpdatedAt,

			"firstname": user.FirstName,
			"lastname":  user.LastName,
			"bio":       user.Bio,
			"verified":  user.Verified,

			"username": user.Username,
			"email":    user.Email,

			"xp":      user.XP,
			"arkhoin": user.Arkhoin,

			"show":         user.Show,
			"notification": user.Notification,
		})

		return
	}

	log.Print(err)

	if err == gorm.ErrRecordNotFound {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error. Please, try again later."})
}
