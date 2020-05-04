package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Ok response with status code 200
func Ok(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"code":    http.StatusOK,
		"data":    data,
	})
	return
}

// Created response with status code 200
func Created(c *gin.Context, data interface{}) {
	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"code":    http.StatusCreated,
		"data":    data,
	})
	return
}

// BadRequest response with status code 400
func BadRequest(c *gin.Context, err interface{}) {
	c.JSON(http.StatusBadRequest, gin.H{
		"success": false,
		"code":    http.StatusBadRequest,
		"error":   err,
	})
	return
}

// Unauthorized response with status code 400
func Unauthorized(c *gin.Context) {
	c.JSON(http.StatusUnauthorized, gin.H{
		"success": false,
		"code":    http.StatusUnauthorized,
		"error":   "Unauthorized request",
	})
	return
}
