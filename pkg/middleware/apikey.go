package middleware

import (
	"fmt"
	"net/http"

	"github.com/devkcud/arkhon-foundation/arkhon-api/internal/model"
	"github.com/gin-gonic/gin"
)

func apiKeyHas(ctx *gin.Context, b bool, field string) {
	if !b {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": fmt.Sprintf("This API key doesn't have the permission to handle %s", field)})
		return
	}

	ctx.Next()
}

func APIKeyHasEnabledKeyGen(ctx *gin.Context) {
	key := ctx.Keys["api_key"].(model.APIKey)
	apiKeyHas(ctx, key.EnabledKeyGen, "key generation")
}

func APIKeyHasEnabledAuth(ctx *gin.Context) {
	key := ctx.Keys["api_key"].(model.APIKey)
	apiKeyHas(ctx, key.EnabledAuth, "auth")
}

func APIKeyHasEnabledSearch(ctx *gin.Context) {
	key := ctx.Keys["api_key"].(model.APIKey)
	apiKeyHas(ctx, key.EnabledSearch, "searches")
}

func APIKeyHasEnabledUserFetch(ctx *gin.Context) {
	key := ctx.Keys["api_key"].(model.APIKey)
	apiKeyHas(ctx, key.EnabledUserFetch, "user fetch")
}

func APIKeyHasEnabledUserActions(ctx *gin.Context) {
	key := ctx.Keys["api_key"].(model.APIKey)
	apiKeyHas(ctx, key.EnabledUserActions, "user actions")
}
