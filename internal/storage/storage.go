package storage

import (
	"sync"
)

type Storage interface {
	Get(collection, key string) ([]byte, error)
	GetAll(collection string) ([][]byte, error)
	GetInRange(collection, minTime, maxTime string) ([][]byte, error)
	Put(collection string, key string, value []byte) error
	Remove(collection string, key string) error
	Close()
}

var once sync.Once

var instance Storage

// GetInstance of the current storage provider
func GetInstance() Storage { return instance }

// SetInstance creates a singleton for the provided storage provider
func SetInstance(s Storage) { once.Do(func() { instance = s }) }

// Get a record by key from the collection
func Get(collection, key string) ([]byte, error) { return GetInstance().Get(collection, key) }

// GetAll records from the collection
func GetAll(collection string) ([][]byte, error) { return GetInstance().GetAll(collection) }

// GetInRange get the records part of the time range.
// Note: this should only be used on collections where the record's key is an RFC3339 encoded time value.
func GetInRange(collection, minTime, maxTime string) ([][]byte, error) {
	return GetInstance().GetInRange(collection, minTime, maxTime)
}

// Put a key/value pair into the collection
func Put(collection string, key string, value []byte) error {
	return GetInstance().Put(collection, key, value)
}

// Remove a value by key from the collection
func Remove(collection string, key string) error { return GetInstance().Remove(collection, key) }

// Close the connection to the storage instance
func Close() { GetInstance().Close() }
