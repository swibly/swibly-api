package v1

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/swibly/swibly-api/config"
	"github.com/swibly/swibly-api/internal/model/dto"
	"github.com/swibly/swibly-api/internal/service"
	"github.com/swibly/swibly-api/pkg/middleware"
	"github.com/swibly/swibly-api/pkg/utils"
	"github.com/swibly/swibly-api/translations"
	"gorm.io/gorm"
)

func newUserRoutes(handler *gin.RouterGroup) {
	h := handler.Group("/user/:username", middleware.APIKeyHasEnabledUserFetch, middleware.Auth)
	{
		h.GET("/profile", GetProfileHandler)
		h.GET("/followers", GetFollowersHandler)
		h.GET("/following", GetFollowingHandler)
	}

	actions := h.Group("", middleware.APIKeyHasEnabledUserActions)
	{
		actions.POST("/follow", FollowUserHandler)
		actions.POST("/unfollow", UnfollowUserHandler)

		actions.GET("/amifollowing", IsFollowingHandler)
	}
}

func GetProfileHandler(ctx *gin.Context) {
	dict := translations.GetTranslation(ctx)

	var issuer *dto.UserProfile = nil
	if p, exists := ctx.Get("auth_user"); exists {
		issuer = p.(*dto.UserProfile)
	}

	username := ctx.Param("username")
	user, err := service.User.GetByUsername(username)
	if err == nil {
		if !utils.HasPermissions(user.Permissions, config.Permissions.ManageUser) {
			if user.Show.Profile == false && (issuer == nil || issuer.ID != user.ID) {
				ctx.JSON(http.StatusForbidden, gin.H{"error": dict.UserDisabledProfile})
				return
			}
		}

		ctx.JSON(http.StatusOK, user)

		return
	}

	log.Print(err)

	if err == gorm.ErrRecordNotFound {
		ctx.JSON(http.StatusNotFound, gin.H{"error": dict.UserNotFound})
		return
	}

	ctx.JSON(http.StatusInternalServerError, gin.H{"error": dict.InternalServerError})
}

func GetFollowersHandler(ctx *gin.Context) {
	dict := translations.GetTranslation(ctx)

	var issuer *dto.UserProfile = nil
	if p, exists := ctx.Get("auth_user"); exists {
		issuer = p.(*dto.UserProfile)
	}

	user, err := service.User.GetByUsername(ctx.Param("username"))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": dict.UserNotFound})
			return
		}

		log.Print(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": dict.InternalServerError})
		return
	}

	if !utils.HasPermissions(user.Permissions, config.Permissions.ManageUser) {
		if user.Show.Profile == false && (issuer == nil || issuer.ID != user.ID) {
			ctx.JSON(http.StatusForbidden, gin.H{"error": dict.UserDisabledProfile})
			return
		}

		if user.Show.Followers == false && (issuer == nil || issuer.ID != user.ID) {
			ctx.JSON(http.StatusForbidden, gin.H{"error": dict.UserDisabledFollowers})
			return
		}
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

	pagination, err := service.Follow.GetFollowers(user.ID, page, perpage)
	if err != nil {
		log.Print(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": dict.InternalServerError})
		return
	}

	ctx.JSON(http.StatusOK, pagination)
}

func GetFollowingHandler(ctx *gin.Context) {
	dict := translations.GetTranslation(ctx)

	var issuer *dto.UserProfile = nil
	if p, exists := ctx.Get("auth_user"); exists {
		issuer = p.(*dto.UserProfile)
	}

	user, err := service.User.GetByUsername(ctx.Param("username"))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": dict.UserNotFound})
			return
		}

		log.Print(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": dict.InternalServerError})
		return
	}

	if !utils.HasPermissions(user.Permissions, config.Permissions.ManageUser) {
		if user.Show.Profile == false && (issuer == nil || issuer.ID != user.ID) {
			ctx.JSON(http.StatusForbidden, gin.H{"error": dict.UserDisabledProfile})
			return
		}

		if user.Show.Following == false && (issuer == nil || issuer.ID != user.ID) {
			ctx.JSON(http.StatusForbidden, gin.H{"error": dict.UserDisabledFollowers})
			return
		}
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

	pagination, err := service.Follow.GetFollowing(user.ID, page, perpage)
	if err != nil {
		log.Print(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": dict.InternalServerError})
		return
	}

	ctx.JSON(http.StatusOK, pagination)
}

func FollowUserHandler(ctx *gin.Context) {
	dict := translations.GetTranslation(ctx)

	issuer := ctx.Keys["auth_user"].(*dto.UserProfile)

	receiver, err := service.User.GetByUsername(ctx.Param("username"))
	if err != nil {
		log.Print(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": dict.UserNotFound})
		return
	}

	if issuer.ID == receiver.ID {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": dict.UserErrorFollowItself})
		return
	}

	if exists, err := service.Follow.Exists(receiver.ID, issuer.ID); err != nil {
		log.Print(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": dict.InternalServerError})
		return
	} else if exists {
		ctx.JSON(http.StatusConflict, gin.H{"error": fmt.Sprintf(dict.UserFollowingAlready, receiver.Username)})
		return
	}

	if err := service.Follow.Follow(receiver.ID, issuer.ID); err != nil {
		log.Print(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": dict.InternalServerError})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf(dict.UserFollowingStarted, receiver.Username)})
}

func UnfollowUserHandler(ctx *gin.Context) {
	dict := translations.GetTranslation(ctx)

	issuer := ctx.Keys["auth_user"].(*dto.UserProfile)

	receiver, err := service.User.GetByUsername(ctx.Param("username"))
	if err != nil {
		log.Print(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": dict.UserNotFound})
		return
	}

	if issuer.ID == receiver.ID {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": dict.UserErrorFollowItself})
		return
	}

	if exists, err := service.Follow.Exists(receiver.ID, issuer.ID); err != nil {
		log.Print(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": dict.InternalServerError})
		return
	} else if !exists {
		ctx.JSON(http.StatusConflict, gin.H{"error": fmt.Sprintf(dict.UserFollowingNot, receiver.Username)})
		return
	}

	if err := service.Follow.Unfollow(receiver.ID, issuer.ID); err != nil {
		log.Print(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": dict.InternalServerError})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf(dict.UserFollowingStopped, receiver.Username)})
}

func IsFollowingHandler(ctx *gin.Context) {
	dict := translations.GetTranslation(ctx)

	issuer := ctx.Keys["auth_user"].(*dto.UserProfile)

	receiver, err := service.User.GetByUsername(ctx.Param("username"))
	if err != nil {
		log.Print(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": dict.UserNotFound})
		return
	}

	if issuer.ID == receiver.ID {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": dict.UserErrorFollowItself})
		return
	}

	exists, err := service.Follow.Exists(receiver.ID, issuer.ID)
	if err != nil {
		log.Print(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": dict.InternalServerError})
		return
	}

	switch exists {
	case true:
		ctx.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf(dict.UserFollowingAlready, receiver.Username)})
	case false:
		ctx.JSON(http.StatusNotFound, gin.H{"message": fmt.Sprintf(dict.UserFollowingNot, receiver.Username)})
	}
}
