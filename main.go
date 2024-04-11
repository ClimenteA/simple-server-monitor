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
	app.Use(cors.New(cors.Config{
		AllowOriginsFunc: func(origin string) bool { return true },
		AllowHeaders:     "*",
	}))
	app.Use(logger.New())
	app.Use(recover.New())

	app.Get("/simple-server-monitor/notifications", func(c *fiber.Ctx) error {
		results, err := handlers.GetAll()
		if err != nil {
			c.Status(500)
			return c.JSON(results)
		}
		parsedResults := handlers.ParseResults(results)
		handlers.Clear()
		return c.JSON(map[string][]handlers.ServerEvent{"data": parsedResults})
	})

	app.Post("/simple-server-monitor/save", func(c *fiber.Ctx) error {
		var err error

		event := new(handlers.ServerEvent)
		if err = c.BodyParser(event); err != nil {
			return err
		}

		err = handlers.SaveEvent(*event)
		if err != nil {
			c.Status(500)
			return c.JSON(map[string]string{"message": "failed to save"})
		}

		c.Status(201)
		return c.JSON(map[string]string{"message": "saved"})

	})

	app.Delete("/simple-server-monitor/delete/:eventId", func(c *fiber.Ctx) error {
		var err error

		eventId := c.Params("eventId")

		err = handlers.DeleteEvent(eventId)
		if err != nil {
			c.Status(500)
			return c.JSON(map[string]string{"message": "failed to delete"})
		}

		c.Status(200)
		return c.JSON(map[string]string{"message": "deleted"})

	})

	app.Delete("/simple-server-monitor/clear-database", func(c *fiber.Ctx) error {
		err := handlers.Clear()
		if err != nil {
			log.Println(err)
			c.Status(500)
			return c.JSON(map[string]string{"message": "cannot clear database"})
		}
		return c.JSON(map[string]string{"message": "database cleared"})
	})

	app.Listen(":" + os.Getenv("SIMPLE_SERVER_MONITOR_PORT"))
}
