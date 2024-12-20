package db

import (
	"errors"
	"fmt"
	"log"
	"reflect"
	"strings"

	"github.com/google/uuid"
	"github.com/swibly/swibly-api/config"
	"github.com/swibly/swibly-api/internal/model"
	"github.com/swibly/swibly-api/pkg/language"
	"github.com/swibly/swibly-api/pkg/notification"
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

	if err := typeCheckAndCreate(db, "notification_type", notification.ArrayString); err != nil {
		log.Fatal(err)
	}

	models := []any{
		&model.APIKey{},
		&model.User{},
		&model.Follower{},
		&model.Permission{},
		&model.UserPermission{},
		&model.PasswordResetKey{},

		&model.Project{},
		&model.ProjectOwner{},
		&model.ProjectPublication{},
		&model.ProjectUserFavorite{},
		&model.ProjectUserPermission{},

		&model.Component{},
		&model.ComponentOwner{},
		&model.ComponentHolder{},
		&model.ComponentPublication{},

		&model.Notification{},
		&model.NotificationUser{},
		&model.NotificationUserRead{},
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

	log.Print("WAIT: Validating API keys")

	var apikey model.APIKey
	if err := db.First(&apikey).Error; err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		key := &model.APIKey{
			Key:                uuid.New().String(),
			EnabledKeyManage:   1,
			EnabledAuth:        1,
			EnabledSearch:      1,
			EnabledUserFetch:   1,
			EnabledUserActions: 1,
			EnabledProjects:    1,
		}

		log.Print("Created API key for the first time: ", key.Key)
		db.Create(&key)
	}

	log.Print("DONE: Validated API keys")
}
