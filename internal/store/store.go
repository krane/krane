package store

import "github.com/boltdb/bolt"

type Store interface {
	Get(collection, key string) ([]byte, error)
	GetAll(collection string) ([][]byte, error)
	GetInRange(collection, minTime, maxTime string) ([][]byte, error)
	Put(collection string, key string, value []byte) error
	Remove(collection string, key string) error
	GetBucket(collection string) *bolt.Bucket
	Shutdown()
}
