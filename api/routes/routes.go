package routes

import (
	"github.com/burakkarasel/URL-Shortener/handlers"
	"github.com/gofiber/fiber/v2"
)

// SetupRoutes sets up the routes for the app
func SetupRoutes(app *fiber.App) {
	app.Get("/:url", handlers.ResolveURL)
	app.Post("/api/v1", handlers.ShortenURL)
}
