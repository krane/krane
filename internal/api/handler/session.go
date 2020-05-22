package handler

import (
	"encoding/json"
	"log"

	"github.com/biensupernice/krane/internal/api/http"
	"github.com/biensupernice/krane/internal/data"
	"github.com/gin-gonic/gin"
)

// GetSessions : get current sessions
func GetSessions(c *gin.Context) {
	sessionData := data.GetAll(data.SessionsBucket)

	var sessions []Session
	for v := 0; v < len(sessionData); v++ {
		var s Session
		err := json.Unmarshal(*sessionData[v], &s)
		if err != nil {
			log.Printf("Unable to parse session [skipping] - %s", err.Error())
			continue
		}
		sessions = append(sessions, s)
	}

	http.Ok(c, sessions)
}
