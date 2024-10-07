package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/swibly/swibly-api/config"
	v1 "github.com/swibly/swibly-api/internal/controller/http/v1"
	"github.com/swibly/swibly-api/internal/service"
	"github.com/swibly/swibly-api/pkg/db"
	"github.com/swibly/swibly-api/pkg/middleware"
	"github.com/swibly/swibly-api/translations"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	config.Parse()
	db.Load()

	service.Init()
	translations.Init("./translations")

	gin.SetMode(config.Router.GinMode)

	router := gin.New()
	router.Use(
		cors.New(cors.Config{
			AllowOrigins:     []string{"*"},
			AllowMethods:     []string{"GET", "POST", "PATCH", "DELETE", "PUT", "OPTIONS"},
			AllowHeaders:     []string{"*"},
			ExposeHeaders:    []string{"*"},
			AllowCredentials: true,
		}),
		gin.Logger(),
		gin.Recovery(),
		middleware.DisableCache,
		middleware.GetLanguage,
		middleware.GetAPIKey,
	)

	// Apparently, some testing frameworks and is-this-alive-checking-tools just can't handle GET and prefer OPTIONS instead :/
	router.Any("/healthcare", func(ctx *gin.Context) {
		ctx.Writer.WriteString(ctx.Keys["lang"].(translations.Translation).Hello)
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
