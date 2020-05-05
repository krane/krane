package server

import (
	"github.com/dgraph-io/badger"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// Run blah
func Run(db *badger.DB) {
	gin.SetMode("release")

	client := gin.Default()

	// Recover from any panics, returns 500
	client.Use(gin.Recovery())
	client.Use(cors.Default())

	// Routes
	client.POST("/login", LoginHandler)
	client.POST("/deploy", StartDeployHandler)

	client.POST("/container/:containerID/stop", StopContainerHandler)
	client.POST("/container/:containerID/start", StartContainerHandler)

	client.Run(":8000")
}
