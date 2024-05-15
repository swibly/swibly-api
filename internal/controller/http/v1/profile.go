package v1

import (
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/devkcud/arkhon-foundation/arkhon-api/internal/model"
	"github.com/devkcud/arkhon-foundation/arkhon-api/internal/service/usecase"
	"github.com/devkcud/arkhon-foundation/arkhon-api/pkg/middleware"
	"github.com/devkcud/arkhon-foundation/arkhon-api/pkg/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func newProfileRoutes(handler *gin.RouterGroup) {
	h := handler.Group("/profile")
	{
		h.GET("/:username", GetProfileHandler)
		h.PATCH("/:idlookup", middleware.AuthMiddleware, UpdateProfileHandler)
	}
}

func GetProfileHandler(ctx *gin.Context) {
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

			"role": user.Role,

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

func UpdateProfileHandler(ctx *gin.Context) {
	// We know it exists, no need to pass in exists variable
	usecaseInterface, _ := ctx.Get("uc")
	// We know it will always be a UserUseCase
	usecase, _ := usecaseInterface.(usecase.UserUseCase)

	uid, _ := ctx.Get("userid")
	issuerID, _ := uid.(uint)

	lid := ctx.Param("idlookup")
	lookupID, err := strconv.Atoi(lid)

	if err != nil {
		log.Print(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Bad ID"})
		return
	}

	issuer, err := usecase.GetByID(issuerID)

	if err != nil {
		log.Print(err)

		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "No user found with that ID"})
			return
		}

		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error. Please, try again later."})
		return
	}

	lookup, err := usecase.GetByID(uint(lookupID))

	if err != nil {
		log.Print(err)

		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "No user found with that ID"})
			return
		}

		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error. Please, try again later."})
		return
	}

	var body model.ProfileUpdate

	if err := ctx.BindJSON(&body); err != nil {
		log.Print(err)
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Bad body format"})
		return
	}

	if errs := utils.ValidateStruct(body); errs != nil {
		err := utils.ValidateErrorMessage(errs[0])

		ctx.JSON(http.StatusBadRequest, gin.H{"error": gin.H{err.Param: err.Message}})
		return
	}

	if issuer.ID == lookup.ID {
		var newModel model.User

		newModel.FirstName = body.FirstName
		newModel.LastName = body.LastName
		newModel.Username = body.Username
		newModel.Email = body.Email

		newModel.Show.Profile = body.Show.Profile
		newModel.Show.Image = body.Show.Image
		newModel.Show.Comments = body.Show.Comments
		newModel.Show.Favorites = body.Show.Favorites
		newModel.Show.Projects = body.Show.Projects
		newModel.Show.Components = body.Show.Components
		newModel.Show.Followers = body.Show.Followers
		newModel.Show.Following = body.Show.Following
		newModel.Show.Inventory = body.Show.Inventory
		newModel.Show.Formations = body.Show.Formations

		newModel.Notification.InApp = body.Notification.InApp
		newModel.Notification.Email = body.Notification.Email

		if err := usecase.Update(lookup.ID, &newModel); err != nil {
			log.Print(err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error. Please, try again later."})
		}
		return
	}

	if issuer.Role != "admin" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "You cannot update another user"})
		return
	}
}
