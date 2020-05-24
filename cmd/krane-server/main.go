package main

import (
	"context"
	"fmt"
	"os"

	"github.com/biensupernice/krane/docker"
	"github.com/biensupernice/krane/internal/api"
	"github.com/biensupernice/krane/internal/auth"
	"github.com/biensupernice/krane/internal/store"

	"github.com/biensupernice/krane/internal/logger"
)

// Env
var (
	RestPort        = os.Getenv("KRANE_REST_PORT")   //  Defaults to 8080
	LogLevel        = os.Getenv("KRANE_LOG_LEVEL")   // Defaults to release
	KranePath       = os.Getenv("KRANE_PATH")        // Defaults to ~/.krane
	KranePrivateKey = os.Getenv("KRANE_PRIVATE_KEY") // Private key for signing server tokens

	config *api.Config
)

func init() {
	// Initialize logger
	l := logger.NewLogger()

	l.Debug("Starting krane server in debug mode")

	// Verify private key is provided
	if KranePrivateKey == "" {
		l.Fatal("Private key [KRANE_PRIVATE_KEY] not set")
	}

	// Set default port to `8080`
	if RestPort == "" {
		RestPort = "8080"
		l.Debugf("[KRANE_REST_PORT] not set using %s", RestPort)
		os.Setenv("KRANE_REST_PORT", RestPort)
	}

	// Set default loglevel to `debug`
	if LogLevel == "" {
		LogLevel = "release"
		l.Debugf("[LogLevel] not set using %s", LogLevel)
		os.Setenv("KRANE_LOG_LEVEL", LogLevel)
	}

	// Set default krane dir to `~/.krane`
	if KranePath == "" {
		homeDir := auth.GetHomeDir()
		KranePath = fmt.Sprintf("%s/.krane", homeDir)
		l.Debugf("[KRANE_PATH] not set using %s", KranePath)
		os.Setenv("KRANE_PATH", KranePath)
	}

	// Set server configuration
	config = &api.Config{
		RestPort: RestPort,
		LogLevel: LogLevel,
	}

	fmt.Printf("🏗 krane path: %s\n", KranePath)
	fmt.Printf("🏗 krane log level: %s\n", LogLevel)
	fmt.Printf("🏗 krane port: %s\n", RestPort)

	// Start db
	_, err := store.NewDB("krane.db")
	if err != nil {
		l.Fatalf("Unable to start db - %s", err.Error())
	}

	// Setup db
	err = store.SetupDB()
	if err != nil {
		l.Fatalf("Unable to setup db - %s", err.Error())
	}

	// Create docker client
	_, err = docker.New()
	if err != nil {
		l.Fatalf("Error with docker - %s", err.Error())
	}

	// Create docker network
	ctx := context.Background()
	netName := "krane"
	netRes, err := docker.CreateBridgeNetwork(&ctx, netName)
	if err != nil {
		l.Fatalf("Error with docker network- %s", err.Error())
	}

	os.Setenv("KRANE_NETWORK_ID", netRes.ID)
	logger.Debugf("Create docker network - %s", netRes.ID)

	ctx.Done()
}

func main() {
	defer store.DB.Close()
	api.Start(*config)
}
