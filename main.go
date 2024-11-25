package main

import (
	"fmt"
	"log"
	"os"

	"github.com/chajaykrishna/go-ecommerce/database"
	"github.com/chajaykrishna/go-ecommerce/models"
	"github.com/chajaykrishna/go-ecommerce/routes"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {

	database.Connect()
	// auto migrate the user models at start.
	if err := database.DB.AutoMigrate(&models.User{}); err != nil {
		log.Fatal("Failed to migrate database", err)
	}
	log.Printf("DB Migration complete")
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "3001"
	}

	app := fiber.New()
	fmt.Printf("%T", app)

	routes.RegisterRoutes(app)
	if err := app.Listen(":" + port); err != nil {
		log.Fatal("Cannot run app, err:", err)
	}
}
