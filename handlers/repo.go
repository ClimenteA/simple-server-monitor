package handlers

import (
	"log"
	"os"
	"path/filepath"

	badger "github.com/dgraph-io/badger/v4"
)

var BadgerDB *badger.DB

func init() {
	cwd, _ := os.Getwd()
	badgerPath := filepath.Join(cwd, ".badger")

	var dberr error
	BadgerDB, dberr = badger.Open(badger.DefaultOptions(badgerPath))
	if dberr != nil {
		log.Fatal(dberr)
	}
}

func Set(key, value string) error {
	err := BadgerDB.Update(func(txn *badger.Txn) error {
		return txn.Set([]byte(key), []byte(value))
	})
	return err
}

func Get(key string) (string, error) {
	var value string
	err := BadgerDB.View(func(txn *badger.Txn) error {
		item, errGet := txn.Get([]byte(key))
		if errGet != nil {
			return errGet
		}
		valueCopy, errCopy := item.ValueCopy(nil)
		if errCopy != nil {
			return errCopy
		}
		value = string(valueCopy)
		return nil
	})
	return value, err
}

func Del(key string) error {
	err := BadgerDB.Update(func(txn *badger.Txn) error {
		return txn.Delete([]byte(key))
	})
	return err
}
