package main

import (
	"context"
	"fmt"
	"os"

	"github.com/biensupernice/krane/api"
	"github.com/biensupernice/krane/auth"
	"github.com/biensupernice/krane/db"
	"github.com/biensupernice/krane/docker"

	"github.com/biensupernice/krane/logger"
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

	l.Debugf("Starting krane server in [%s] mode", LogLevel)

	// Verify private key is provided
	if KranePrivateKey == "" {
		privK := "KbVHZLjpM3IUprwTSRvteRx+d8kmVecnEKvwAuJIaaw="
		l.Debugf("[KRANE_PRIVATE_KEY] not set, defaulting to %s", privK)
		KranePrivateKey = privK
	}

	// Set default port to `8080`
	if RestPort == "" {
		RestPort = "8080"
		l.Debugf("[KRANE_REST_PORT] not set, defaulting to %s", RestPort)
		os.Setenv("KRANE_REST_PORT", RestPort)
	}

	// Set default loglevel to `debug`
	if LogLevel == "" {
		LogLevel = "release"
		l.Debugf("[LogLevel] not set, defaulting to %s", LogLevel)
		os.Setenv("KRANE_LOG_LEVEL", LogLevel)
	}

	// Set default krane dir to `~/.krane`
	if KranePath == "" {
		homeDir := auth.GetHomeDir()
		KranePath = fmt.Sprintf("%s/.krane", homeDir)
		l.Debugf("[KRANE_PATH] not set, defaulting to %s", KranePath)
		os.Setenv("KRANE_PATH", KranePath)
	}

	// Set server configuration
	config = &api.Config{
		RestPort: RestPort,
		LogLevel: LogLevel,
	}

	l.Infof("🏗 krane path: %s\n", KranePath)
	l.Infof("🏗 krane log level: %s\n", LogLevel)
	l.Infof("🏗 krane port: %s\n", RestPort)

	// Create db
	err := db.New("krane.db")
	if err != nil {
		l.Fatalf("Unable to start db - %s", err.Error())
	}

	// Setup db - creates required buckets
	err = db.Setup()
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
	logger.Debugf("Using docker network - %s", netRes.ID)

	ctx.Done()
}

func main() {
	defer db.GetDB().Close()
	api.Start(*config)
}
