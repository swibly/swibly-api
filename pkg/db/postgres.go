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
	createTypeSQL := fmt.Sprintf(`DO $$ BEGIN IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = '%s') THEN CREATE TYPE %s AS ENUM ('%s'); END IF; END $$;`,
		typeName,
		typeName,
		strings.Join(values, "', '"),
	)

	if err := db.Exec(createTypeSQL).Error; err != nil {
		return fmt.Errorf("error creating type: %w", err)
	}

	return nil
}

func dropUnusedColumns(db *gorm.DB, dsts ...interface{}) {
	for _, dst := range dsts {
		stmt := &gorm.Statement{DB: db}
		stmt.Parse(dst)
		fields := stmt.Schema.Fields
		columns, _ := db.Debug().Migrator().ColumnTypes(dst)

		for i := range columns {
			found := false

			for j := range fields {
				if columns[i].Name() == fields[j].DBName {
					found = true
					break
				}
			}

			if !found {
				db.Migrator().DropColumn(dst, columns[i].Name())
			}
		}
	}
}

func Load() {
  log.Print("WAIT: Loading database")

	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  config.Postgres.ConnectionString,
		PreferSimpleProtocol: true,
	}), &gorm.Config{
		PrepareStmt: false,
	})
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	// NOTE: Currently, we are handling every error case when doing DB operations. There is no need for an extra logger
	db.Logger = logger.Discard

	Postgres = db

  log.Print("DONE: Loaded database")

  log.Print("WAIT: Loading migrations")

	if err := typeCheckAndCreate(db, "enum_language", language.ArrayString); err != nil {
		log.Fatal(err)
	}

	models := []any{
		&model.APIKey{},
		&model.User{},
		&model.Follower{},
		&model.Permission{},
		&model.UserPermission{},
	}

	if err := db.AutoMigrate(models...); err != nil {
		log.Fatal(err)
	}

	dropUnusedColumns(db, models...)

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

  log.Print("DONE: Loaded migrations")
}
