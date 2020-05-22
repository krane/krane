package handler

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/biensupernice/krane/auth"
	"github.com/biensupernice/krane/internal/api/http"
	"github.com/biensupernice/krane/internal/store"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

var (
	// Server private key
	serverPrivKey = []byte(os.Getenv("KRANE_PRIVATE_KEY"))
)

// AuthRequest : to authenticate with krane-server
type AuthRequest struct {
	RequestID string `json:"request_id" binding:"required"`
	Token     string `json:"token" binding:"required"`
}

// AuthResponse : response from auth request
type AuthResponse struct {
	Session Session `json:"session"`
}

// Session : relevant data for authenticated sessions
type Session struct {
	ID        string `json:"id"`
	Token     string `json:"token"`
	ExpiresAt string `json:"expires_at"`
}

// SessionToken :  token for the authenticated session, sign token with server private key
type SessionToken struct {
	SessionID string `json:"session_id"`
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

	err := store.Put(store.AuthBucket, key, val)
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
		http.BadRequest(c, err.Error())
		return
	}

	// Check if request id is valid, get phrase stored on the server
	serverPhrase := string(store.Get(store.AuthBucket, req.RequestID))
	if serverPhrase == "" {
		errMsg := "Unable to authenticate"
		http.BadRequest(c, errMsg)
		return
	}

	// Get authorized_keys on the server
	authorizedKeys, err := auth.GetAuthorizedKeys("")
	if err != nil {
		http.BadRequest(c, err.Error())
		return
	}

	// Check if any key in authorized_keys serves as a valid key to parse incoming jwt token
	// If an authorized_key can parse incoming token, authClaims will be returned containing the phrase from the token
	authClaims, err := auth.VerifyAuthTokenWithAuthorizedKeys(authorizedKeys, req.Token)
	if err != nil {
		errMsg := "Unable to verify"
		http.BadRequest(c, errMsg)
		return
	}

	// Compare phrase from server against phrase from auth claims
	if strings.Compare(serverPhrase, authClaims.Phrase) != 0 {
		errMsg := "Invalid token"
		http.BadRequest(c, errMsg)
		return
	}

	// If reached here, authentication was succesful
	// create a new token and assign it to a session
	// Remove auth data from auth bucket
	err = store.Remove(store.AuthBucket, req.RequestID)
	if err != nil {
		errMsg := errors.Errorf("Something went wrong - %s", err.Error())
		http.BadRequest(c, errMsg)
		return
	}

	sessionID := uuid.New().String()
	sessionTkn := &SessionToken{SessionID: sessionID}
	signedataessionTkn, err := auth.CreateToken(serverPrivKey, sessionTkn)
	if err != nil {
		errMsg := errors.Errorf("Invalid request - %s", err.Error())
		http.BadRequest(c, errMsg)
		return
	}

	// Create a session with relevant data
	session := &Session{
		ID:        sessionID,
		Token:     signedataessionTkn, // Token used for authenticating subsequent requests for a session
		ExpiresAt: UnixToDate(auth.OneYear),
	}

	// Store session into sessions bucket
	sessionBytes, _ := json.Marshal(session)
	store.Put(store.SessionsBucket, sessionID, sessionBytes)

	http.Ok(c, &AuthResponse{Session: *session})
}

// UnixToDate : format MM/DD/YYYY
func UnixToDate(u int64) string {
	t := time.Unix(u, 0)
	year, month, day := t.Date()
	return fmt.Sprintf("%02d/%d/%d", int(month), day, year)
}
