package handler

import (
	"github.com/biensupernice/krane/api/response"
	"github.com/gin-gonic/gin"
)

// HealthHandler : returns 200, used for checking the servers health
func HealthHandler(c *gin.Context) { response.Ok(c, map[string]string{"status": "Up"}) }
