package middleware

import (
	"net/http"

	"github.com/devkcud/arkhon-foundation/arkhon-api/pkg/utils"
	"github.com/devkcud/arkhon-foundation/arkhon-api/translations"
	"github.com/gin-gonic/gin"
)

func HasPermissions(permissions ...string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		dict := translations.GetTranslation(ctx)

		list, exists := ctx.Get("permissions")
		if !exists {
			ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": dict.Unauthorized})
			return
		}

		if !utils.HasPermissions(list.([]string), permissions...) {
			ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": dict.Unauthorized})
			return
		}

		ctx.Next()
	}
}
