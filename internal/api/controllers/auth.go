package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/docker/distribution/uuid"

	"github.com/biensupernice/krane/internal/api/response"
	"github.com/biensupernice/krane/internal/auth"
	"github.com/biensupernice/krane/internal/constants"
	"github.com/biensupernice/krane/internal/logger"
	"github.com/biensupernice/krane/internal/session"
	"github.com/biensupernice/krane/internal/store"
	"github.com/biensupernice/krane/internal/utils"
)

// AuthRequest : the payload expected when authenticating with krane
type AuthRequest struct {
	RequestID string `json:"request_id" binding:"required"`
	Token     string `json:"token" binding:"required"`
}

// AuthenticateClientJWT  : authenticate a client signed jwt token. This usually gets called after a call to /login.
// The login route returns a server phrase which is signed using the clients private auth. To authenticate the clients signed
// token this route is called and the token is validated against the clients public auth on the host machine.
func AuthenticateClientJWT(w http.ResponseWriter, r *http.Request) {
	var body AuthRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		response.HTTPBad(w, err)
		return
	}

	// Check if request id is valid, get phrase stored on the server
	serverPhraseBytes, err := store.Client().Get(constants.AuthenticationCollectionName, body.RequestID)
	if err != nil {
		response.HTTPBad(w, err)
		return
	}

	serverPhrase := string(serverPhraseBytes)
	if serverPhrase == "" {
		logger.Debug("Invalid request id")
		response.HTTPBad(w, errors.New("unable to authenticate"))
		return
	}

	authKeys := auth.GetAuthorizeKeys()
	if len(authKeys) == 0 || authKeys[0] == "" {
		logger.Info("no authorized keys found on the server")
		response.HTTPBad(w, errors.New("unable to authenticate"))
		return
	}

	// If any public key can be used to parse the incoming jwt token
	// the decode it, and use the phrase in the token with the one on the server
	claims := auth.VerifyAuthTokenWithAuthorizedKeys(authKeys, body.Token)
	if claims == nil || strings.Compare(serverPhrase, claims.Phrase) != 0 {
		response.HTTPBad(w, errors.New("invalid token"))
		return
	}

	// Create a new token and assign it to a session
	// Remove auth data from auth bucket
	err = store.Client().Remove(constants.AuthenticationCollectionName, body.RequestID)
	if err != nil {
		response.HTTPBad(w, err)
		return
	}

	sessionTkn := &session.Token{SessionID: uuid.Generate().String()}
	signedTkn, err := session.CreateSessionToken(auth.GetServerPrivateKey(), *sessionTkn)
	if err != nil {
		response.HTTPBad(w, err)
		return
	}

	newSession := session.Session{
		ID:        sessionTkn.SessionID,
		Token:     signedTkn,
		ExpiresAt: UnixToDate(utils.OneYear),
		User:      "root",
	}

	if err := session.Save(newSession); err != nil {
		response.HTTPBad(w, err)
		return
	}

	response.HTTPOk(w, newSession)
	return
}

// UnixToDate : format unix date into MM/DD/YYYY
func UnixToDate(u int64) string {
	t := time.Unix(u, 0)
	year, month, day := t.Date()
	return fmt.Sprintf("%02d/%d/%d", int(month), day, year)
}
