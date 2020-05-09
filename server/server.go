package server

import (
	"github.com/dgraph-io/badger"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Config struct {
	Port     string
	LogLevel string // release | debug
}

// Run blah
func Run(cnf Config, db *badger.DB) {
	gin.SetMode(cnf.LogLevel)

	client := gin.New()

	// Recover from any panics, returns 500
	client.Use(gin.Recovery())
	client.Use(cors.Default())

	// Routes
	client.POST("/login", Login)
	client.POST("/deploy", DeployApp)

	client.Run(cnf.Port)
}
