package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"os"
	"path/filepath"
	"time"

	badger "github.com/dgraph-io/badger/v4"
)

var BadgerDB *badger.DB

var ErrorExpiredApiKey = errors.New("apikey expired")

type ClientInfo struct {
	Client          string `json:"client"`
	PostUrl         string `json:"postUrl"`
	RequestInterval int    `json:"requestInterval"`
	ApiKey          string `json:"apiKey"`
	Expire          string `json:"expire"`
}

func init() {
	cwd, _ := os.Getwd()
	badgerPath := filepath.Join(cwd, ".badger")

	var err error
	BadgerDB, err = badger.Open(badger.DefaultOptions(badgerPath))
	if err != nil {
		log.Fatal(err)
	}
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

func LoadServerClientsInDb() {
	cwd, _ := os.Getwd()
	clientsFilePath := filepath.Join(cwd, "server-users")

	clients := getServerClients(clientsFilePath)

	currentDate := time.Now()

	for _, client := range clients {

		givenDate, err := time.Parse("2006-01-02", client.Expire)
		if err != nil {
			panic("Error parsing date")
		}

		if givenDate.Before(currentDate) {
			panic("Client apiKey expired")
		}

		errdb := BadgerDB.Update(func(txn *badger.Txn) error {
			return txn.Set([]byte(client.ApiKey), []byte(client.Expire))
		})

		if errdb != nil {
			panic("Error saving client ApiKey")
		}
	}
}

func validApiKey(key string) error {

	errdb := BadgerDB.View(func(txn *badger.Txn) error {

		item, err := txn.Get([]byte(key))
		if err != nil {
			return err
		}

		expire, err := item.ValueCopy(nil)
		if err != nil {
			return err
		}

		currentDate := time.Now()

		givenDate, err := time.Parse("2006-01-02", (string(expire)))
		if err != nil {
			return err
		}

		if givenDate.Before(currentDate) {
			return ErrorExpiredApiKey
		}

		return nil
	})

	if errdb != nil {
		return errdb
	}

	return nil
}
