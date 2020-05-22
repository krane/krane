package middleware

import (
	"encoding/json"
	"errors"
	"os"
	"strings"

	"github.com/biensupernice/krane/internal/api/handler"
	"github.com/biensupernice/krane/internal/api/http"
	"github.com/biensupernice/krane/internal/auth"
	"github.com/biensupernice/krane/internal/logger"
	"github.com/biensupernice/krane/internal/store"
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
			logger.Debug("No token provided")
			http.Unauthorized(c)
			return
		}

		// Split on the space of the token ex. Bearer XXXXX
		splitTkn := strings.Split(bearerTkn, " ")
		jwtTkn := splitTkn[1]

		// Check token is a bearer token
		if strings.Compare(splitTkn[0], "Bearer") != 0 {
			logger.Debug("Not a `Bearer` token")
			msg := errors.New("Invalid token")
			http.BadRequest(c, msg)
			return
		}

		// Authenticate token using server private key
		tkn := auth.ParseJWTToken(serverPrivKey, jwtTkn)
		if tkn == nil {
			logger.Debug("Unable to authenticate token with server private key")
			http.Unauthorized(c)
			return
		}

		// Parse token claims into custom claims
		customClaims, ok := tkn.Claims.(*auth.CustomClaims)
		if !ok {
			logger.Debug("Unable to parse the claims for the provided token")
			http.Unauthorized(c)
			return
		}

		// The data inside custom claims should be of type SessionToken
		// Convert custom claims data into SessionToken
		dataBytes, _ := json.Marshal(customClaims.Data)
		var sessionTkn handler.SessionToken
		err := json.Unmarshal(dataBytes, &sessionTkn)
		if err != nil {
			logger.Debug("Unable to convert custom claims into a session token")
			http.Unauthorized(c)
			return
		}

		// Check if session is valid by retrieving the sessions from the servers datastore
		sessionBytes := store.Get(store.SessionsBucket, sessionTkn.SessionID)

		var session handler.Session
		err = json.Unmarshal(sessionBytes, &session) // convert bytes to session struct
		if err != nil {
			logger.Debug("Unable to convert token from the store into a session token")
			http.Unauthorized(c)
			return
		}

		// Compare if session token from the server matches the incoming bearer token
		if session.Token == "" || strings.Compare(jwtTkn, session.Token) != 0 {
			logger.Debug("Token did not match against the server token - try loggin in again")
			http.Unauthorized(c)
			return
		}

		// Add session to router context
		c.Set("session", session)

		// Continue to the next handler
		c.Next()
	}
}
