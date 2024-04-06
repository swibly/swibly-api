package main

import (
	"fmt"
	"os"

	"github.com/devkcud/arkhon-foundation/arkhon-api/internal/controller"
	"github.com/devkcud/arkhon-foundation/arkhon-api/internal/loader"
	"github.com/devkcud/arkhon-foundation/arkhon-api/internal/middleware"
	"github.com/gin-gonic/gin"
)

func main() {
	loader.LoadEnv()
	loader.LoadDB()

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

	// XXX: API v1
	v1 := router.Group("/v1")
	user := v1.Group("/user")
	{
		user.POST("/register", controller.Register)
		user.POST("/login", controller.Login)
		user.POST("/logout", middleware.RequireAuth, controller.Logout)
	}

	return router
}
