package main

import (
	"log"
	"os"

	"github.com/sirupsen/logrus"

	"github.com/biensupernice/krane/api"
	"github.com/biensupernice/krane/internal/logging"
	"github.com/biensupernice/krane/internal/storage"
	"github.com/biensupernice/krane/internal/storage/boltdb"
	"github.com/biensupernice/krane/pkg/bbq"
	"github.com/biensupernice/krane/pkg/docker"
)

func init() {
	log.Println("Starting Krane...")

	logLvl := os.Getenv("LOG_LEVEL")
	if logLvl == "" {
		_ = os.Setenv("LOG_LEVEL", logging.INFO)
	}
	logging.ConfigureLogrus()

	listenAddr := os.Getenv("LISTEN_ADDRESS")
	if listenAddr == "" {
		logrus.Fatal("LISTEN_ADDRESS must be set and non-empty")
	}

	docker.Init()
	boltdb.Init()
	bbq.InitJobQueue()
}

func main() {
	api.Run()
	defer storage.Close()
}
