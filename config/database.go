package config

import (
	"fmt"
	"log"
	"os"

	"roottrack-backend/models"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
	// Load .env only for local development
	_ = godotenv.Load()

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL is not set")
	}

	database, err := gorm.Open(postgres.Open(dbURL), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Auto Migration
	err = database.AutoMigrate(
		&models.User{},
		&models.Routine{},
		&models.SheddingLog{},
		&models.Product{},
		&models.ProgressPhoto{},
	)
	if err != nil {
		log.Fatal("Migration failed:", err)
	}

	DB = database
	fmt.Println("Database connected and migrated successfully.")
}
