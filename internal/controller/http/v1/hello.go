package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func newWorldRoutes(handler *gin.RouterGroup) {
	h := handler.Group("/world")
	{
		h.GET("/hello", HelloHandler)
	}
}

func HelloHandler(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"message": "Hello, world!"})
}
