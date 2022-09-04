package handlers

import (
	"net/http"

	"github.com/burakkarasel/URL-Shortener/database"
	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
)

// ResolveURL redirects the user to real url
func ResolveURL(c *fiber.Ctx) error {
	// first we get the url from the request's URI
	url := c.Params("url")

	// we create a new client and we close it after func runs
	r := database.CreateClient(0)
	defer r.Close()

	// then we check in database for record
	val, err := r.Get(database.Ctx, url).Result()

	// then we check for error
	if err != nil {
		if err == redis.Nil {
			return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": "custom url is not found in the database"})
		}
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "cannot connect to the database"})
	}

	rInr := database.CreateClient(1)
	defer rInr.Close()

	_ = rInr.Incr(database.Ctx, "counter")

	// finally if no error occurs we redirect user to the specified URL
	return c.Redirect(val, http.StatusPermanentRedirect)
}
