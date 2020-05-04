package main

import (
	"fmt"
	"log"
	"os/user"

	"github.com/biensupernice/krane/server"
	"github.com/dgraph-io/badger"
)

func main() {
	db := getDB()
	server.Run(db)
}

func getDB() *badger.DB {

	usr, err := user.Current()
	if err != nil {
		log.Fatalf("Unable to make ~/.krane dir %s", err.Error())
	}

	dbDir := fmt.Sprintf("%s/%s", usr.HomeDir, ".krane/db")
	db, err := badger.Open(badger.DefaultOptions(dbDir))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	return db
}
