package handler

import (
	"github.com/biensupernice/krane/internal/api/http"
	"github.com/gin-gonic/gin"
)

// HealthHandler : returns 200, used for checking the servers health
func HealthHandler(c *gin.Context) { http.Ok(c, map[string]string{"status": "Up"}) }
