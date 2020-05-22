package handler

import (
	"github.com/biensupernice/krane/internal/api/http"
	"github.com/gin-gonic/gin"
)

// HealthHandler : returns 200
func HealthHandler(c *gin.Context) { http.Ok(c, map[string]string{"status": "healthy"}) }
