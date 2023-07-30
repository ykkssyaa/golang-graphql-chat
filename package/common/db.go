package common

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	model "graphql_chat/package/model"
	"os"
)

func InitPostgres() (*gorm.DB, error) {

	password := os.Getenv("PASSWORD")

	dsn := fmt.Sprintf("host=localhost user=postgres password=%s dbname=postgres port=5432 sslmode=disable", password)

	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN: dsn,
	}), &gorm.Config{})

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
