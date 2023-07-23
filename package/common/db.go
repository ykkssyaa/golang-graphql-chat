package common

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	model "graphql_chat/package/model"
)

func InitPostgres() (*gorm.DB, error) {

	dsn := "host=localhost user=postgres password=yksadm dbname=postgres port=5432 sslmode=disable"

	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN: dsn,
	}), &gorm.Config{})

	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(model.UserDB{},
		model.MessageDB{},
		model.ChatDB{})

	if err != nil {
		return nil, err
	}

	return db, err
}
