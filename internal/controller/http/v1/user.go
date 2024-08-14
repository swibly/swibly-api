package v1

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/devkcud/arkhon-foundation/arkhon-api/config"
	"github.com/devkcud/arkhon-foundation/arkhon-api/internal/model/dto"
	"github.com/devkcud/arkhon-foundation/arkhon-api/internal/service"
	"github.com/devkcud/arkhon-foundation/arkhon-api/pkg/middleware"
	"github.com/devkcud/arkhon-foundation/arkhon-api/pkg/utils"
	"github.com/devkcud/arkhon-foundation/arkhon-api/translations"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func newUserRoutes(handler *gin.RouterGroup) {
	h := handler.Group("/user")
	{
		h.GET("/:username/profile", middleware.APIKeyHasEnabledUserFetch, middleware.OptionalAuthMiddleware, middleware.GetPermissionsMiddleware, GetProfileHandler)

		h.GET("/:username/followers", middleware.APIKeyHasEnabledUserFetch, middleware.OptionalAuthMiddleware, middleware.GetPermissionsMiddleware, GetFollowersHandler)
		h.GET("/:username/following", middleware.APIKeyHasEnabledUserFetch, middleware.OptionalAuthMiddleware, middleware.GetPermissionsMiddleware, GetFollowingHandler)

		h.GET("/:username/permissions", middleware.APIKeyHasEnabledUserFetch, GetUserPermissions)

		h.POST("/:username/follow", middleware.APIKeyHasEnabledUserActions, middleware.AuthMiddleware, FollowUserHandler)
		h.POST("/:username/unfollow", middleware.APIKeyHasEnabledUserActions, middleware.AuthMiddleware, UnfollowUserHandler)
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
		if !utils.HasPermissionsByContext(ctx, config.Permissions.ManageUser) {
			if user.Show.Profile == -1 && (issuer == nil || issuer.ID != user.ID) {
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

	username := ctx.Param("username")
	user, err := service.User.GetByUsername(username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": dict.UserNotFound})
			return
		}

		log.Print(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": dict.InternalServerError})
		return
	}

	if !utils.HasPermissionsByContext(ctx, config.Permissions.ManageUser) {
		if user.Show.Profile == -1 && (issuer == nil || issuer.ID != user.ID) {
			ctx.JSON(http.StatusForbidden, gin.H{"error": dict.UserDisabledProfile})
			return
		}

		if user.Show.Followers == -1 && (issuer == nil || issuer.ID != user.ID) {
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

	username := ctx.Param("username")
	user, err := service.User.GetByUsername(username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": dict.UserNotFound})
			return
		}

		log.Print(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": dict.InternalServerError})
		return
	}

	if !utils.HasPermissionsByContext(ctx, config.Permissions.ManageUser) {
		if user.Show.Profile == -1 && (issuer == nil || issuer.ID != user.ID) {
			ctx.JSON(http.StatusForbidden, gin.H{"error": dict.UserDisabledProfile})
			return
		}

		if user.Show.Following == -1 && (issuer == nil || issuer.ID != user.ID) {
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

func GetUserPermissions(ctx *gin.Context) {
	dict := translations.GetTranslation(ctx)

	username := ctx.Param("username")
	user, err := service.User.GetByUsername(username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": dict.UserNotFound})
			return
		}

		log.Print(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": dict.InternalServerError})
		return
	}

	permissions, err := service.Permission.GetPermissions(user.ID)
	if err != nil {
		log.Print(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": dict.InternalServerError})
		return
	}

	list := []string{}

	for _, permission := range permissions {
		list = append(list, permission.Name)
	}

	ctx.JSON(http.StatusOK, list)
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

	if err := service.Follow.FollowUser(receiver.ID, issuer.ID); err != nil {
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

	if err := service.Follow.UnfollowUser(receiver.ID, issuer.ID); err != nil {
		log.Print(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": dict.InternalServerError})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf(dict.UserFollowingStopped, receiver.Username)})
}
