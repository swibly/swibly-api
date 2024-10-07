package middleware

import (
	"errors"
	"log"
	"net/http"

	"github.com/swibly/swibly-api/internal/service"
	"github.com/swibly/swibly-api/translations"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Why? It solves a very specific problem where responses should not be cached,
// ensuring that each request receives fresh data based on its parameters.
// Caused due to cloud caching headers or smt like that idk
//
// HACK: Please, future me, fix this when you got time, it's kinda hacky
func DisableCache(ctx *gin.Context) {
	ctx.Header("Cache-Control", "no-store, no-cache, must-revalidate, proxy-revalidate, max-age=0")
	ctx.Header("Pragma", "no-cache")
	ctx.Header("Expires", "1")
	ctx.Next()
}

func UserLookup(ctx *gin.Context) {
	dict := translations.GetTranslation(ctx)

	username := ctx.Param("username")
	user, err := service.User.GetByUsername(username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": dict.UserNotFound})
			return
		}

		log.Print(err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": dict.InternalServerError})
		return
	}

	ctx.Set("user_lookup", user)
	ctx.Next()
}
