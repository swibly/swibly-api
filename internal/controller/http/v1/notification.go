package v1

import (
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/swibly/swibly-api/internal/model/dto"
	"github.com/swibly/swibly-api/internal/service"
	"github.com/swibly/swibly-api/internal/service/repository"
	"github.com/swibly/swibly-api/pkg/middleware"
	"github.com/swibly/swibly-api/translations"
)

func newNotificationRoutes(handler *gin.RouterGroup) {
	h := handler.Group("/notification")
	h.Use(middleware.APIKeyHasEnabledUserFetch, middleware.Auth)
	{
		h.GET("", GetOwnNotificationsHandler)

		specific := h.Group("/:id")
		{
			specific.POST("/read", PostReadNotificationHandler)
			specific.DELETE("/unread", DeleteUnreadNotificationHandler)
		}
	}
}

func GetOwnNotificationsHandler(ctx *gin.Context) {
	dict := translations.GetTranslation(ctx)

	page := 1
	perPage := 10
	onlyUnread := false

	if i, e := strconv.Atoi(ctx.Query("page")); e == nil && ctx.Query("page") != "" {
		page = i
	}

	if i, e := strconv.Atoi(ctx.Query("perpage")); e == nil && ctx.Query("perpage") != "" {
		perPage = i
	}

	unreadFlag := strings.ToLower(ctx.Query("unread"))
	if unreadFlag == "true" || unreadFlag == "t" || unreadFlag == "1" {
		onlyUnread = true
	}

	issuer := ctx.Keys["auth_user"].(*dto.UserProfile)

	notifications, err := service.Notification.GetForUser(issuer.ID, onlyUnread, page, perPage)
	if err != nil {
		log.Print(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": dict.InternalServerError})
		return
	}

	ctx.JSON(http.StatusOK, notifications)
}

func PostReadNotificationHandler(ctx *gin.Context) {
	dict := translations.GetTranslation(ctx)
	issuer := ctx.Keys["auth_user"].(*dto.UserProfile)

	notificationID, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		log.Print(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": dict.NotificationInvalid})
		return
	}

	if err := service.Notification.MarkAsRead(*issuer, uint(notificationID)); err != nil {
		switch {
		case errors.Is(err, repository.ErrNotificationNotAssigned):
			ctx.JSON(http.StatusForbidden, gin.H{"error": dict.NotificationNotAssigned})
		case errors.Is(err, repository.ErrNotificationAlreadyRead):
			ctx.JSON(http.StatusConflict, gin.H{"error": dict.NotificationAlreadyRead})
		default:
			log.Print(err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": dict.InternalServerError})
		}
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": dict.NotificationMarkedAsRead})
}

func DeleteUnreadNotificationHandler(ctx *gin.Context) {
	dict := translations.GetTranslation(ctx)
	issuer := ctx.Keys["auth_user"].(*dto.UserProfile)

	notificationID, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		log.Print(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": dict.NotificationInvalid})
		return
	}

	if err := service.Notification.MarkAsUnread(*issuer, uint(notificationID)); err != nil {
		switch {
		case errors.Is(err, repository.ErrNotificationNotAssigned):
			ctx.JSON(http.StatusForbidden, gin.H{"error": dict.NotificationNotAssigned})
		case errors.Is(err, repository.ErrNotificationNotRead):
			ctx.JSON(http.StatusConflict, gin.H{"error": dict.NotificationNotRead})
		default:
			log.Print(err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": dict.InternalServerError})
		}
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": dict.NotificationMarkedAsUnread})
}
