package database

import (
	"log"
	"os"
	"stocks/backend/models"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func Connect() error {
	var err error

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // IO writer
		logger.Config{
			SlowThreshold:             time.Second, // Slow SQL threshold
			LogLevel:                  logger.Warn, // Only log warnings and errors
			IgnoreRecordNotFoundError: true,        // Ignore "record not found" errors
			Colorful:                  true,        // Keep colors enabled
		},
	)

	DB, err = gorm.Open(sqlite.Open("data/stocks.db"), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		log.Fatal("Failed to connect to database", err)
	}

	DB.AutoMigrate(&models.Investment{}, &models.Transaction{})

	return nil
}
