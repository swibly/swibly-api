package db

import (
	"fmt"
	"log"

	"github.com/devkcud/arkhon-foundation/arkhon-api/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var Postgres *gorm.DB

func Load() {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=5432 sslmode=disable", config.Postgres.Host, config.Postgres.User, config.Postgres.Password, config.Postgres.DB)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	Postgres = db
}
