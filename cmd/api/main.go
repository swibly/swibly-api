package main

import (
	"fmt"
	"os"

	"github.com/devkcud/arkhon-foundation/arkhon-api/internal/controller/auth"
	"github.com/devkcud/arkhon-foundation/arkhon-api/internal/utils"
	"github.com/gin-gonic/gin"
)

func main() {
	utils.LoadEnv()
	utils.LoadDB()

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
		user.POST("/register", auth.RegisterHandler)
		user.POST("/login", auth.LoginHandler)
	}

	return router
}
