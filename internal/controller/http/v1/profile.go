package v1

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/devkcud/arkhon-foundation/arkhon-api/internal/model"
	"github.com/devkcud/arkhon-foundation/arkhon-api/internal/service/usecase"
	"github.com/devkcud/arkhon-foundation/arkhon-api/pkg/middleware"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func newProfileRoutes(handler *gin.RouterGroup) {
	h := handler.Group("/profile")
	{
		h.GET("/view/:username", middleware.OptionalAuthMiddleware, GetProfileHandler)
		h.GET("/view/:username/followers", middleware.OptionalAuthMiddleware, GetFollowersHandler)
		h.GET("/view/:username/following", middleware.OptionalAuthMiddleware, GetFollowingHandler)
	}
}

func GetProfileHandler(ctx *gin.Context) {
	var issuer *model.User

	idFromJWT, exists := ctx.Get("id_from_jwt")
	if exists {
		log.Print(idFromJWT)
		id, err := strconv.Atoi(fmt.Sprintf("%v", idFromJWT))
		if err != nil {
			log.Print(err)
		} else {
			issuer, err = usecase.UserInstance.GetByID(uint(id))
			if err != nil {
				log.Print(err)
			}
		}
	}

	username := ctx.Param("username")
	user, err := usecase.UserInstance.GetByUsername(username)
	if err == nil {
		if user.Show.Profile == -1 && (issuer == nil || issuer.ID != user.ID) {
			ctx.JSON(http.StatusForbidden, gin.H{"message": "User disabled viewing their profile"})
			return
		}

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

func GetFollowersHandler(ctx *gin.Context) {
	var issuer *model.User

	idFromJWT, exists := ctx.Get("id_from_jwt")
	if exists {
		log.Print(idFromJWT)
		id, err := strconv.Atoi(fmt.Sprintf("%v", idFromJWT))
		if err != nil {
			log.Print(err)
		} else {
			issuer, err = usecase.UserInstance.GetByID(uint(id))
			if err != nil {
				log.Print(err)
			}
		}
	}

	username := ctx.Param("username")
	user, err := usecase.UserInstance.GetByUsername(username)
	if err != nil {
		log.Print(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error. Please, try again later."})
		return
	}

	if user.Show.Followers == -1 && (issuer == nil || issuer.ID != user.ID) {
		ctx.JSON(http.StatusForbidden, gin.H{"message": "User disabled viewing whom are following them"})
		return
	}

	followers, err := usecase.FollowInstance.GetFollowers(user.ID)
	if err != nil {
		log.Print(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error. Please, try again later."})
		return
	}

	ctx.JSON(http.StatusOK, followers)
}

func GetFollowingHandler(ctx *gin.Context) {
	var issuer *model.User

	idFromJWT, exists := ctx.Get("id_from_jwt")
	if exists {
		log.Print(idFromJWT)
		id, err := strconv.Atoi(fmt.Sprintf("%v", idFromJWT))
		if err != nil {
			log.Print(err)
		} else {
			issuer, err = usecase.UserInstance.GetByID(uint(id))
			if err != nil {
				log.Print(err)
			}
		}
	}

	username := ctx.Param("username")
	user, err := usecase.UserInstance.GetByUsername(username)
	if err != nil {
		log.Print(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error. Please, try again later."})
		return
	}

	if user.Show.Following == -1 && (issuer == nil || issuer.ID != user.ID) {
		ctx.JSON(http.StatusForbidden, gin.H{"message": "User disabled viewing whom they are following"})
		return
	}

	following, err := usecase.FollowInstance.GetFollowing(user.ID)
	if err != nil {
		log.Print(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error. Please, try again later."})
		return
	}

	ctx.JSON(http.StatusOK, following)
}
