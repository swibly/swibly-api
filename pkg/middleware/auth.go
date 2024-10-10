package middleware

import (
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/swibly/swibly-api/internal/service"
	"github.com/swibly/swibly-api/pkg/utils"
	"github.com/swibly/swibly-api/translations"
	"github.com/gin-gonic/gin"
)

func Auth(ctx *gin.Context) {
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

func OptionalAuth(ctx *gin.Context) {
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
