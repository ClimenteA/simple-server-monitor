package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/keyauth"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/joho/godotenv"
)

func validateApiKey(c *fiber.Ctx, key string) (bool, error) {

	expire := os.Getenv(key)

	currentDate := time.Now()

	givenDate, err := time.Parse("2006-01-02", expire)
	if err != nil {
		fmt.Println("Error parsing date:", err)
		return false, keyauth.ErrMissingOrMalformedAPIKey
	}

	if givenDate.Before(currentDate) {
		fmt.Println("Client apiKey expired")
		return false, keyauth.ErrMissingOrMalformedAPIKey
	}

	return true, nil

}

type ClientInfo struct {
	Client          string `json:"client"`
	PostUrl         string `json:"postUrl"`
	RequestInterval string `json:"requestInterval"`
	ApiKey          string `json:"apiKey"`
	Expire          string `json:"expire"`
}

func getServerClients(fp string) []ClientInfo {

	entries, err := os.ReadDir(fp)
	if err != nil {
		log.Fatal(err)
	}

	var clients []ClientInfo

	for _, entry := range entries {
		filePath := fp + "/" + entry.Name()
		file, _ := os.Open(filePath)
		defer file.Close()

		var client ClientInfo
		decoder := json.NewDecoder(file)
		if err := decoder.Decode(&client); err != nil {
			log.Printf("Error decoding JSON from file %s: %v\n", filePath, err)
			continue
		}

		clients = append(clients, client)
	}

	return clients

}

func setClientsEnvs(clients []ClientInfo) {

	currentDate := time.Now()

	for _, client := range clients {

		givenDate, err := time.Parse("2006-01-02", client.Expire)
		if err != nil {
			fmt.Println("Error parsing date:", err)
			continue
		}

		if givenDate.Before(currentDate) {
			fmt.Println("Client apiKey expired")
			continue
		}

		os.Setenv(client.ApiKey, client.Expire)
	}

}

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	PORT := os.Getenv("PORT")
	SERVER_USERS_PATH := os.Getenv("SERVER_USERS_PATH")

	clients := getServerClients(SERVER_USERS_PATH)
	setClientsEnvs(clients)

	app := fiber.New()

	app.Use(keyauth.New(keyauth.Config{
		KeyLookup: "header:ApiKey",
		Validator: validateApiKey,
	}))

	app.Use(helmet.New())
	app.Use(cors.New())
	app.Use(logger.New())
	app.Use(recover.New())

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	app.Listen(":" + PORT)
}

// TODO - API with storage a fast KV store

// List of:
// {
// 	id: 'bd7acbea-c1b1-46c2-aed5-3ad53abb28ba',
// 	title: 'First Item',
// 	message: 'Long message here lorem The title and onPress handler are required. It is recommended to set accessibilityLabel to help make your app usable by everyone.',
// 	timestamp: '08:30/23-03-2024',
// }
