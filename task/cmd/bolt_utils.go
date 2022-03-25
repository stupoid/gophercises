package cmd

import (
	"log"
	"os"
	"path/filepath"

	"github.com/boltdb/bolt"
	"github.com/mitchellh/go-homedir"
)

// ResetDB removes .db file relative to homedir
func ResetDB(path string) {
	home, err := homedir.Dir()
	if err != nil {
		log.Fatal(err)
	}

	path = home + path
	err = os.Remove(path)
	if err != nil {
		log.Fatal(err)
	}
}

func Put(path string, bucketName, key, value []byte) error {
	home, err := homedir.Dir()
	if err != nil {
		log.Fatal(err)
	}

	path = home + path
	err = os.MkdirAll(filepath.Dir(path), os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}
	db, err := bolt.Open(path, 0600, nil)
	if err != nil {
		return err
	}
	defer db.Close()

	return db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(bucketName))
		if err != nil {
			return err
		}
		return b.Put(key, value)
	})
}

func Get(path string, bucketName, key []byte) ([]byte, error) {
	var value []byte

	home, err := homedir.Dir()
	if err != nil {
		log.Fatal(err)
	}

	path = home + path
	err = os.MkdirAll(filepath.Dir(path), os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}

	db, err := bolt.Open(path, 0600, nil)
	if err != nil {
		return value, err
	}
	defer db.Close()

	err = db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(bucketName))
		if err != nil {
			return err
		}
		v := b.Get(key)
		if v != nil {
			value = append([]byte{}, v...)
		}
		return nil
	})

	if err != nil {
		return value, err
	}

	return value, err
}
