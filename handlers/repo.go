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

func GetAll() ([]map[string]string, error) {

	values := []map[string]string{}

	err := BadgerDB.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchValues = false
		it := txn.NewIterator(opts)
		defer it.Close()
		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()
			key := item.Key()
			value, err := item.ValueCopy(nil)
			if err != nil {
				return err
			}
			values = append(values, map[string]string{string(key): string(value)})
		}
		return nil
	})

	return values, err

}

func Del(key string) error {
	err := BadgerDB.Update(func(txn *badger.Txn) error {
		return txn.Delete([]byte(key))
	})
	return err
}

func Clear() error {
	BadgerDB.RunValueLogGC(0.9)

	var keys [][]byte
	err := BadgerDB.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchValues = false
		it := txn.NewIterator(opts)
		defer it.Close()
		for it.Rewind(); it.Valid(); it.Next() {
			keys = append(keys, it.Item().KeyCopy(nil))
		}
		return nil
	})

	if err != nil {
		return err
	}

	for _, key := range keys {
		err := BadgerDB.Update(func(txn *badger.Txn) error {
			return txn.Delete(key)
		})
		if err != nil {
			return err
		}
	}

	return nil
}
