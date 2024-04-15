package utils

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func LoadDB() {
	// NOTE: It's assuming the .env is already loaded OR env variables are already defined in system
	var err error

	postgres_host := os.Getenv("POSTGRES_HOST")
	postgres_db := os.Getenv("POSTGRES_DB")
	postgres_user := os.Getenv("POSTGRES_USER")
	postgres_password := os.Getenv("POSTGRES_PASSWORD")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=5432 sslmode=disable", postgres_host, postgres_user, postgres_password, postgres_db)
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal("Error connecting to database:", err)
	}

	initValidator()
}
