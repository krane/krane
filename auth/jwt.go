package auth

import (
	"fmt"
	"log"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var (
	jwtKey = []byte("private.key")
	Bucket = "AuthBucket"

	OneYear = time.Now().Add(time.Minute * 525600).Unix()
)

type Claims struct {
	Phrase string `json:"phrase"`
	jwt.StandardClaims
}

func CreateToken(phrase string) (string, error) {
	// Create the JWT claims
	claims := &Claims{
		Phrase: phrase,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: OneYear,
			Issuer:    "krane-server",
		},
	}

	// Declare the token with the algorithm used for signing, and the claims
	tkn := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with our secret
	tknStr, err := tkn.SignedString(jwtKey)
	if err != nil {
		return "", err
	}

	return tknStr, nil
}

func ValidateToken(pubKey string, tknStr string, phrase string) bool {
	// Initialize a new instance of `Claims`
	claims := &Claims{}

	// Parse the JWT string and store the result in `claims`.
	// Note that we are passing the key in this method as well. This method will return an error
	// if the token is invalid (if it has expired according to the expiry time we set on sign in),
	// or if the signature does not match
	tkn, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil {
		msg := fmt.Sprintf("Unable to verify token - %s", err.Error())
		log.Printf(msg)
		return false
	}

	// Verify token is valid
	if !tkn.Valid {
		msg := fmt.Sprintf("Invalid token - %s", err.Error())
		log.Printf(msg)
		return false
	}

	return true
}
