package handler

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/biensupernice/krane/auth"
	"github.com/biensupernice/krane/ds"
	"github.com/biensupernice/krane/server/http"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// AuthRequest : to authenticate with krane-server
type AuthRequest struct {
	RequestID string `json:"request_id" binding:"required"`
	Token     string `json:"token" binding:"required"`
}

// AuthResponse : response from auth request
type AuthResponse struct {
	Session Session     `json:"session"`
	Error   interface{} `json:"error"`
}

// Session : relevant data for authenticated sessions
type Session struct {
	ID        uuid.UUID `json:"id"`
	Token     string    `json:"token"`
	ExpiresAt string    `json:"expires_at"`
}

// LoginResponse : client response when attempting to login
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
	var req AuthRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		http.BadRequest(c, &AuthResponse{Error: err.Error()})
		return
	}

	// Check if request id is valid, get phrase stored on the server
	serverPhrase := string(ds.Get(auth.AuthBucket, req.RequestID))
	if serverPhrase == "" {
		err := fmt.Errorf("Unable to authenticate")
		http.BadRequest(c, &AuthResponse{Error: err.Error()})
		return
	}

	// Get authorized_keys
	authorizedKeys, err := auth.GetAuthorizedKeys("")
	if err != nil {
		http.BadRequest(c, &AuthResponse{Error: err.Error()})
		return
	}

	// Check if any key in authorized_keys serves as a valid key to parse incoming jwt token
	// If an authorized_key can parse incoming token, authClaims will be returned containing the phrase from the token
	authClaims, err := auth.VerifyAuthTokenWithAuthorizedKeys(authorizedKeys, req.Token)
	if err != nil {
		msg := "Unable to verify with authorized_keys, your .pub key is part of the authorized_keys file on your server"
		http.BadRequest(c, &AuthResponse{Error: msg})
		return
	}

	// Compare phrase from server against pharse from auth claims
	if strings.Compare(serverPhrase, authClaims.Phrase) != 0 {
		http.BadRequest(c, &AuthResponse{Error: "Invalid token"})
		return
	}

	// If reached here, authentication was succesful
	// create a new token and assign it to a session
	// Server private key used to create session token
	var serverPrivKey = []byte(os.Getenv("KRANE_PRIVATE_KEY"))

	// Create new token including sessionID, sign token with server private key
	sessionID := uuid.New()
	serverSignedTkn, err := auth.CreateToken(serverPrivKey, sessionID)
	if err != nil {
		errMsg := fmt.Sprintf("Invalid request - %s", err.Error())
		http.BadRequest(c, &AuthResponse{Error: errMsg})
		return
	}

	// Remove auth data from auth bucket
	err = ds.Remove(auth.AuthBucket, req.RequestID)
	if err != nil {
		errMsg := fmt.Sprintf("Something went wrong - %s", err.Error())
		http.BadRequest(c, &AuthResponse{Error: errMsg})
		return
	}

	// Create a session with relevant data
	session := &Session{
		ID:        sessionID,
		Token:     serverSignedTkn, // Token used for authenticating subsequent requests
		ExpiresAt: UnixToDate(auth.OneYear),
	}

	// Store session into sessions bucket
	sessionBytes, _ := json.Marshal(session)
	ds.Put(auth.SessionsBucket, sessionID.String(), sessionBytes)

	http.Ok(c, &AuthResponse{Session: *session})
}

// UnixToDate : format MM/DD/YYYY
func UnixToDate(u int64) string {
	t := time.Unix(u, 0)
	year, month, day := t.Date()
	return fmt.Sprintf("%02d/%d/%d", int(month), day, year)
}
