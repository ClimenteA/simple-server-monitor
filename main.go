package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/keyauth"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"

	"github.com/ClimenteA/simple-server-monitor/handlers"

	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	go handlers.MonitorServerUsage()

	defer handlers.BadgerDB.Close()

	app := fiber.New()

	app.Use(keyauth.New(keyauth.Config{
		KeyLookup: "header:ApiKey",
		Validator: handlers.ValidateApiKey,
	}))

	app.Use(helmet.New())
	app.Use(cors.New())
	app.Use(logger.New())
	app.Use(recover.New())

	app.Post("/clear-database", func(c *fiber.Ctx) error {
		err := handlers.Clear()
		if err != nil {
			log.Println(err)
			c.Status(500)
			return c.JSON(map[string]string{"message": "cannot clear database"})
		}
		return c.JSON(map[string]string{"message": "database cleared"})
	})

	app.Get("/", func(c *fiber.Ctx) error {
		values, err := handlers.GetAll()
		if err != nil {
			c.Status(500)
			return c.JSON(values)
		}
		return c.JSON(values)
	})

	app.Listen(":" + os.Getenv("PORT"))
}

// List of:
// {
// 	id: 'bd7acbea-c1b1-46c2-aed5-3ad53abb28ba',
// 	title: 'First Item',
// 	message: 'Long message here lorem The title and onPress handler are required. It is recommended to set accessibilityLabel to help make your app usable by everyone.',
// 	timestamp: '08:30/23-03-2024',
// }
