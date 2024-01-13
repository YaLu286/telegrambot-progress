package models

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {

	dsn := "host=10.5.0.4 user=postgres password=123 dbname=postgres port=5432 sslmode=disable TimeZone=Europe/Moscow"

	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		panic("Failed to connect to database!")
	}

	err = database.AutoMigrate(&Beer{}, &Admin{}, &Location{}, &UserSession{})

	if err != nil {
		return
	}

	DB = database
}
