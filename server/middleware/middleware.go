package middleware

import (
	"encoding/json"
	"net/http"
	"os"
	"strings"

	"github.com/biensupernice/krane/auth"
	"github.com/biensupernice/krane/ds"
	"github.com/biensupernice/krane/server/handler"
	"github.com/gin-gonic/gin"
)

// AuthSessionMiddleware : validate a session bearer token from incoming http request
func AuthSessionMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get Bearer token
		bearerTkn := c.Request.Header.Get("Authorization")

		// Check token is provided
		if bearerTkn == "" {
			c.JSON(http.StatusUnauthorized, "Unauthorized")
			c.Abort()
			return
		}

		// Split on the space of the token ex. Bearer XXXXX
		splitTkn := strings.Split(bearerTkn, " ")
		jwtTkn := splitTkn[1]

		// Check token type is bearer token
		if strings.Compare(splitTkn[0], "Bearer") != 0 {
			c.JSON(http.StatusUnauthorized, "Unauthorized")
			c.Abort()
			return
		}

		// Server private key
		serverPrivKey := os.Getenv("KRANE_PRIVATE_KEY")

		// Authenticate token using server private key
		tkn := auth.ParseJWTToken(serverPrivKey, jwtTkn)
		if tkn == nil {
			c.JSON(http.StatusUnauthorized, "Unauthorized")
			c.Abort()
			return
		}

		// Parse token claims into our custom claims
		customClaims, ok := tkn.Claims.(*auth.CustomClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, "Invalid token")
			c.Abort()
			return
		}

		// THe data inside the custom claims should be if type SessionToken
		// Convert custom claims into SessionToken type
		dataBytes, _ := json.Marshal(customClaims.Data)
		var sessionTkn handler.SessionToken
		err := json.Unmarshal(dataBytes, &sessionTkn)
		if err != nil {
			c.JSON(http.StatusUnauthorized, "Invalid token")
			c.Abort()
			return
		}

		// Validate if session is still valid by retrieving the session from the datastore
		sessionBytes := ds.Get(auth.SessionsBucket, sessionTkn.SessionID) // get session from db
		var session handler.Session
		err = json.Unmarshal(sessionBytes, &session) // convert bytes to session struct
		if err != nil {
			c.JSON(http.StatusUnauthorized, "Invalid token")
			c.Abort()
			return
		}

		// Compare if session token from the server matches the incoming jwt token
		if session.Token == "" || strings.Compare(jwtTkn, session.Token) != 0 {
			c.JSON(http.StatusUnauthorized, "Unauthorized")
			c.Abort()
			return
		}

		// Add session struct into router context
		c.Set("session", session)

		// Continue to the next handler
		c.Next()
	}
}
