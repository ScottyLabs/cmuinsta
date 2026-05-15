package db

import (
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/scottylabs/cmuinsta/backend/models"
)

var DB *gorm.DB

func Init() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL is not set")
	}

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	DB.AutoMigrate(
		&models.User{},
		&models.AppState{},
	)

	DB.FirstOrCreate(
		&models.AppState{ID: 1},
		models.AppState{
			ID: 1, QueuePosition: 0,
			QueueSize: 0,
		})
}
