package middleware

import (
	"net/http"
	"os"
	"strings"

	"github.com/biensupernice/krane/auth"
	"github.com/gin-gonic/gin"
)

// TokenAuthMiddleware : authenticate incoming http requests
func TokenAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		bearerTkn := c.Request.Header.Get("Authorization")
		splitTkn := strings.Split(bearerTkn, " ")
		jwtTkn := splitTkn[1]

		// Check type bearer token
		if strings.Compare(splitTkn[0], "Bearer") != 0 {
			c.JSON(http.StatusUnauthorized, "Unauthorized")
			c.Abort()
			return
		}

		// Private key on the server
		serverPrivKey := os.Getenv("KRANE_PRIVATE_KEY")

		// Authenticate token using server private key
		tkn := auth.ParseJWTToken(serverPrivKey, jwtTkn)
		if tkn == nil {
			c.JSON(http.StatusUnauthorized, "Unauthorized")
			c.Abort()
			return
		}

		c.Next()
	}
}
