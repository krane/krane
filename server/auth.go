package server

import (
	"fmt"

	"github.com/biensupernice/krane/auth"
	"github.com/biensupernice/krane/ds"
	"github.com/biensupernice/krane/http"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// LoginRequest : to authenticate with krane-server
type LoginRequest struct {
	RequestID string `json:"request_id" binding:"required"`
	Token     string `json:"token" binding:"required"`
	User      string `json:"user" binding:"required"`
}

// PreLogin : id to authenticate login attempt
func PreLogin(c *gin.Context) {
	reqID := uuid.New()

	// Store `reqID` in auth bucket
	key := reqID.String()
	val := []byte(fmt.Sprintf("Hey krane, %s", key))

	err := ds.Put(auth.Bucket, key, val)
	if err != nil {
		http.BadRequest(c, err.Error())
		return
	}

	http.Ok(c, map[string]string{"request_id": reqID.String()})
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

	phrase, len := ds.Get(auth.Bucket, req.RequestID)
	if phrase == nil || len == -1 { // verify requestID is not nil
		err := fmt.Errorf("Unable to authenticate, login request not found")
		http.BadRequest(c, err)
	}

	auth.ValidateToken("", req.Token, string(phrase))

	http.Ok(c, req)
}
