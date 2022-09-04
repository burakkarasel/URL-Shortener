package handlers

import (
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/burakkarasel/URL-Shortener/database"
	"github.com/burakkarasel/URL-Shortener/helpers"
	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// request holds the request data
type request struct {
	URL         string        `json:"url"`
	CustomShort string        `json:"short"`
	Expiry      time.Duration `json:"expiry"`
}

// response holds the response data
type response struct {
	URL             string        `json:"url"`
	CustomShort     string        `json:"short"`
	Expiry          time.Duration `json:"expiry"`
	XRateRemaining  int           `json:"rate_limit"`
	XRateLimitReset time.Duration `json:"rate_limit_reset"`
}

// ShortenURL shortens a given URL
func ShortenURL(c *fiber.Ctx) error {
	// first we parse the body
	body := new(request)

	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse JSON"})
	}

	//* rate limiting

	// then we create a client for checking quota
	r2 := database.CreateClient(1)
	defer r2.Close()

	// we get the result for request's IP Address
	val, err := r2.Get(database.Ctx, c.IP()).Result()

	// if address's value is nil we create a new record for 30 minutes
	if err == redis.Nil {
		_ = r2.Set(database.Ctx, c.IP(), os.Getenv("API_QUOTA"), time.Minute*30).Err()
	} else {
		// we check for remaning limit
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Cannot connect to the database"})
		}
		valInt, _ := strconv.Atoi(val)
		// if it's exceeded we return service unavailable
		if valInt <= 0 {
			limit, _ := r2.TTL(database.Ctx, c.IP()).Result()
			return c.Status(http.StatusServiceUnavailable).JSON(fiber.Map{"error": "Rate Limit exceeded", "rate_limit_rest": limit / time.Nanosecond / time.Minute})
		}
	}

	//* check if the input is an actual URL

	if !govalidator.IsURL(body.URL) {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid URL"})
	}

	//* check for domain error

	if !helpers.RemoveDomainError(body.URL) {
		return c.Status(http.StatusServiceUnavailable).JSON(fiber.Map{"error": "You cannot "})
	}

	//* enforce https, SSL

	body.URL = helpers.EnforceHTTP(body.URL)

	//! first we check if the user gave us an custom short or not
	var id string

	if body.CustomShort == "" {
		id = uuid.New().String()[:6]
	} else {
		id = body.CustomShort
	}

	//* we create a new client to check urls
	r := database.CreateClient(0)
	defer r.Close()

	//! we check if the custom short in use
	val, _ = r.Get(database.Ctx, id).Result()
	if val != "" {
		return c.Status(http.StatusForbidden).JSON(fiber.Map{"error": "URL customs short is already in use"})
	}

	//* if user didnt pass expiry we set it 24 hours
	if body.Expiry == 0 {
		body.Expiry = 24 * time.Hour
	}

	//! id is custom short URL for our url if user set one it will be it otherwise its gonna be 6 digit uuid
	err = r.Set(database.Ctx, id, body.URL, body.Expiry).Err()

	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "cannot connect to the database"})
	}

	// we prepare the response
	resp := response{
		URL:             body.URL,
		CustomShort:     "",
		Expiry:          body.Expiry,
		XRateRemaining:  10,
		XRateLimitReset: 30 * time.Minute,
	}

	// if no error occurs we decrement quota of the user by 1
	r2.Decr(database.Ctx, c.IP())

	//* here we set remaining rate by fetching it from the DB
	val, _ = r2.Get(database.Ctx, c.IP()).Result()
	resp.XRateRemaining, _ = strconv.Atoi(val)

	//* here we set remaining time reset by fetching it from the DB
	ttl, _ := r2.TTL(database.Ctx, c.IP()).Result()
	resp.XRateLimitReset = ttl / time.Nanosecond / time.Minute

	//* here we set the custom short with our domain and id
	resp.CustomShort = os.Getenv("DOMAIN") + "/" + id

	return c.Status(http.StatusOK).JSON(resp)
}
