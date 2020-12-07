package controllers

import (
	"net/http"

	"github.com/biensupernice/krane/internal/api/response"
	"github.com/biensupernice/krane/internal/session"
)

// GetSessions : get user sessions. A session is an authenticated user with a valid token.
func GetSessions(w http.ResponseWriter, r *http.Request) {
	sessions, err := session.GetAllSessions()
	if err != nil {
		response.HTTPBad(w, err)
		return
	}

	response.HTTPOk(w, sessions)
	return
}
