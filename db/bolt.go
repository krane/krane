package db

/**
Persistent key/value datastore using bolt
Operations:
- Get : get value by key
- GetAll : get all values in a bucket
- Put: store key-value pait
- Newdb: new instance of boltdb
- CreateBucket: new bucket that collects relevant data
- Setupdb : setup intial db buckets
- StartdbMetrics:  start a gp routine to fetch bolt metrics

**/

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/user"
	"time"

	"github.com/biensupernice/krane/logger"
	bolt "go.etcd.io/bbolt"
)

const (
	// AuthBucket : bucket used for storing auth related key-value data
	AuthBucket = "Auth"

	// SessionsBucket : bucket used for storing session related key-value data
	SessionsBucket = "Sessions"

	// DeploymentsBucket : bucket used for storing deployment related key-value data
	DeploymentsBucket = "Deployments"

	// SpecsBucket : bucket used for storing deployment spec
	SpecsBucket = "Specs"
)

var db *bolt.DB

// GetDB : return instance of boltdb
func GetDB() *bolt.DB { return db }

// Setup : initial db buckets
func Setup() (err error) {
	if db == nil {
		return fmt.Errorf("db not initiated")
	}

	// Buckets to create
	bkts := []string{
		AuthBucket,
		SessionsBucket,
		DeploymentsBucket,
		SpecsBucket,
	}

	// Iterate and create buckets
	for _, bkt := range bkts {
		err = CreateBucket(bkt)
		if err != nil {
			return
		}

		msg := fmt.Sprintf("Created %s Bucket", bkt)
		logger.Debug(msg)
	}

	return
}

// New : instance of bolt
func New(dbName string) (err error) {
	if db != nil {
		return
	}

	// Get base krane directory
	kPath := os.Getenv("KRANE_PATH")

	dbDir := fmt.Sprintf("%s/db", kPath)
	if kPath == "" {
		dbDir = "./db"
	}

	// Make db directory
	os.Mkdir(dbDir, 0777)
	logger.Debug("Created db")

	// Open the `dbName` data file in your current directory.
	// It will be created if it doesn't exist.
	options := &bolt.Options{Timeout: 1 * time.Second}

	dbPath := fmt.Sprintf("%s/%s", dbDir, dbName)
	dbInstance, err := bolt.Open(dbPath, 0600, options)
	if err != nil {
		return
	}

	db = dbInstance

	return
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
func CreateBucket(bktName string) (err error) {
	if db == nil {
		return fmt.Errorf("db not initialized")
	}

	return db.Update(func(tx *bolt.Tx) (err error) {
		_, err = tx.CreateBucketIfNotExists([]byte(bktName))
		return
	})
}

// Put : a key/value pair
func Put(bktName string, k string, v []byte) error {
	if db == nil {
		return fmt.Errorf("db not initialized")
	}

	return db.Update(func(tx *bolt.Tx) error {
		bkt := tx.Bucket([]byte(bktName))
		return bkt.Put([]byte(k), v)
	})
}

// Get : a value by key
func Get(bktName string, key string) (val []byte) {
	if db == nil {
		return nil
	}

	err := db.View(func(tx *bolt.Tx) error {
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
	if db == nil {
		return nil
	}

	err := db.View(func(tx *bolt.Tx) error {
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
	if db == nil {
		return fmt.Errorf("db not initialized")
	}

	return db.Update(func(tx *bolt.Tx) error {
		bkt := tx.Bucket([]byte(bktName))
		return bkt.Delete([]byte(key))
	})
}

// StartDBMetrics : capture db metrics
func StartDBMetrics() {
	go func() {
		// Grab the initial stats.
		prev := db.Stats()

		for {
			// Wait for 10s.
			time.Sleep(10 * time.Second)

			// Grab the current stats and diff them.
			stats := db.Stats()
			diff := stats.Sub(&prev)

			// Encode stats to JSON and print to STDERR.
			json.NewEncoder(os.Stderr).Encode(diff)

			// Save stats for the next loop.
			prev = stats
		}
	}()
}
