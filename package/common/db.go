package common

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	model "graphql_chat/package/model"
	"log"
	"os"
	"time"
)

func InitPostgres() (*gorm.DB, error) {

	password := os.Getenv("PASSWORD")

	dsn := fmt.Sprintf("host=localhost user=postgres password=%s dbname=postgres port=5432 sslmode=disable", password)

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second,   // Slow SQL threshold
			LogLevel:                  logger.Silent, // Log level
			IgnoreRecordNotFoundError: false,         // Ignore ErrRecordNotFound error for logger
			ParameterizedQueries:      false,         // Don't include params in the SQL log
			Colorful:                  true,          // Disable color
		},
	)

	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN: dsn,
	}), &gorm.Config{
		TranslateError: true,
		Logger:         newLogger})

	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(model.UserDB{},
		model.ChatDB{})

	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(model.MessageDB{})

	if err != nil {
		return nil, err
	}

	return db, err
}
