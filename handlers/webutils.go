package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/keyauth"
)

func ValidateApiKey(c *fiber.Ctx, key string) (bool, error) {

	err := validApiKey(key)
	if err != nil {
		return false, keyauth.ErrMissingOrMalformedAPIKey
	}
	return true, nil
}
