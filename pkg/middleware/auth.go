package middleware

import (
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/devkcud/arkhon-foundation/arkhon-api/internal/service"
	"github.com/devkcud/arkhon-foundation/arkhon-api/pkg/utils"
	"github.com/devkcud/arkhon-foundation/arkhon-api/translations"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware(ctx *gin.Context) {
	dict := translations.GetTranslation(ctx)

	tokenString := strings.TrimPrefix(ctx.GetHeader("Authorization"), "Bearer ")

	if tokenString == "" {
		log.Print("Couldn't find JWT token in header")

		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": dict.Unauthorized})
		return
	}

	claims, err := utils.GetClaimsJWT(tokenString)
	if err != nil {
		log.Print(err)

		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": dict.Unauthorized})
		return
	}

	id, err := strconv.Atoi(claims.Subject)
	if err != nil || id < 0 {
		log.Print(err)

		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": dict.Unauthorized})
		return
	}

	pf, err := service.User.GetByID(uint(id))
	if pf == nil && err != nil {
		log.Print(err)

		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": dict.Unauthorized})
		return
	}

	ctx.Set("auth_user", pf)
	ctx.Next()
}

func OptionalAuthMiddleware(ctx *gin.Context) {
	dict := translations.GetTranslation(ctx)
	tokenString := strings.TrimPrefix(ctx.GetHeader("Authorization"), "Bearer ")

	if tokenString == "" {
		ctx.Next()
		return
	}

	claims, err := utils.GetClaimsJWT(tokenString)

	if err != nil {
		log.Print(err)
		ctx.Next()
		return
	}

	id, err := strconv.Atoi(claims.Subject)
	if err != nil || id < 0 {
		log.Print(err)

		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": dict.Unauthorized})
		return
	}

	pf, err := service.User.GetByID(uint(id))
	if pf == nil && err != nil {
		log.Print(err)

		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": dict.Unauthorized})
		return
	}

	ctx.Set("auth_user", pf)
	ctx.Next()
}
