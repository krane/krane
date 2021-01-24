package controllers

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/docker/distribution/uuid"
	"github.com/gorilla/mux"

	"github.com/krane/krane/internal/api/response"
	"github.com/krane/krane/internal/auth"
	"github.com/krane/krane/internal/logger"
	"github.com/krane/krane/internal/session"
	"github.com/krane/krane/internal/utils"
)

// GetSessions returns a list of user sessions. A session is an authenticated client with a valid access token
func GetSessions(w http.ResponseWriter, _ *http.Request) {
	sessions, err := session.GetAllSessions()
	if err != nil {
		response.HTTPBad(w, err)
		return
	}

	response.HTTPOk(w, sessions)
	return
}

// CreateSession returns a session with an active access token. Access tokens are useful for CI
func CreateSession(w http.ResponseWriter, r *http.Request) {
	user := utils.QueryParamOrDefault(r, "user", "")

	if user == "" {
		response.HTTPBad(w, errors.New("a user identifier is required to create a session"))
		return
	}

	if !utils.IsAlphaNumeric(user) {
		response.HTTPBad(w, errors.New("user must be alphanumeric"))
		return
	}

	token := session.Token{SessionID: uuid.Generate().String()}
	signedTkn, err := session.CreateSessionToken(auth.GetServerPrivateKey(), token)
	if err != nil {
		logger.Errorf("unable to create session %v", err)
		response.HTTPBad(w, err)
		return
	}

	newSession := session.Session{
		ID:        token.SessionID,
		Token:     signedTkn,
		ExpiresAt: utils.UnixToDate(utils.OneYear),
		User:      strings.ToLower(user),
	}

	if err := session.Save(newSession); err != nil {
		logger.Errorf("unable to save session %v", err)
		response.HTTPBad(w, err)
		return
	}

	response.HTTPOk(w, newSession)
	return
}

// DeleteSession will delete a session revoking any further request using that session token
func DeleteSession(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	sessionID := params["id"]

	if sessionID == "" {
		response.HTTPBad(w, errors.New("session id required"))
		return
	}

	if !session.Exist(sessionID) {
		response.HTTPBad(w, fmt.Errorf("sessions with id %s does not exist", sessionID))
		return
	}

	if err := session.Delete(sessionID); err != nil {
		logger.Errorf("unable to delete session %v", err)
		response.HTTPBad(w, err)
		return
	}

	response.HTTPOk(w, nil)
	return
}
