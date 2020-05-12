package server

import (
	"github.com/biensupernice/krane/server/handler"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// Config : server config
type Config struct {
	Port     string
	LogLevel string // release | debug
}

// Run : start server
func Run(cnf Config) {
	gin.SetMode(cnf.LogLevel)

	client := gin.New()

	// Middleware
	client.Use(gin.Recovery())
	client.Use(gin.Logger())
	client.Use(cors.Default())

	// Routes
	client.POST("/auth", handler.Auth)
	client.GET("/login", handler.Login)
	client.POST("/deploy", handler.DeployApp)

	client.PUT("/container/:containerID/stop", handler.StopContainer)
	client.PUT("/container/:containerID/start", handler.StartContainer)

	client.Run(cnf.Port)
}
