package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/devkcud/arkhon-foundation/arkhon-api/config"
	v1 "github.com/devkcud/arkhon-foundation/arkhon-api/internal/controller/http/v1"
	"github.com/devkcud/arkhon-foundation/arkhon-api/internal/service"
	"github.com/devkcud/arkhon-foundation/arkhon-api/pkg/db"
	"github.com/gin-gonic/gin"
)

func main() {
	config.Parse()
	db.Load()

	service.Init()

	gin.SetMode(config.Router.GinMode)

	router := gin.New()
	router.Use(
		gin.Logger(),
		gin.Recovery(),
		func(ctx *gin.Context) {
			key, err := service.APIKey.Find(ctx.GetHeader("X-API-KEY"))
			if err != nil {
				ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid API key"})
				return
			}

			ctx.Set("api_key", key)
			ctx.Next()
		},
	)

	router.GET("/healthz", func(ctx *gin.Context) {
		ctx.Writer.WriteString("Hello, world!")
	})

	v1.NewRouter(router)

	// NOTE: Prioritize the PORT env variable, as some web services may set it
	port := os.Getenv("PORT")

	if port == "" {
		log.Printf("PORT env variable not found, using default: %d", config.Router.Port)
		port = fmt.Sprint(config.Router.Port)
	}

	log.Printf("Using port %s", port)

	go func() {
		log.Print("Starting API")

		if err := router.Run(fmt.Sprintf("%s:%s", config.Router.Address, port)); err != nil {
			log.Fatal(err)
		}
	}()

	// Graceful Shutdown
	exit := make(chan os.Signal, 1)
	signal.Notify(exit, os.Interrupt, syscall.SIGTERM)
	<-exit // Keep process alive

	log.Print("Server stopped. Graceful Shutdown (CTRL+C)")
}
