package main

import (
	"log"
	"os"

	"github.com/biensupernice/krane/ds"
	"github.com/biensupernice/krane/server"
)

// Env
var (
	Port     = os.Getenv("PORT")
	LogLevel = os.Getenv("LOG_LEVEL")

	config *server.Config
)

func init() {
	err := os.Setenv("KRANE_PRIVATE_KEY", "biensupernice") // Change this :]
	if err != nil {
		log.Panicf("Unable to set KRANE_PRIVATE_KEY")
	}

	// Create db
	_, err = ds.New("krane.db")
	if err != nil {
		log.Panicf("Unable to start db - %s", err.Error())
	}

	// Setup db
	err = ds.SetupDB()
	if err != nil {
		log.Panicf("Unable to setup db - %s", err.Error())
	}

	// Set default port
	if Port == "" {
		Port = "8080"
	}

	// Set default loglevel
	if LogLevel == "" {
		LogLevel = "debug"
	}

	// Set server configuration
	config = &server.Config{
		Port:     ":" + Port,
		LogLevel: LogLevel,
	}
}

func main() {
	defer ds.DB.Close()
	server.Run(*config)
}
