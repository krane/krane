package main

import (
	"fmt"
	"log"
	"os"
	"os/user"

	"github.com/biensupernice/krane/server"
	"github.com/dgraph-io/badger"
)

// Env
var Port = os.Getenv("PORT")
var LogLevel = os.Getenv("LOG_LEVEL")

func main() {
	db := startDB()

	// Configure default env
	cnfEnv()

	cnf := &server.Config{Port: fmt.Sprintf(":%s", Port), LogLevel: LogLevel}
	server.Run(*cnf, db)
}

// Configure default env
func cnfEnv() {
	// Set default port
	if Port == "" {
		Port = "8080"
	}

	// Set default loglevel
	if LogLevel == "" {
		LogLevel = "debug"
	}
}

func startDB() *badger.DB {
	usr, err := user.Current()
	if err != nil {
		log.Fatalf("Unable to make ~/.krane dir %s\n", err.Error())
	}

	dbDir := fmt.Sprintf("%s/%s", usr.HomeDir, ".krane/db")
	db, err := badger.Open(badger.DefaultOptions(dbDir))
	if err != nil {
		log.Fatalf("Unable to open db %s\n", err.Error())
	}
	defer db.Close()

	return db
}
