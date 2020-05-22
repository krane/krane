package store

/**
Persistent key/value datastore using bolt
Operations:
- Get : get value by key
- GetAll : get all values in a bucket
- Put: store key-value pait
- NewDB: new instance of boltdb
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

	"github.com/biensupernice/krane/internal/logger"
	bolt "go.etcd.io/bbolt"
)

const (
	// AuthBucket : bucket used for storing auth related key-value data
	AuthBucket = "Auth"

	// SessionsBucket : bucket used for storing session related key-value data
	SessionsBucket = "Sessions"

	// DeploymentsBucket : bucket used for storing deployment related key-value data
	DeploymentsBucket = "Deployments"

	// TemplatesBucket : bucket used for storing deployment templates
	TemplatesBucket = "Templates"
)

var (
	// DB : refrence to an instance of bolt
	DB *bolt.DB
)

// SetupDB : initial db buckets
func SetupDB() error {
	if DB == nil {
		return fmt.Errorf("Unable to setup db")
	}

	// Bucket to create
	bkts := []string{
		AuthBucket,
		SessionsBucket,
		DeploymentsBucket,
		TemplatesBucket,
	}

	// Iterate and create buckets
	for i := 0; i < len(bkts); i++ {
		err := CreateBucket(bkts[i])
		if err != nil {
			return err
		}

		msg := fmt.Sprintf("Created %s Bucket", bkts[i])
		logger.Debug(msg)
	}

	return nil
}

// NewDB : instance of bolt
func NewDB(dbName string) (*bolt.DB, error) {
	if DB != nil {
		return nil, nil
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
