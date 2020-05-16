package main

import (
	"fmt"
	"log"
	"os"

	"github.com/biensupernice/krane/auth"
	"github.com/biensupernice/krane/ds"
	"github.com/biensupernice/krane/server"
)

// Env
var (
	RestPort        = os.Getenv("KRANE_REST_PORT")
	LogLevel        = os.Getenv("KRANE_LOG_LEVEL")
	KranePath       = os.Getenv("KRANE_PATH")
	KranePrivateKey = os.Getenv("KRANE_PRIVATE_KEY")

	config *server.Config
)

func init() {

	// Set default krane dir
	if KranePath == "" {
		dir := auth.GetHomeDir()
		KranePath = fmt.Sprintf("%s/%s", dir, "/.krane")
	}

	log.Printf("🏗 krane path: %s", KranePath)
	log.Printf("🏗 krane log level: %s", LogLevel)
	log.Printf("🏗 krane port: %s", RestPort)

	// Create db
	_, err := ds.New("krane.db")
	if err != nil {
		log.Panicf("Unable to start db - %s", err.Error())
	}

	// Setup db
	err = ds.SetupDB()
	if err != nil {
		log.Panicf("Unable to setup db - %s", err.Error())
	}

	if KranePrivateKey == "" {
		log.Panicf("Private key [KRANE_PRIVATE_KEY] not set")
	}

	// Set default port
	if RestPort == "" {
		RestPort = "8080"
	}

	// Set default loglevel
	if LogLevel == "" {
		LogLevel = "debug"
	}

	// Set server configuration
	config = &server.Config{
		Port:     ":" + RestPort,
		LogLevel: LogLevel,
	}
}

func main() {
	defer ds.DB.Close()
	server.Run(*config)
}
