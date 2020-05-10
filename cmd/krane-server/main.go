package main

import (
	"log"
	"os"

	"github.com/biensupernice/krane/auth"
	"github.com/biensupernice/krane/server"
	"github.com/biensupernice/krane/store"
)

// Env
var (
	Port     = os.Getenv("PORT")
	LogLevel = os.Getenv("LOG_LEVEL")

	db     *store.DB
	config *server.Config
)

func init() {
	// Setup db
	db, _ = store.New("krane.db")

	store.CreateBucket(db, auth.Bucket)

	// Example usage
	store.Put(db, auth.Bucket, "123", []byte("True"))
	val, _ := store.Get(db, auth.Bucket, "123")
	log.Printf("---->%s\n", string(val))

	// Set default port
	if Port == "" {
		Port = "8080"
	}

	// Set default loglevel
	if LogLevel == "" {
		LogLevel = "release"
	}

	// Set server configuration
	config = &server.Config{
		Port:     ":" + Port,
		LogLevel: LogLevel,
	}
}

func main() {
	defer db.Close()
	server.Run(*config, db)
}
