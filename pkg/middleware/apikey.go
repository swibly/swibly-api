package middleware

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/devkcud/arkhon-foundation/arkhon-api/config"
	"github.com/devkcud/arkhon-foundation/arkhon-api/internal/model/dto"
	"github.com/devkcud/arkhon-foundation/arkhon-api/internal/service"
	"github.com/devkcud/arkhon-foundation/arkhon-api/pkg/utils"
	"github.com/devkcud/arkhon-foundation/arkhon-api/translations"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func APIKeyLookup(ctx *gin.Context) {
	dict := translations.GetTranslation(ctx)

	key, err := service.APIKey.GetByKey(ctx.Param("key"))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": dict.NoAPIKeyFound})
			return
		}

		log.Print(err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": dict.InternalServerError})
		return
	}

	user := ctx.Keys["auth_user"].(*dto.UserProfile)

	if key.Owner != user.Username && !utils.HasPermissions(user.Permissions, config.Permissions.ManageAPIKey) {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": dict.Unauthorized})
		return
	}

	ctx.Set("api_key_lookup", key)
	ctx.Next()
}

func GetAPIKey(ctx *gin.Context) {
	dict := translations.GetTranslation(ctx)

	if strings.TrimSpace(ctx.GetHeader("X-API-KEY")) == "" {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": dict.InvalidAPIKey})
		return
	}

	key, err := service.APIKey.GetByKey(ctx.GetHeader("X-API-KEY"))
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": dict.InvalidAPIKey})
		return
	}

	ctx.Set("api_key", key)
	ctx.Next()
}

func apiKeyHas(ctx *gin.Context, b int, permission string) {
	dict := translations.GetTranslation(ctx)

	key := ctx.Keys["api_key"].(*dto.ReadAPIKey)

	if key.MaxUsage != 0 && key.TimesUsed >= key.MaxUsage {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": dict.MaximumAPIKey})
		return
	}

	if err := service.APIKey.RegisterUse(key.Key); err != nil {
		log.Print(err)
	}

	if b == -1 {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": fmt.Sprintf(dict.RequirePermissionAPIKey, permission)})
		return
	}

	ctx.Next()
}

func APIKeyHasEnabledKeyManage(ctx *gin.Context) {
	key := ctx.Keys["api_key"].(*dto.ReadAPIKey)
	apiKeyHas(ctx, key.EnabledKeyManage, "manage.api")
}

func APIKeyHasEnabledAuth(ctx *gin.Context) {
	key := ctx.Keys["api_key"].(*dto.ReadAPIKey)
	apiKeyHas(ctx, key.EnabledAuth, "manage.auth")
}

func APIKeyHasEnabledSearch(ctx *gin.Context) {
	key := ctx.Keys["api_key"].(*dto.ReadAPIKey)
	apiKeyHas(ctx, key.EnabledSearch, "query.search")
}

func APIKeyHasEnabledUserFetch(ctx *gin.Context) {
	key := ctx.Keys["api_key"].(*dto.ReadAPIKey)
	apiKeyHas(ctx, key.EnabledUserFetch, "query.user")
}

func APIKeyHasEnabledUserActions(ctx *gin.Context) {
	key := ctx.Keys["api_key"].(*dto.ReadAPIKey)
	apiKeyHas(ctx, key.EnabledUserActions, "actions")
}

func APIKeyHasEnabledProjects(ctx *gin.Context) {
	key := ctx.Keys["api_key"].(*dto.ReadAPIKey)
	apiKeyHas(ctx, key.EnabledProjects, "query.project")
}
