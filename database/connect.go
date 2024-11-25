package database

import (
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect(postgres_url string) error {
	dsn := postgres_url
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to Connect to database: ", err)
	}
	DB = db
	log.Println("Database connection successful!")
	return nil
}
