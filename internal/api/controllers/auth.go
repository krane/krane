package controllers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/docker/distribution/uuid"

	"github.com/krane/krane/internal/api/response"
	"github.com/krane/krane/internal/auth"
	"github.com/krane/krane/internal/constants"
	"github.com/krane/krane/internal/logger"
	"github.com/krane/krane/internal/session"
	"github.com/krane/krane/internal/store"
	"github.com/krane/krane/internal/utils"
)

// AuthRequest : the payload expected when authenticating with Krane
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
		logger.Warn("Invalid request id")
		response.HTTPBad(w, errors.New("unable to authenticate"))
		return
	}

	authKeys := auth.GetAuthorizeKeys()
	if len(authKeys) == 0 || authKeys[0] == "" {
		logger.Warn("no authorized keys found on the server")
		response.HTTPBad(w, errors.New("unable to authenticate"))
		return
	}

	// If any public key can be used to parse the incoming jwt token
	// the decode it, and use the phrase in the token with the one on the server
	claims := auth.VerifyAuthTokenWithAuthorizedKeys(authKeys, body.Token)
	if claims == nil || strings.Compare(serverPhrase, claims.Phrase) != 0 {
		logger.Warn("no authorized keys found on the server")
		response.HTTPBad(w, errors.New("invalid token"))
		return
	}

	// Create a new token and assign it to a session
	// Remove auth data from auth bucket
	err = store.Client().Remove(constants.AuthenticationCollectionName, body.RequestID)
	if err != nil {
		logger.Errorf("unable to remove authentication request %v", err)
		response.HTTPBad(w, err)
		return
	}

	sessionTkn := session.Token{SessionID: uuid.Generate().String()}
	signedTkn, err := session.CreateSessionToken(auth.GetServerPrivateKey(), sessionTkn)
	if err != nil {
		logger.Errorf("unable to create session token %v", err)
		response.HTTPBad(w, err)
		return
	}

	newSession := session.Session{
		ID:        sessionTkn.SessionID,
		Token:     signedTkn,
		ExpiresAt: utils.UnixToDate(utils.OneYear),
		User:      "root", // TODO: handle unique users
	}

	if err := session.Save(newSession); err != nil {
		response.HTTPBad(w, err)
		return
	}

	response.HTTPOk(w, newSession)
	return
}
