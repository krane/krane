package boltdb

import (
	"bytes"
	"os"
	"time"

	"github.com/boltdb/bolt"
	"github.com/sirupsen/logrus"

	"github.com/biensupernice/krane/internal/storage"
)

type BoltDB struct {
	*bolt.DB
}

type BoltConfig struct {
	Path    string
	Buckets []string
}

// 0600 - Sets permissions so that:
// (U)ser / owner can read, can write and can't execute.
// (G)roup can't read, can't write and can't execute.
// (O)thers can't read, can't write and can't execute.
var fileMode os.FileMode = 0600

func Init() {
	options := &bolt.Options{Timeout: 10 * time.Second}

	db, err := bolt.Open("/tmp/krane.db", fileMode, options)

	if err != nil {
		logrus.Fatal(err.Error())
	}

	// set the instance for the current storage provider to boltdb
	storage.SetInstance(&BoltDB{db})

	return
}

func (b *BoltDB) Close() { b.Close() }

func (b *BoltDB) Put(collection string, key string, value []byte) error {
	return b.Update(func(tx *bolt.Tx) error {
		bkt, err := tx.CreateBucketIfNotExists([]byte(collection))
		if err != nil {
			return err
		}

		return bkt.Put([]byte(key), value)
	})
}

func (b *BoltDB) Get(collection, key string) (data []byte, err error) {
	b.View(func(tx *bolt.Tx) (err error) {
		bkt := tx.Bucket([]byte(collection))
		if bkt == nil {
			return
		}

		data = bkt.Get([]byte(key))
		return
	})

	if err != nil {
		return
	}

	return
}

func (b *BoltDB) GetAll(collection string) (data [][]byte, err error) {
	err = b.View(func(tx *bolt.Tx) (err error) {
		bkt := tx.Bucket([]byte(collection))
		if bkt == nil {
			return
		}

		bkt.ForEach(func(k, v []byte) (err error) {
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

// Iterate over a time range.
// minDate: RFC3339 sortable time string ie. 1990-01-01T00:00:00Z
// maxDate example: RFC3339 sortable time string ie. 2000-01-01T00:00:00Z
func (b *BoltDB) GetInRange(collection, minDate, maxDate string) (data [][]byte, err error) {
	err = b.View(func(tx *bolt.Tx) (err error) {
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
	return b.Update(func(tx *bolt.Tx) error {
		bkt := tx.Bucket([]byte(collection))
		if bkt == nil {
			return nil
		}
		return bkt.Delete([]byte(key))
	})
}
