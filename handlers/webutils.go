package handlers

import (
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/keyauth"
)

func ValidateApiKey(c *fiber.Ctx, key string) (bool, error) {

	if key == os.Getenv("SIMPLE_SERVER_MONITOR_APIKEY") {
		return true, nil
	}
	return false, keyauth.ErrMissingOrMalformedAPIKey
}
