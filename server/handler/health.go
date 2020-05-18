package handler

import (
	"github.com/biensupernice/krane/server/http"
	"github.com/gin-gonic/gin"
)

// HealthHandler
func HealthHandler(c *gin.Context) {
	http.Ok(c, map[string]string{"message": "healthy"})
}
