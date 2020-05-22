package api

import (
	"strings"

	"github.com/biensupernice/krane/internal/api/handler"
	"github.com/biensupernice/krane/internal/api/middleware"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// Config : server config
type Config struct {
	RestPort string // port to use for rest api
	LogLevel string // release | debug
}

// Start : api server
func Start(cnf Config) {
	gin.SetMode(strings.ToLower(cnf.LogLevel))

	client := gin.New()

	// Middleware
	client.Use(gin.Recovery())
	client.Use(gin.Logger())
	client.Use(cors.Default())

	// Routes
	client.POST("/health", handler.HealthHandler)
	client.POST("/auth", handler.Auth)
	client.GET("/login", handler.Login)

	client.GET("/sessions", middleware.AuthSessionMiddleware(), handler.GetSessions)

	client.GET("/deployments", middleware.AuthSessionMiddleware(), handler.GetDeployments)
	client.GET("/deployments/:name", middleware.AuthSessionMiddleware(), handler.GetDeployment)
	client.POST("/deployments", middleware.AuthSessionMiddleware(), handler.CreateDeployment)

	client.GET("/containers", middleware.AuthSessionMiddleware(), handler.ListContainers)
	client.PUT("/containers/:containerID/stop", middleware.AuthSessionMiddleware(), handler.StopContainer)
	client.PUT("/containers/:containerID/start", middleware.AuthSessionMiddleware(), handler.StartContainer)

	client.Run(":" + cnf.RestPort)
}
