package server

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/biensupernice/krane/auth"
	"github.com/biensupernice/krane/ds"
	"github.com/biensupernice/krane/http"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

var PrivateKey = []byte(os.Getenv("KRANE_PRIVATE_KEY"))

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
		return
	}

	// Login token should be encrypted with privKey
	ok, err := auth.ValidateTokenWithPubKey(req.Token)
	if !ok {
		http.BadRequest(c, err.Error())
		return
	}

	// Create new token with the valid req as the payload
	payload, _ := json.Marshal(req)
	reqID, err := auth.CreateToken(PrivateKey, payload)
	if err != nil {
		errMsg := fmt.Sprintf("Invalid request - %s", err.Error())
		http.BadRequest(c, errMsg)
		return
	}

	http.Ok(c, reqID)
}
