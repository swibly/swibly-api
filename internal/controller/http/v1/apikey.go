package v1

import (
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/devkcud/arkhon-foundation/arkhon-api/internal/model/dto"
	"github.com/devkcud/arkhon-foundation/arkhon-api/internal/service"
	"github.com/devkcud/arkhon-foundation/arkhon-api/pkg/middleware"
	"github.com/devkcud/arkhon-foundation/arkhon-api/pkg/utils"
	"github.com/devkcud/arkhon-foundation/arkhon-api/translations"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func newAPIKeyRoutes(handler *gin.RouterGroup) {
	h := handler.Group("/key")
	h.Use(middleware.APIKeyHasEnabledKeyManage)
	{
		h.GET("/all", GetAllAPIKeys)
		h.GET("/mine", middleware.AuthMiddleware, GetMyAPIKeys)
		h.POST("/create", middleware.OptionalAuthMiddleware, CreateAPIKey)
	}

	specific := h.Group("/:key")
	specific.Use(func(ctx *gin.Context) {
		dict := translations.GetTranslation(ctx)

		key, err := service.APIKey.Find(ctx.Param("key"))
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": dict.NoAPIKeyFound})
				return
			}

			log.Print(err)
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": dict.InternalServerError})
			return
		}

		ctx.Set("api_key_lookup", key)
		ctx.Next()
	})
	{
		specific.GET("", GetAPIKeyInfo)
		specific.DELETE("", DestroyAPIKey)
		specific.PATCH("", UpdateAPIKey)
	}
}

func GetAllAPIKeys(ctx *gin.Context) {
	dict := translations.GetTranslation(ctx)

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

	keys, err := service.APIKey.FindAll(page, perpage)
	if err != nil {
		log.Print(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": dict.InternalServerError})
		return
	}

	ctx.JSON(http.StatusOK, keys)
}

func GetMyAPIKeys(ctx *gin.Context) {
	dict := translations.GetTranslation(ctx)

	issuer := ctx.Keys["auth_user"].(*dto.ProfileSearch)

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

	keys, err := service.APIKey.FindByOwnerUsername(issuer.Username, page, perpage)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": dict.InternalServerError})
		return
	}

	ctx.JSON(http.StatusOK, keys)
}

func CreateAPIKey(ctx *gin.Context) {
	dict := translations.GetTranslation(ctx)

	var issuerUsername string = ""
	if u, exists := ctx.Get("auth_user"); exists {
		issuerUsername = u.(*dto.ProfileSearch).Username
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
	dict := translations.GetTranslation(ctx)

	key, err := service.APIKey.Find(ctx.Param("key"))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": dict.NoAPIKeyFound})
			return
		}

		log.Print(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": dict.InternalServerError})
		return
	}

	ctx.JSON(http.StatusOK, key)
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

	if body.OwnerUsername != "" {
		if _, err := service.User.GetByUsername(body.OwnerUsername); err != nil {
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
