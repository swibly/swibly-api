package main

import (
	"fmt"
	"os"

	"github.com/devkcud/arkhon-foundation/arkhon-api/internal/controller/auth"
	"github.com/devkcud/arkhon-foundation/arkhon-api/internal/model"
	"github.com/devkcud/arkhon-foundation/arkhon-api/internal/utils"
	"github.com/gin-gonic/gin"
)

func main() {
	// Load funcs cannot run in goroutines due to dependency of the program on their global variables
	// e.g.: utils.DB and env variables
	utils.LoadEnv()
	utils.LoadDB()

	// Using a goroutine for AutoMigrate prevents thread blocking,
	// allowing the rest of the application to run smoothly.
	go utils.DB.AutoMigrate(&model.User{})

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
