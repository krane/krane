package server

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Config struct {
	Port     string
	LogLevel string // release | debug
}

// Run : start server
func Run(cnf Config) {
	gin.SetMode(cnf.LogLevel)

	client := gin.New()

	// Gin middleware
	client.Use(gin.Recovery())
	client.Use(gin.Logger())
	client.Use(cors.Default())

	// Routes
	client.GET("/login", Login)
	client.POST("/auth", Auth)
	client.POST("/deploy", DeployApp)

	client.PUT("/container/:containerID/stop", StopContainerHandler)
	client.PUT("/container/:containerID/start", StartContainerHandler)

	client.Run(cnf.Port)
}
