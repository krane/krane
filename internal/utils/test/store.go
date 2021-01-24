package test

import (
	"os"

	"github.com/krane/krane/internal/store"
)

const DbPath = "./krane_test.db"

func SetupDb() {
	store.Connect(DbPath)
}

func TeardownDb() {
	os.Remove(DbPath)
	store.Client().Disconnect()
}
