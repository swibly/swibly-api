package v1

import (
	"github.com/gin-gonic/gin"
)

func NewRouter(handler *gin.Engine) {
	g := handler.Group("/v1")
	{
		newWorldRoutes(g)
	}
}
