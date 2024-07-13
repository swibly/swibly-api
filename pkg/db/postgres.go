package db

import (
	"fmt"
	"log"
	"reflect"
	"strings"

	"github.com/devkcud/arkhon-foundation/arkhon-api/config"
	"github.com/devkcud/arkhon-foundation/arkhon-api/internal/model"
	"github.com/devkcud/arkhon-foundation/arkhon-api/pkg/language"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/logger"
)

var Postgres *gorm.DB

func typeCheckAndCreate(db *gorm.DB, typeName string, values []string) error {
	var typeExists bool
	err := db.Raw("SELECT EXISTS (SELECT 1 FROM pg_type WHERE typname =?)", typeName).Scan(&typeExists).Error
	if err != nil {
		return fmt.Errorf("error checking type existence: %w", err)
	}

	if !typeExists {
		err := db.Exec(fmt.Sprintf("CREATE TYPE %s AS ENUM ('%s')", typeName, strings.Join(values, "', '"))).Error
		if err != nil {
			return fmt.Errorf("error creating type: %w", err)
		}
	}

	return nil
}

func Load() {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s", config.Postgres.Host, config.Postgres.User, config.Postgres.Password, config.Postgres.DB, config.Postgres.Port, config.Postgres.SSLMode)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{PrepareStmt: false})
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	// NOTE: Currently, we are handling every error case when doing DB operations. There is no need for an extra logger
	db.Logger = logger.Discard

	Postgres = db

	log.Print("Loaded Database")

	if err := typeCheckAndCreate(Postgres, "enum_language", language.ArrayString); err != nil {
		log.Fatal(err)
	}

	log.Print("Loaded migrations")

	if err := Postgres.AutoMigrate(
		&model.APIKey{},
		&model.User{},
		&model.Follower{},
		&model.Permission{},
		&model.UserPermission{},
	); err != nil {
		log.Fatal(err)
	}

	var permissions []model.Permission
	v := reflect.ValueOf(config.Permissions)
	for i := 0; i < v.NumField(); i++ {
		permissions = append(permissions, model.Permission{Name: v.Field(i).String()})
	}

	if err := db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "name"}},
		DoUpdates: clause.AssignmentColumns([]string{"name"}),
	}).Create(&permissions).Error; err != nil {
		log.Println(err)
	}
}
