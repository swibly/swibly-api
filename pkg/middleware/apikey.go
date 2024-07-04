package middleware

import (
	"fmt"
	"net/http"

	"github.com/devkcud/arkhon-foundation/arkhon-api/internal/model"
	"github.com/devkcud/arkhon-foundation/arkhon-api/internal/service"
	"github.com/gin-gonic/gin"
)

func GetAPIKey(ctx *gin.Context) {
	key, err := service.APIKey.Find(ctx.GetHeader("X-API-KEY"))
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid API key"})
		return
	}

	ctx.Set("api_key", key)
	ctx.Next()
}

func apiKeyHas(ctx *gin.Context, b int, field string) {
	if b == -1 {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": fmt.Sprintf("This API key doesn't have the permission to handle %s", field)})
		return
	}

	ctx.Next()
}

func APIKeyHasEnabledKeyManage(ctx *gin.Context) {
	key := ctx.Keys["api_key"].(*model.APIKey)
	apiKeyHas(ctx, key.EnabledKeyManage, "api key manipulation")
}

func APIKeyHasEnabledAuth(ctx *gin.Context) {
	key := ctx.Keys["api_key"].(*model.APIKey)
	apiKeyHas(ctx, key.EnabledAuth, "auth")
}

func APIKeyHasEnabledSearch(ctx *gin.Context) {
	key := ctx.Keys["api_key"].(*model.APIKey)
	apiKeyHas(ctx, key.EnabledSearch, "searches")
}

func APIKeyHasEnabledUserFetch(ctx *gin.Context) {
	key := ctx.Keys["api_key"].(*model.APIKey)
	apiKeyHas(ctx, key.EnabledUserFetch, "user fetch")
}

func APIKeyHasEnabledUserActions(ctx *gin.Context) {
	key := ctx.Keys["api_key"].(*model.APIKey)
	apiKeyHas(ctx, key.EnabledUserActions, "user actions")
}
