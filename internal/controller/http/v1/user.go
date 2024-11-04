package v1

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/swibly/swibly-api/config"
	"github.com/swibly/swibly-api/internal/model/dto"
	"github.com/swibly/swibly-api/internal/service"
	"github.com/swibly/swibly-api/pkg/middleware"
	"github.com/swibly/swibly-api/pkg/notification"
	"github.com/swibly/swibly-api/pkg/utils"
	"github.com/swibly/swibly-api/translations"
)

func newUserRoutes(handler *gin.RouterGroup) {
	h := handler.Group("/user/:username", middleware.APIKeyHasEnabledUserFetch, middleware.Auth, middleware.UserLookup)
	{
		h.GET("/profile", middleware.UserPrivacy(dto.UserShow{Profile: true}), GetProfileHandler)
		h.GET("/followers", middleware.UserPrivacy(dto.UserShow{Followers: true}), GetFollowersHandler)
		h.GET("/following", middleware.UserPrivacy(dto.UserShow{Following: true}), GetFollowingHandler)
	}

	actions := h.Group("", middleware.APIKeyHasEnabledUserActions)
	{
		actions.POST("/follow", FollowUserHandler)
		actions.POST("/unfollow", UnfollowUserHandler)

		actions.GET("/amifollowing", IsFollowingHandler)
	}
}

func GetProfileHandler(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, ctx.Keys["user_lookup"])
}

func GetFollowersHandler(ctx *gin.Context) {
	dict := translations.GetTranslation(ctx)

	user := ctx.Keys["user_lookup"].(*dto.UserProfile)

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

	user := ctx.Keys["user_lookup"].(*dto.UserProfile)

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
	receiver := ctx.Keys["user_lookup"].(*dto.UserProfile)

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

	service.CreateNotification(dto.CreateNotification{
		Title:    dict.CategoryFollowers,
		Message:  fmt.Sprintf(dict.NotificationUserFollowedYou, issuer.FirstName+issuer.LastName),
		Type:     notification.Information,
		Redirect: utils.ToPtr(fmt.Sprintf(config.Redirects.Profile, issuer.Username)),
	}, receiver.ID)

	ctx.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf(dict.UserFollowingStarted, receiver.Username)})
}

func UnfollowUserHandler(ctx *gin.Context) {
	dict := translations.GetTranslation(ctx)

	issuer := ctx.Keys["auth_user"].(*dto.UserProfile)
	receiver := ctx.Keys["user_lookup"].(*dto.UserProfile)

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
	receiver := ctx.Keys["user_lookup"].(*dto.UserProfile)

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
