package server

import (
	"github.com/dgraph-io/badger"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func Run(db *badger.DB) {
	gin.SetMode("release")

	client := gin.Default()

	// Recover from any panics, returns 500
	client.Use(gin.Recovery())
	client.Use(cors.Default())

	// Routes
	client.POST("/login", LoginHandler)
	client.POST("/deploy", DeployHandler)

	client.Run(":8000")
}
