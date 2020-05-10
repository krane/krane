package server

import (
	"github.com/biensupernice/krane/auth"
	"github.com/biensupernice/krane/ds"
	"github.com/biensupernice/krane/http"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type LoginRequest struct {
	Token string `json:"token" binding:"required"`
}

// PreLogin : id to authenticate login attempt
func PreLogin(c *gin.Context) {
	reqID := uuid.New().String()

	// Store `reqID` id in authentication bucket
	err := ds.Put(auth.Bucket, reqID, []byte(reqID))
	if err != nil {
		http.BadRequest(c, err.Error())
		return
	}

	http.Ok(c, map[string]string{"request_id": string(reqID)})
	return
}

// Login : handle login attempt
func Login(c *gin.Context) {
	var req LoginRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		http.BadRequest(c, err.Error())
		return
	}

	http.Ok(c, req)
}
