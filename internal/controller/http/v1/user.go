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

func newUserRoutes(handler *gin.RouterGroup) {
	h := handler.Group("/user")
	{
		h.GET("/:username/profile", middleware.OptionalAuthMiddleware, GetProfileHandler)
		h.GET("/:username/followers", middleware.OptionalAuthMiddleware, GetFollowersHandler)
		h.GET("/:username/following", middleware.OptionalAuthMiddleware, GetFollowingHandler)

		h.POST("/:username/follow", middleware.AuthMiddleware, FollowUserHandler)
		h.POST("/:username/unfollow", middleware.AuthMiddleware, UnfollowUserHandler)
	}
}

func GetProfileHandler(ctx *gin.Context) {
	var issuer *model.User

	idFromJWT, exists := ctx.Get("id_from_jwt")
	if exists {
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
			ctx.JSON(http.StatusForbidden, gin.H{"error": "User disabled viewing their profile"})
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

	if user.Show.Profile == -1 && (issuer == nil || issuer.ID != user.ID) {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "User disabled viewing their profile"})
		return
	}

	if user.Show.Followers == -1 && (issuer == nil || issuer.ID != user.ID) {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "User disabled viewing whom are following them"})
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

	if user.Show.Profile == -1 && (issuer == nil || issuer.ID != user.ID) {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "User disabled viewing their profile"})
		return
	}

	if user.Show.Following == -1 && (issuer == nil || issuer.ID != user.ID) {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "User disabled viewing whom they are following"})
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

func FollowUserHandler(ctx *gin.Context) {
	idFromJWT, _ := ctx.Get("id_from_jwt")
	id, _ := strconv.Atoi(fmt.Sprintf("%v", idFromJWT))
	issuer, _ := usecase.UserInstance.GetByID(uint(id))

	receiver, err := usecase.UserInstance.GetByUsername(ctx.Param("username"))
	if err != nil {
		log.Print(err)
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user"})
		return
	}

	if issuer.ID == receiver.ID {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Users cannot follow or unfollow themselves"})
		return
	}

	if exists, err := usecase.FollowInstance.Exists(receiver.ID, issuer.ID); err != nil {
		log.Print(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error. Please, try again later."})
		return
	} else if exists {
		ctx.JSON(http.StatusConflict, gin.H{"error": fmt.Sprintf("Already following %s", receiver.Username)})
		return
	}

	if err := usecase.FollowInstance.FollowUser(receiver.ID, issuer.ID); err != nil {
		log.Print(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error. Please, try again later."})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("Started following %s", receiver.Username)})
}

func UnfollowUserHandler(ctx *gin.Context) {
	idFromJWT, _ := ctx.Get("id_from_jwt")
	id, _ := strconv.Atoi(fmt.Sprintf("%v", idFromJWT))
	issuer, _ := usecase.UserInstance.GetByID(uint(id))

	receiver, err := usecase.UserInstance.GetByUsername(ctx.Param("username"))
	if err != nil {
		log.Print(err)
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user"})
		return
	}

	if issuer.ID == receiver.ID {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Users cannot follow or unfollow themselves"})
		return
	}

	if exists, err := usecase.FollowInstance.Exists(receiver.ID, issuer.ID); err != nil {
		log.Print(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error. Please, try again later."})
		return
	} else if !exists {
		ctx.JSON(http.StatusConflict, gin.H{"error": fmt.Sprintf("Not following %s", receiver.Username)})
		return
	}

	if err := usecase.FollowInstance.UnfollowUser(receiver.ID, issuer.ID); err != nil {
		log.Print(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error. Please, try again later."})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("Stopped following %s", receiver.Username)})
}
