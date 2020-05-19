package handler

import (
	"encoding/json"
	"log"

	"github.com/biensupernice/krane/data"
	"github.com/biensupernice/krane/server/http"
	"github.com/gin-gonic/gin"
)

// GetSessions : get current sessions
func GetSessions(c *gin.Context) {
	mySession, _ := c.Get("session")
	log.Printf("My session: %v", mySession)

	// Get sessions
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

	http.Created(c, sessions)
}
