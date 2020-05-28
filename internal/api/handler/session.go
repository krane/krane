package handler

import (
	"encoding/json"
	"log"

	"github.com/biensupernice/krane/internal/api/http"
	"github.com/biensupernice/krane/internal/store"
	"github.com/gin-gonic/gin"
)

// GetSessions : currently active on the krane-server, a session consists
// of a user sucessfully authenticating and receiving a session token
func GetSessions(c *gin.Context) {
	sessionsBytes := store.GetAll(store.SessionsBucket)

	var sessions []Session
	for _, sessionBytes := range sessionsBytes {
		var s Session
		err := json.Unmarshal(*sessionBytes, &s)
		if err != nil {
			log.Printf("Unable to parse session [skipping] - %s", err.Error())
			continue
		}

		sessions = append(sessions, s)
	}

	http.Ok(c, sessions)
}
