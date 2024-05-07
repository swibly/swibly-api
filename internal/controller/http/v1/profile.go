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
		h.Use(func(ctx *gin.Context) {
			ctx.Set("uc", usecase.NewUserUseCase())
			ctx.Next()
		})

		h.GET("/:username", GetUserRoute)
	}
}

func GetUserRoute(ctx *gin.Context) {
	// We know it exists, no need to pass in exists variable
	usecaseInterface, _ := ctx.Get("uc")
	// We know it will always be a UserUseCase
	usecase, _ := usecaseInterface.(usecase.UserUseCase)

	username := ctx.Param("username")

	user, err := usecase.GetByUsername(username)

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
