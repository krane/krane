package controllers

import (
	"net/http"

	"github.com/biensupernice/krane/internal/api/status"
	"github.com/biensupernice/krane/internal/session"
)

func GetSessions(w http.ResponseWriter, r *http.Request) {
	sessions, err := session.GetAllSessions()
	if err != nil {
		status.HTTPBad(w, err)
		return
	}

	status.HTTPOk(w, sessions)
	return
}
