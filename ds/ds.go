package ds

/**
Persistent key/value datastore using bolt
Operations:
- Get : get value by key
- GetAll : get all values in a bucket
- Put: store key-value pait
- New: new instance of boltdb
- CreateBucket: new bucket that collects relevant data
- SetupDB : setup intial db buckets
- StartDBMetrics:  start a gp routine to fetch bolt metrics
**/

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/user"
	"time"

	"github.com/biensupernice/krane/auth"
	bolt "go.etcd.io/bbolt"
)

var (
	// DB : refrence to an instance of bolt
	DB *bolt.DB
)

// SetupDB : create initial db buckets
func SetupDB() error {
	if DB == nil {
		return fmt.Errorf("Unable to setup db")
	}

	// Setup auth bucket
	err := CreateBucket(auth.AuthBucket)
	if err != nil {
		return err
	}

	// Setup sessions bucket
	err = CreateBucket(auth.SessionsBucket)
	if err != nil {
		return err
	}

	return nil
}

// New : instance of bolt
func New(dbName string) (*bolt.DB, error) {
	if DB != nil {
		return nil, nil
	}

	// Open the `dbName` data file in your current directory.
	// It will be created if it doesn't exist.
	options := &bolt.Options{Timeout: 1 * time.Second}

	dbPath := fmt.Sprintf("%s/%s", BoltPath(), dbName)
	db, err := bolt.Open(dbPath, 0600, options)
	if err != nil {
		return nil, err
	}

	DB = db

	return db, nil
}

// BoltPath : location of boltdb
func BoltPath() string {
	usr, err := user.Current()
	if err != nil {
		log.Fatalf("Unable to get current user - %s\n", err.Error())
	}

	path := fmt.Sprintf("%s/%s", usr.HomeDir, ".config/@krane/db")
	return path
}

// CreateBucket : new bucket
func CreateBucket(bktName string) error {
	if DB == nil {
		return fmt.Errorf("db not initialized")
	}

	return DB.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(bktName))
		if err != nil {
			return err
		}
		return nil
	})
}

// Put : a key/value pair
func Put(bktName string, k string, v []byte) error {
	if DB == nil {
		return fmt.Errorf("db not initialized")
	}

	return DB.Update(func(tx *bolt.Tx) error {
		bkt := tx.Bucket([]byte(bktName))
		return bkt.Put([]byte(k), v)
	})
}

// Get : a value by key
func Get(bktName string, key string) (val []byte) {
	if DB == nil {
		return nil
	}

	err := DB.View(func(tx *bolt.Tx) error {
		bkt := tx.Bucket([]byte(bktName))
		if bkt == nil {
			return fmt.Errorf("Bucket %s not found", bktName)
		}
		val = bkt.Get([]byte(key))
		return nil
	})

	if err != nil {
		log.Fatal(err)
		return nil
	}

	return val
}

// GetAll : the values in a bucket
func GetAll(bktName string) (data []*[]byte) {
	if DB == nil {
		return nil
	}

	err := DB.View(func(tx *bolt.Tx) error {
		bkt := tx.Bucket([]byte(bktName))

		// Iterate all the keys in the bucket
		bkt.ForEach(func(k, v []byte) error {
			data = append(data, &v)
			return nil
		})

		return nil
	})

	if err != nil {
		return nil
	}

	return
}

// Remove : item by key
func Remove(bktName string, key string) error {
	if DB == nil {
		return fmt.Errorf("db not initialized")
	}

	return DB.Update(func(tx *bolt.Tx) error {
		bkt := tx.Bucket([]byte(bktName))
		return bkt.Delete([]byte(key))
	})
}

// StartDBMetrics : start a go routine capturing db metrics
func StartDBMetrics() {
	go func() {
		// Grab the initial stats.
		prev := DB.Stats()

		for {
			// Wait for 10s.
			time.Sleep(10 * time.Second)

			// Grab the current stats and diff them.
			stats := DB.Stats()
			diff := stats.Sub(&prev)

			// Encode stats to JSON and print to STDERR.
			json.NewEncoder(os.Stderr).Encode(diff)

			// Save stats for the next loop.
			prev = stats
		}
	}()
}
