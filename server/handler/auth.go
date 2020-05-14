package handler

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/biensupernice/krane/auth"
	"github.com/biensupernice/krane/ds"
	"github.com/biensupernice/krane/server/http"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// LoginRequest : to authenticate with krane-server
type LoginRequest struct {
	RequestID string `json:"request_id" binding:"required"`
	Token     string `json:"token" binding:"required"`
}

// LoginResponse : response sent to client pre-login
type LoginResponse struct {
	RequestID string `json:"request_id" binding:"required"`
	Phrase    string `json:"phrase" binding:"required"`
}

// Login : returns request_id to authenticate login attempt
func Login(c *gin.Context) {
	reqID := uuid.New()

	// Store `reqID` in auth bucket
	key := reqID.String()
	val := []byte(fmt.Sprintf("Hey krane, %s", key))

	err := ds.Put(auth.AuthBucket, key, val)
	if err != nil {
		http.BadRequest(c, err.Error())
		return
	}

	resp := &LoginResponse{RequestID: key, Phrase: string(val)}
	http.Ok(c, resp)
}

// Auth : handle login attempt
func Auth(c *gin.Context) {
	var KranePrivateKey = []byte(os.Getenv("KRANE_PRIVATE_KEY"))

	var req LoginRequest

	err := c.ShouldBindJSON(&req)
	if err != nil {
		http.BadRequest(c, err.Error())
		return
	}

	// Check if request id is valid
	phrase, len := ds.Get(auth.AuthBucket, req.RequestID)
	if phrase == nil || len == -1 {
		err := fmt.Errorf("Unable to authenticate, login request not found")
		http.BadRequest(c, err)
		return
	}

	// Validate token was encrypted with private key
	ok, err := auth.ValidateTokenWithPubKey(req.Token)
	if !ok {
		http.BadRequest(c, err.Error())
		return
	}

	// Read pub key
	pubKey, err := auth.ReadPubKeyFile("")
	if err != nil {
		http.BadRequest(c, err.Error())
		return
	}

	// Parse token & get phrase
	p, err := auth.ParseAuthToken(pubKey, req.Token)
	if err != nil {
		http.BadRequest(c, err.Error())
		return
	}

	// Compare phrase
	if string(p) != string(phrase) {
		http.BadRequest(c, fmt.Errorf("Invalid token"))
		return
	}

	// Create new token with req id in payload
	data := map[string]string{"request_id": req.RequestID}
	tkn, err := auth.CreateToken(KranePrivateKey, data)
	if err != nil {
		errMsg := fmt.Sprintf("Invalid request - %s", err.Error())
		http.BadRequest(c, errMsg)
		return
	}

	// Remove request_id from auth bucket
	err = ds.Remove(auth.AuthBucket, req.RequestID)
	if err != nil {
		errMsg := fmt.Sprintf("Something went wrong - %s", err.Error())
		http.BadRequest(c, errMsg)
		return
	}

	// token expiration date
	exp := UnixToDate(auth.OneYear)

	type IdentityData struct {
		Token     string `json:"token"`
		ExpiresAt string `json:"expires_at"`
	}

	// Store identity info to identity bucket
	identityData := &IdentityData{Token: tkn, ExpiresAt: exp}
	b, _ := json.Marshal(identityData)
	ds.Put(auth.IdentityBucket, req.RequestID, b)

	http.Ok(c, identityData)
}

// UnixToDate : format MM/DD/YYYY
func UnixToDate(u int64) string {
	t := time.Unix(u, 0)
	year, month, day := t.Date()
	return fmt.Sprintf("%02d/%d/%d", int(month), day, year)
}
