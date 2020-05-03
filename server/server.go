package server

import (
	"github.com/dgraph-io/badger"
	"github.com/gin-gonic/gin"
)

func Run(db *badger.DB) {
	r := gin.Default()

	r.POST("/login", LoginHandler)
	r.POST("/deploy", DeployHandler)

	r.Run()
}
