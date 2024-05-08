package v1

import (
	"github.com/devkcud/arkhon-foundation/arkhon-api/internal/service/usecase"
	"github.com/gin-gonic/gin"
)

func NewRouter(handler *gin.Engine) {
	g := handler.Group("/v1")
	{
		g.Use(func(ctx *gin.Context) {
			ctx.Set("uc", usecase.NewUserUseCase())
			ctx.Next()
		})

		newAuthRoutes(g)
		newProfileRoutes(g)
	}
}
