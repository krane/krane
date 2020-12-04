package store

type Store interface {
	Disconnect()
	Get(collection, key string) ([]byte, error)
	GetAll(collection string) ([][]byte, error)
	GetInRange(collection, minTime, maxTime string) ([][]byte, error)
	Put(collection string, key string, value []byte) error
	Remove(collection string, key string) error
	DeleteCollection(collection string) error
	CreateCollection(collection string) error
}
