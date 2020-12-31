package controllers

import (
	"net/http"

	"github.com/biensupernice/krane/internal/api/response"
	"github.com/biensupernice/krane/internal/session"
)

// GetSessions returns a list of user sessions. A session is an authenticated client with a valid token.
func GetSessions(w http.ResponseWriter, _ *http.Request) {
	sessions, err := session.GetAllSessions()
	if err != nil {
		response.HTTPBad(w, err)
		return
	}

	response.HTTPOk(w, sessions)
	return
}
