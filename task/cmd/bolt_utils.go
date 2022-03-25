package cmd

import (
	"github.com/boltdb/bolt"
)

func Put(path string, bucketName, key, value []byte) error {
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
