package handlers

import (
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/keyauth"
)

var APIKEY = os.Getenv("APIKEY")

func ValidateApiKey(c *fiber.Ctx, key string) (bool, error) {
	if APIKEY == key {
		return true, nil
	}
	return false, keyauth.ErrMissingOrMalformedAPIKey
}
