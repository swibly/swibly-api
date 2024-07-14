package middleware

import (
	"slices"
	"strings"

	"github.com/devkcud/arkhon-foundation/arkhon-api/pkg/language"
	"github.com/devkcud/arkhon-foundation/arkhon-api/translations"
	"github.com/gin-gonic/gin"
)

func DetectLanguage(ctx *gin.Context) {
	xlang := strings.ToLower(strings.TrimSpace(ctx.GetHeader("X-Lang")))

	if !slices.Contains(language.ArrayString, xlang) {
		xlang = string(language.PT)
	}

	ctx.Set("lang", translations.Translations[xlang])

	ctx.Next()
}
