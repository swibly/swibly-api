package middleware

import (
	"fmt"
	"log"
	"net/http"

	"github.com/devkcud/arkhon-foundation/arkhon-api/internal/model"
	"github.com/devkcud/arkhon-foundation/arkhon-api/internal/service"
	"github.com/devkcud/arkhon-foundation/arkhon-api/translations"
	"github.com/gin-gonic/gin"
)

func GetAPIKey(ctx *gin.Context) {
	dict := translations.GetLang(ctx)

	key, err := service.APIKey.Find(ctx.GetHeader("X-API-KEY"))
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": dict.InvalidAPIKey})
		return
	}

	ctx.Set("api_key", key)
	ctx.Next()
}

func apiKeyHas(ctx *gin.Context, b int, permission string) {
	dict := translations.GetLang(ctx)

	key := ctx.Keys["api_key"].(*model.APIKey)

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
	key := ctx.Keys["api_key"].(*model.APIKey)
	apiKeyHas(ctx, key.EnabledKeyManage, "manage.api")
}

func APIKeyHasEnabledAuth(ctx *gin.Context) {
	key := ctx.Keys["api_key"].(*model.APIKey)
	apiKeyHas(ctx, key.EnabledAuth, "manage.auth")
}

func APIKeyHasEnabledSearch(ctx *gin.Context) {
	key := ctx.Keys["api_key"].(*model.APIKey)
	apiKeyHas(ctx, key.EnabledSearch, "query.search")
}

func APIKeyHasEnabledUserFetch(ctx *gin.Context) {
	key := ctx.Keys["api_key"].(*model.APIKey)
	apiKeyHas(ctx, key.EnabledUserFetch, "query.user")
}

func APIKeyHasEnabledUserActions(ctx *gin.Context) {
	key := ctx.Keys["api_key"].(*model.APIKey)
	apiKeyHas(ctx, key.EnabledUserActions, "actions")
}
