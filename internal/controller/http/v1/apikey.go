package v1

import (
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/devkcud/arkhon-foundation/arkhon-api/internal/model"
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
		key, err := service.APIKey.Find(ctx.Param("key"))
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "No API key found."})
				return
			}

			log.Print(err)
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": ctx.Keys["lang"].(translations.Translation).InternalServerError})
			return
		}

		ctx.Set("api_key_lookup", key)
		ctx.Next()
	})
	{
		specific.GET("", GetAPIKeyInfo)
		specific.DELETE("/destroy", DestroyAPIKey)
		specific.PATCH("/update", UpdateAPIKey)
	}
}

func GetAllAPIKeys(ctx *gin.Context) {
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
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": ctx.Keys["lang"].(translations.Translation).InternalServerError})
		return
	}

	ctx.JSON(http.StatusOK, keys)
}

func GetMyAPIKeys(ctx *gin.Context) {
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

	keys, err := service.APIKey.FindByOwnerID(issuer.ID, page, perpage)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": ctx.Keys["lang"].(translations.Translation).InternalServerError})
		return
	}

	ctx.JSON(http.StatusOK, keys)
}

func CreateAPIKey(ctx *gin.Context) {
	var issuerID uint = 0
	if u, exists := ctx.Get("auth_user"); exists {
		issuerID = u.(*dto.ProfileSearch).ID
	}

	maxUsage, err := strconv.ParseUint(ctx.Query("maxusage"), 10, 64)
	if err != nil {
		maxUsage = 0
	}

	newKey, err := service.APIKey.Create(issuerID, uint(maxUsage))
	if err != nil {
		log.Printf("Error generating new API key: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Couldn't generate new key"})
		return
	}

	ctx.JSON(http.StatusCreated, newKey)
}

func GetAPIKeyInfo(ctx *gin.Context) {
	key, err := service.APIKey.Find(ctx.Param("key"))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "No API key found."})
			return
		}

		log.Print(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": ctx.Keys["lang"].(translations.Translation).InternalServerError})
		return
	}

	ctx.JSON(http.StatusOK, key)
}

func DestroyAPIKey(ctx *gin.Context) {
	if err := service.APIKey.Delete(ctx.Keys["api_key_lookup"].(*model.APIKey).Key); err != nil {
		log.Printf("Error destroying API key: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Couldn't destroy key"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Destroyied key"})
}

func UpdateAPIKey(ctx *gin.Context) {
	key := ctx.Keys["api_key_lookup"].(*model.APIKey)

	var body dto.APIKey
	if err := ctx.BindJSON(&body); err != nil {
		log.Print(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": ctx.Keys["lang"].(translations.Translation).InvalidBody})
		return
	}

	if errs := utils.ValidateStruct(&body); errs != nil {
		err := utils.ValidateErrorMessage(errs[0])

		log.Print(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": gin.H{err.Param: err.Message}})
		return
	}

	if err := service.APIKey.Update(key.Key, &body); err != nil {
		log.Print(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": ctx.Keys["lang"].(translations.Translation).InternalServerError})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "API key updated"})
}
