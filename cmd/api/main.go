package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	gin.SetMode(gin.ReleaseMode)

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	createRouter().Run(fmt.Sprintf(":%s", port))
}

func createRouter() *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger())

	v1 := router.Group("/v1")
	{
		v1.GET("/hello", func(ctx *gin.Context) {
			ctx.String(http.StatusOK, "Hello, world!")
		})

		v1.GET("/hello/:name", func(ctx *gin.Context) {
			ctx.String(http.StatusOK, fmt.Sprintf("Hello, %s!", ctx.Param("name")))
		})
	}

	return router
}
