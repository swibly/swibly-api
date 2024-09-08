package v1

import (
	"log"
	"net/http"
	"strconv"

	"github.com/devkcud/arkhon-foundation/arkhon-api/internal/model/dto"
	"github.com/devkcud/arkhon-foundation/arkhon-api/internal/service"
	"github.com/devkcud/arkhon-foundation/arkhon-api/pkg/middleware"
	"github.com/devkcud/arkhon-foundation/arkhon-api/pkg/utils"
	"github.com/devkcud/arkhon-foundation/arkhon-api/translations"
	"github.com/gin-gonic/gin"
)

func newAPIKeyRoutes(handler *gin.RouterGroup) {
	h := handler.Group("/keys", middleware.APIKeyHasEnabledKeyManage)
	{
		h.GET("", middleware.OptionalAuth, GetAllAPIKeys)
		h.POST("/create", middleware.OptionalAuth, CreateAPIKey)
	}

	specific := h.Group("/:key", middleware.Auth, middleware.APIKeyLookup)
	{
		specific.GET("", GetAPIKeyInfo)
		specific.DELETE("", DestroyAPIKey)
		specific.PATCH("", UpdateAPIKey)
		specific.POST("regenerate", RegenerateAPIKey)
	}
}

func GetAllAPIKeys(ctx *gin.Context) {
	dict := translations.GetTranslation(ctx)

	var (
		page    int  = 1
		perpage int  = 10
		own     bool = false
	)

	if i, e := strconv.Atoi(ctx.Query("page")); e == nil && ctx.Query("page") != "" {
		page = i
	}

	if i, e := strconv.Atoi(ctx.Query("perpage")); e == nil && ctx.Query("perpage") != "" {
		perpage = i
	}

	if _, exists := ctx.GetQuery("own"); exists {
		own = true
	}

	if own {
		issuer, exists := ctx.Get("auth_user")
		if !exists {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": dict.Unauthorized})
			return
		}

		keys, err := service.APIKey.GetByOwner(issuer.(*dto.UserProfile).Username, page, perpage)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": dict.InternalServerError})
			return
		}

		ctx.JSON(http.StatusOK, keys)
		return
	}

	keys, err := service.APIKey.GetAll(page, perpage)
	if err != nil {
		log.Print(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": dict.InternalServerError})
		return
	}

	ctx.JSON(http.StatusOK, keys)
}

func CreateAPIKey(ctx *gin.Context) {
	dict := translations.GetTranslation(ctx)

	var issuerUsername string = ""
	if u, exists := ctx.Get("auth_user"); exists {
		issuerUsername = u.(*dto.UserProfile).Username
	}

	maxUsage, err := strconv.ParseUint(ctx.Query("maxusage"), 10, 64)
	if err != nil {
		maxUsage = 0
	}

	newKey, err := service.APIKey.Create(issuerUsername, uint(maxUsage))
	if err != nil {
		log.Printf("Error generating new API key: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": dict.InternalServerError})
		return
	}

	ctx.JSON(http.StatusCreated, newKey)
}

func GetAPIKeyInfo(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, ctx.Keys["api_key_lookup"].(*dto.ReadAPIKey))
}

func DestroyAPIKey(ctx *gin.Context) {
	dict := translations.GetTranslation(ctx)

	if err := service.APIKey.Delete(ctx.Keys["api_key_lookup"].(*dto.ReadAPIKey).Key); err != nil {
		log.Printf("Error destroying API key: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": dict.InternalServerError})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": dict.APIKeyDestroyed})
}

func UpdateAPIKey(ctx *gin.Context) {
	dict := translations.GetTranslation(ctx)

	key := ctx.Keys["api_key_lookup"].(*dto.ReadAPIKey)

	var body dto.UpdateAPIKey
	if err := ctx.BindJSON(&body); err != nil {
		log.Print(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": dict.InvalidBody})
		return
	}

	if errs := utils.ValidateStruct(&body); errs != nil {
		err := utils.ValidateErrorMessage(ctx, errs[0])

		log.Print(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": gin.H{err.Param: err.Message}})
		return
	}

	if body.Owner != "" {
		if _, err := service.User.GetByUsername(body.Owner); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": dict.UserNotFound})
			return
		}
	}

	if err := service.APIKey.Update(key.Key, &body); err != nil {
		log.Print(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": dict.InternalServerError})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": dict.APIKeyUpdated})
}

func RegenerateAPIKey(ctx *gin.Context) {
	dict := translations.GetTranslation(ctx)

	oldKey := ctx.Keys["api_key_lookup"].(*dto.ReadAPIKey)

	err := service.APIKey.Regenerate(oldKey.Key)
	if err != nil {
		log.Print(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": dict.InternalServerError})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": dict.APIKeyUpdated})
}
