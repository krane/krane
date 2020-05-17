package middleware

import (
	"encoding/json"
	"errors"
	"os"
	"strings"

	"github.com/biensupernice/krane/auth"
	"github.com/biensupernice/krane/ds"
	"github.com/biensupernice/krane/server/handler"
	"github.com/biensupernice/krane/server/http"
	"github.com/gin-gonic/gin"
)

// Get server private key
var (
	serverPrivKey = os.Getenv("KRANE_PRIVATE_KEY")
)

// AuthSessionMiddleware : validate a session bearer token from incoming http request
func AuthSessionMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get token from headers
		bearerTkn := c.Request.Header.Get("Authorization")

		// Check token is provided
		if bearerTkn == "" {
			msg := errors.New("Token not provided")
			http.BadRequest(c, msg)
			return
		}

		// Split on the space of the token ex. Bearer XXXXX
		splitTkn := strings.Split(bearerTkn, " ")
		jwtTkn := splitTkn[1]

		// Check token is a bearer token
		if strings.Compare(splitTkn[0], "Bearer") != 0 {
			msg := errors.New("Invalid token")
			http.BadRequest(c, msg)
			return
		}

		// Authenticate token using server private key
		tkn := auth.ParseJWTToken(serverPrivKey, jwtTkn)
		if tkn == nil {
			http.Unauthorized(c)
			return
		}

		// Parse token claims into custom claims
		customClaims, ok := tkn.Claims.(*auth.CustomClaims)
		if !ok {
			http.Unauthorized(c)
			return
		}

		// The data inside custom claims should be of type SessionToken
		// Convert custom claims data into SessionToken
		dataBytes, _ := json.Marshal(customClaims.Data)
		var sessionTkn handler.SessionToken
		err := json.Unmarshal(dataBytes, &sessionTkn)
		if err != nil {
			http.Unauthorized(c)
			return
		}

		// Check if session is valid by retrieving the sessions from the servers datastore
		sessionBytes := ds.Get(auth.SessionsBucket, sessionTkn.SessionID)

		var session handler.Session
		err = json.Unmarshal(sessionBytes, &session) // convert bytes to session struct
		if err != nil {
			http.Unauthorized(c)
			return
		}

		// Compare if session token from the server matches the incoming bearer token
		if session.Token == "" || strings.Compare(jwtTkn, session.Token) != 0 {
			http.Unauthorized(c)
			return
		}

		// Add session to router context
		c.Set("session", session)

		// Continue to the next handler
		c.Next()
	}
}
