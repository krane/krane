package store

import (
	"fmt"
	"log"
	"os/user"
	"time"

	"github.com/boltdb/bolt"
)

type DB = bolt.DB

// New : instance of bolt
func New(dbName string) (*DB, error) {
	// Open the `dbName` data file in your current directory.
	// It will be created if it doesn't exist.
	options := &bolt.Options{Timeout: 1 * time.Second}
	return bolt.Open(dbPath(dbName), 0600, options)
}

func dbPath(dbName string) string {
	usr, err := user.Current()
	if err != nil {
		log.Fatalf("Unable to make ~/.krane dir %s\n", err.Error())
	}

	path := fmt.Sprintf("%s/%s/%s", usr.HomeDir, ".krane/db", dbName)
	return path
}

// CreateBucket : new bucket
func CreateBucket(db *bolt.DB, bktName string) error {
	return db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucket([]byte(bktName))
		if err != nil {
			return err
		}
		return nil
	})
}

// Put : store data
func Put(db *bolt.DB, bktName string, k string, v []byte) error {
	return db.Update(func(tx *bolt.Tx) error {
		bkt := tx.Bucket([]byte(bktName))
		err := bkt.Put([]byte(k), v)
		return err
	})
}

// Get : retrieve data
func Get(db *DB, bucketName string, key string) (val []byte, length int) {
	err := db.View(func(tx *bolt.Tx) error {
		bkt := tx.Bucket([]byte(bucketName))
		if bkt == nil {
			return fmt.Errorf("Bucket %q not found!", bucketName)
		}
		val = bkt.Get([]byte(key))
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
	return val, len(string(val))
}
