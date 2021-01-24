package controllers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/docker/distribution/uuid"

	"github.com/krane/krane/internal/api/response"
	"github.com/krane/krane/internal/auth"
	"github.com/krane/krane/internal/logger"
	"github.com/krane/krane/internal/session"
	"github.com/krane/krane/internal/utils"
)

// AuthRequest represents the payload expected when authenticating with Krane
type AuthRequest struct {
	RequestID string `json:"request_id" binding:"required"`
	Token     string `json:"token" binding:"required"`
}

// LoginResponse is the response received when you initially want to authenticate.
// The request_id is a uuid stored for future validation and the phrase is a generated phrased
// containing that request_id meant to be signed by the clients private key to later be unsigned
// by the clients public key to establish an authenticated sessions
type LoginResponse struct {
	RequestID string `json:"request_id"`
	Phrase    string `json:"phrase"`
}

// RequestLoginPhrase request a preliminary login request for authentication with the krane server.
// This will return a request id and phrase. The phrase should be encrypted using the clients private key.
// This route does not return a token. You must use /auth and provide the signed phrase.
func RequestLoginPhrase(w http.ResponseWriter, _ *http.Request) {
	reqID, phrase, err := auth.CreateAuthenticationPhrase()
	if err != nil {
		response.HTTPBad(w, err)
	}

	response.HTTPOk(w, LoginResponse{
		RequestID: reqID,
		Phrase:    phrase,
	})
}

// AuthenticateClientJWT authenticates a client signed jwt token. This usually gets called after a call to /login.
// The login route returns a server phrase which is signed using the clients private auth. To authenticate the clients signed
// token this route is called and the token is validated against the clients public auth on the host machine.
func AuthenticateClientJWT(w http.ResponseWriter, r *http.Request) {
	var body AuthRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		response.HTTPBad(w, err)
		return
	}

	// We check if the request id is valid to ensure we've stored a generate server
	// phrase for that client id during the initial login request
	serverPhrase, err := auth.GetAuthenticationPhrase(body.RequestID)
	if err != nil {
		response.HTTPBad(w, err)
		return
	}

	// Grab all the authorized keys on the host machine
	// which will be used to decode the jwt token
	authKeys := auth.GetServerAuthorizeKeys()
	if len(authKeys) == 0 || authKeys[0] == "" {
		logger.Warn("no authorized keys found on the server")
		response.HTTPBad(w, errors.New("unable to authenticate"))
		return
	}

	// If any public key can be used to parse the incoming jwt token decode it,
	// and passes the phrase comparison between incoming and server phrase,
	// that token will be
	claims := session.VerifyAuthTokenWithAuthorizedKeys(authKeys, body.Token)
	if claims == nil || strings.Compare(serverPhrase, claims.Phrase) != 0 {
		logger.Warn("no authorized keys found on the server")
		response.HTTPBad(w, errors.New("invalid token"))
		return
	}
	// revoke the request id to ensure no one else can
	// use the same request id to create tokens
	if err := auth.RevokeAuthenticationRequest(body.RequestID); err != nil {
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
