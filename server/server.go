package server

import (
	"github.com/biensupernice/krane/auth"
	"github.com/biensupernice/krane/ds"
	"github.com/biensupernice/krane/http"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Config struct {
	Port     string
	LogLevel string // release | debug
}

// Run blah
func Run(cnf Config) {
	gin.SetMode(cnf.LogLevel)

	client := gin.New()

	// Recover from any panics, returns 500
	client.Use(gin.Recovery())
	client.Use(cors.Default())

	// Routes
	client.GET("/login", func(c *gin.Context) {
		id := uuid.New().String()

		ds.Put(auth.Bucket, id, []byte(id))

		val, _ := ds.Get(auth.Bucket, id)

		http.Ok(c, map[string]string{"uid": string(val)})
	})
	client.POST("/login", Login)
	client.POST("/deploy", DeployApp)

	client.Run(cnf.Port)
}
