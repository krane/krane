package main

import (
	"log"

	"github.com/biensupernice/krane/server"
	"github.com/dgraph-io/badger"
)

func main() {
	db := getDB()
	server.Run(db)
}

func getDB() *badger.DB {
	db, err := badger.Open(badger.DefaultOptions("~/.krane"))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	return db
}
