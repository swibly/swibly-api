package middleware

import (
	"slices"
	"strings"

	"github.com/swibly/swibly-api/pkg/language"
	"github.com/swibly/swibly-api/translations"
	"github.com/gin-gonic/gin"
)

func GetLanguage(ctx *gin.Context) {
	xlang := strings.ToLower(strings.TrimSpace(ctx.GetHeader("X-Lang")))

	if !slices.Contains(language.ArrayString, xlang) {
		xlang = string(language.PT)
	}

	ctx.Set("lang", translations.Translations[xlang])

	ctx.Next()
}
