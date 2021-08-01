package store

import (
	"bytes"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/boltdb/bolt"

	"github.com/krane/krane/internal/logger"
)

type BoltDB struct {
	*bolt.DB
}

type BoltConfig struct {
	Path    string
	Buckets []string
}

var once sync.Once
var instance *BoltDB

var (
	defaultBoltPath = "/tmp/krane.db"

	// 0600 - Sets permissions so that:
	// (U)ser / owner can read, can write and can't execute.
	// (G)roup can't read, can't write and can't execute.
	// (O)thers can't read, can't write and can't execute.
	fileMode os.FileMode = 0600
)

// Client boltdb client instance
func Client() Store { return instance }

// Connect connect to boltdb
func Connect(path string) *BoltDB {
	if instance != nil {
		logger.Info("Bolt instance already exists")
		return instance
	}

	logger.Info("Opening boltdb")

	options := &bolt.Options{Timeout: 30 * time.Second}

	if path == "" {
		path = defaultBoltPath
	}

	if err := os.MkdirAll(path, fileMode); err != nil {
		logger.Fatalf("Failed to create store at %s: %s", path, err.Error())
	}

	db, err := bolt.Open(path, fileMode, options)
	if err != nil {
		logger.Fatalf("Failed to open store at %s: %s", path, err.Error())
		return nil
	}

	once.Do(func() { instance = &BoltDB{db} })
	return instance
}

// Disconnect close boltdb client
func (b *BoltDB) Disconnect() {
	logger.Debug("Closing boltdb")
	if err := b.Close(); err != nil {
		logger.Errorf("Error closing boltdb %v", err)
	}
}

// Put upsert a key/value pair
func (b *BoltDB) Put(collection string, key string, value []byte) error {
	return instance.Update(func(tx *bolt.Tx) error {
		bkt, err := tx.CreateBucketIfNotExists([]byte(collection))
		if err != nil {
			return fmt.Errorf("unable to create bucket for %s", collection)
		}

		return bkt.Put([]byte(key), value)
	})
}

// Get get a key/value pair from a bucket
func (b *BoltDB) Get(collection, key string) (data []byte, err error) {
	err = instance.View(func(tx *bolt.Tx) error {
		bkt := tx.Bucket([]byte(collection))
		if bkt == nil {
			return nil
		}

		data = bkt.Get([]byte(key))
		return nil
	})

	if err != nil {
		return
	}

	return
}

// GetAll get all key/value pairs in a collection
func (b *BoltDB) GetAll(collection string) (data [][]byte, err error) {
	err = instance.View(func(tx *bolt.Tx) (err error) {
		bkt := tx.Bucket([]byte(collection))
		if bkt == nil {
			return
		}

		_ = bkt.ForEach(func(k, v []byte) (err error) {
			data = append(data, v)
			return
		})

		return
	})

	if err != nil {
		return
	}

	return
}

// GetInRange get key/value pairs within a time range
// minDate: RFC3339 sortable time string ie. 1990-01-01T00:00:00Z
// maxDate example: RFC3339 sortable time string ie. 2000-01-01T00:00:00Z
func (b *BoltDB) GetInRange(collection, minDate, maxDate string) (data [][]byte, err error) {
	err = instance.View(func(tx *bolt.Tx) (err error) {
		bkt := tx.Bucket([]byte(collection))
		if bkt == nil {
			return nil
		}

		c := bkt.Cursor()

		for k, v := c.Seek([]byte(minDate)); k != nil && bytes.Compare(k, []byte(maxDate)) <= 0; k, v = c.Next() {
			data = append(data, v)
		}
		return
	})
	return
}

func (b *BoltDB) Remove(collection string, key string) error {
	return instance.Update(func(tx *bolt.Tx) error {
		bkt := tx.Bucket([]byte(collection))
		if bkt == nil {
			// dont return err if bkt does not exists
			return nil
		}
		return bkt.Delete([]byte(key))
	})
}

func (b *BoltDB) DeleteCollection(collection string) error {
	return instance.Update(func(tx *bolt.Tx) error {
		return tx.DeleteBucket([]byte(collection))
	})
}

func (b *BoltDB) CreateCollection(collection string) error {
	return instance.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(collection))
		return err
	})
}
