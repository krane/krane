package middleware

import (
	"log"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

var KranePrivateKey = []byte(os.Getenv("KRANE_PRIVATE_KEY"))

func TokenAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tkn := c.Request.Header.Get("Authorization")
		splitTkn := strings.Split(tkn, "Bearer")
		jwtTkn := splitTkn[1]

		log.Println(jwtTkn)

		// _, err := auth.ParseToken(KranePrivateKey, jwtTkn)

		// if err != nil {
		// 	c.JSON(http.StatusUnauthorized, err.Error())
		// 	c.Abort()
		// 	return
		// }
		c.Next()
	}
}
