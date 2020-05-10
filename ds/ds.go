package ds

/**
	ds : Datastore is a simple key value wrapper around bolt

	Operations:
		- Get : get value by key
		- Put: store key-value pait
		- New: new instance of boltdb
		- CreateBucket: new bucket that collects relevant data
**/

import (
	"fmt"
	"log"
	"os/user"
	"time"

	"github.com/boltdb/bolt"
)

var (
	DB *bolt.DB
)

// New : instance of bolt
func New(dbName string) error {

	// Open the `dbName` data file in your current directory.
	// It will be created if it doesn't exist.
	options := &bolt.Options{Timeout: 1 * time.Second}

	dbPath := fmt.Sprintf("%s/%s", BoltPath(), dbName)
	db, err := bolt.Open(dbPath, 0600, options)
	if err != nil {
		return err
	}

	DB = db

	return nil
}

// BoltPath : location of boltdb
func BoltPath() string {
	usr, err := user.Current()
	if err != nil {
		log.Fatalf("Unable to get current user - %s\n", err.Error())
	}

	path := fmt.Sprintf("%s/%s", usr.HomeDir, ".krane/db")
	return path
}

// CreateBucket : new bucket
func CreateBucket(bktName string) error {
	return DB.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucket([]byte(bktName))
		if err != nil {
			return err
		}
		return nil
	})
}

// Put : store data
func Put(bktName string, k string, v []byte) error {
	return DB.Update(func(tx *bolt.Tx) error {
		bkt := tx.Bucket([]byte(bktName))
		err := bkt.Put([]byte(k), v)
		return err
	})
}

// Get : retrieve data
func Get(bucketName string, key string) (val []byte, length int) {
	err := DB.View(func(tx *bolt.Tx) error {
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
