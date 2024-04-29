package middleware

import (
	"log"
	"net/http"
	"strings"

	"github.com/devkcud/arkhon-foundation/arkhon-api/pkg/utils"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware(ctx *gin.Context) {
	tokenString := strings.TrimPrefix(ctx.GetHeader("Authorization"), "Bearer ")

	if tokenString == "" {
		log.Print("Couldn't find JWT token in header")

		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	_, err := utils.GetClaimsJWT(tokenString)

	if err != nil {
		log.Print(err)

		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	ctx.Next()
}
