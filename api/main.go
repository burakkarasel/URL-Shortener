package main

import (
	"log"
	"os"

	"github.com/burakkarasel/URL-Shortener/routes"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("cannot load env variables:", err)
	}

	app := fiber.New()
	app.Use(logger.New())

	routes.SetupRoutes(app)

	log.Fatal(app.Listen(os.Getenv("APP_PORT")))
}
