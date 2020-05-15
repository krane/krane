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

type AuthResponse struct {
	Session Session     `json:"session"`
	Error   interface{} `json:"error"`
}

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
	var KranePrivateKey = []byte(os.Getenv("KRANE_PRIVATE_KEY"))

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

	// Read public key
	dir := fmt.Sprintf("%s/.ssh/pinga.pub", auth.GetHomeDir())
	pubKey, err := auth.ReadPubKeyFile(dir)
	if err != nil {
		http.BadRequest(c, &AuthResponse{Error: err.Error()})
		return
	}

	// Get claims from token
	claims, err := auth.ParseToken(string(pubKey), req.Token)
	if err != nil {
		http.BadRequest(c, &AuthResponse{Error: err.Error()})
		return
	}

	// Read contents of token
	authClaims, ok := claims.(*auth.AuthClaims)
	if !ok {
		http.BadRequest(c, &AuthResponse{Error: "Invalid token"})
		return
	}

	// Compare phrase from server against pharse from auth claims
	if strings.Compare(authClaims.Phrase, serverPhrase) != 0 {
		http.BadRequest(c, &AuthResponse{Error: "Invalid token"})
		return
	}

	// If reached here, authentication was succesful
	// create a new token and assign it to a session
	// New token expiration date
	exp := UnixToDate(auth.OneYear)

	// Create new token including sessionID
	sessionID := uuid.New()
	serverTkn, err := auth.CreateToken(KranePrivateKey, sessionID)
	if err != nil {
		errMsg := fmt.Sprintf("Invalid request - %s", err.Error())
		http.BadRequest(c, &AuthResponse{Error: errMsg})
		return
	}

	// Create a session with relevant data
	session := &Session{ID: sessionID, Token: serverTkn, ExpiresAt: exp}

	// Store session info to sessions bucket
	sessionBytes, _ := json.Marshal(session)
	ds.Put(auth.SessionsBucket, sessionID.String(), sessionBytes)

	// Remove auth info from auth bucket
	err = ds.Remove(auth.AuthBucket, req.RequestID)
	if err != nil {
		errMsg := fmt.Sprintf("Something went wrong - %s", err.Error())
		http.BadRequest(c, &AuthResponse{Error: errMsg})
		return
	}

	http.Ok(c, &AuthResponse{Session: *session})
}

// UnixToDate : format MM/DD/YYYY
func UnixToDate(u int64) string {
	t := time.Unix(u, 0)
	year, month, day := t.Date()
	return fmt.Sprintf("%02d/%d/%d", int(month), day, year)
}
