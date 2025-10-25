package db

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"GuildVault/internal/models"
)

var DB *gorm.DB

func InitDB() {
	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	port := os.Getenv("DB_PORT")

	if host == "" {
		host = "localhost"
		user = "postgres"
		password = "postgres"
		dbname = "dao_db"
		port = "5432"
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		host, user, password, dbname, port)

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect to PostgreSQL:", err)
	}

	DB.Exec(`CREATE EXTENSION IF NOT EXISTS "pgcrypto";`)

	err = DB.AutoMigrate(&models.User{}, &models.Tournament{})
	if err != nil {
		log.Fatal("failed to migrate:", err)
	}

	log.Println("Database connected and migrated successfully ✅")
}
