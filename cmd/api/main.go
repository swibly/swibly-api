package main

import (
	"fmt"
	"os"

	"github.com/devkcud/arkhon-foundation/arkhon-api/config"
	v1 "github.com/devkcud/arkhon-foundation/arkhon-api/internal/controller/http/v1"
	"github.com/gin-gonic/gin"
)

func main() {
	config.Parse()

	gin.SetMode(config.Router.GinMode)

	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	v1.NewRouter(router)

	// NOTE: Prioritize the PORT env variable, as some web services may set it
	port := os.Getenv("PORT")

	if port == "" {
		port = fmt.Sprint(config.Router.Port)
	}

	router.Run(fmt.Sprintf("%s:%s", config.Router.Address, port))
}
