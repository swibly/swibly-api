package db

import (
	"fmt"
	"log"

	"github.com/devkcud/arkhon-foundation/arkhon-api/config"
	"github.com/devkcud/arkhon-foundation/arkhon-api/internal/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var Postgres *gorm.DB

func Load() {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=5432 sslmode=%s", config.Postgres.Host, config.Postgres.User, config.Postgres.Password, config.Postgres.DB, config.Postgres.SSLMode)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	// NOTE: Currently, we are handling every error case when doing DB operations. There is no need for an extra logger
	db.Logger = logger.Discard

	Postgres = db

	log.Print("Loaded Database")

	if err := Postgres.AutoMigrate(&model.User{}); err != nil {
		log.Fatal(err)
	}

	log.Print("Loaded migrations")
}
