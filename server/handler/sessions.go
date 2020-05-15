package handler

import (
	"encoding/json"
	"log"

	"github.com/biensupernice/krane/auth"
	"github.com/biensupernice/krane/ds"
	"github.com/biensupernice/krane/server/http"
	"github.com/gin-gonic/gin"
)

// GetSessions : get current sessions
func GetSessions(c *gin.Context) {
	mySession, _ := c.Get("session")
	log.Printf("My session: %v", mySession)

	// Get sessions
	data := ds.GetAll(auth.SessionsBucket)

	var sessions []Session
	for v := 0; v < len(data); v++ {
		var s Session
		err := json.Unmarshal(*data[v], &s)
		if err != nil {
			log.Printf("Unable to parse session [skipping] - %s", err.Error())
			continue
		}
		sessions = append(sessions, s)
	}

	http.Ok(c, sessions)
}
