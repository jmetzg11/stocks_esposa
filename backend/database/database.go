package database

import (
	"log"
	"stocks/backend/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() error {
	var err error
	DB, err = gorm.Open(sqlite.Open("data/stocks.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database", err)
	}

	DB.AutoMigrate(&models.Investments{}, &models.Transactions{})

	return nil
}
