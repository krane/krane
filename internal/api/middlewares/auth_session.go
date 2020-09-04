package middlewares

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/sirupsen/logrus"

	"github.com/biensupernice/krane/internal/api/status"
	"github.com/biensupernice/krane/internal/auth"
)

// AuthSessionMiddleware : middleware to authenticate a client token against an active session
func AuthSessionMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logrus.Infof("Authenticating session")
		// Get token from headers
		tkn := r.Header.Get("Authorization")

		isValidToken := isValidTokenFormat(tkn)
		if !isValidToken {
			logrus.Infof("Invalid token provided")
			status.HTTPBad(w, errors.New("invalid token"))
			r.Context().Done()
			return
		}

		// Authenticate token using server private key
		pk := auth.GetServerPrivateKey()
		_, tknValue := parseToken(tkn)
		decodedTkn, err := auth.DecodeJWTToken(pk, tknValue)
		if err != nil {
			logrus.Infof("Unable to decode token", err.Error())
			status.HTTPBad(w, err)
			r.Context().Done()
			return
		}

		// Parse token claims into custom claims
		sessionTkn, err := parseSessionTokenFromJWTClaims(decodedTkn)
		if err != nil {
			logrus.Infof("Unable to parse token claims", err.Error())
			status.HTTPBad(w, err)
			r.Context().Done()
			return
		}

		session, err := auth.GetSessionByID(sessionTkn.SessionID)
		if err != nil {
			logrus.Infof("Unable to find a valid session", err.Error())
			status.HTTPBad(w, err)
			r.Context().Done()
			return
		}

		// Add the session as part of the request context
		ctx := context.WithValue(r.Context(), "session", session)
		next.ServeHTTP(w, r.WithContext(ctx))
	})

}

func parseSessionTokenFromJWTClaims(tkn jwt.Token) (auth.SessionToken, error) {
	claims, ok := tkn.Claims.(*auth.CustomClaims)
	if !ok {
		return auth.SessionToken{}, errors.New("unable to parse the claims for the provided token")
	}

	var sessionTkn auth.SessionToken
	bytes, _ := json.Marshal(claims.Data)
	_ = json.Unmarshal(bytes, &sessionTkn)
	return sessionTkn, nil
}

// Parse JWT token. Returns the type and value of the token
func parseToken(tkn string) (string, string) {
	if tkn == "" {
		logrus.Error("No token provided")
		return "", ""
	}

	splitTkn := strings.Split(tkn, " ")
	return splitTkn[0], splitTkn[1]
}

// Check if token is a well formatter Bearer token
func isValidTokenFormat(tkn string) bool {
	if tkn == "" {
		logrus.Info("No token provided")
		return false
	}

	// Split on the space of the token ex. Bearer XXXXX
	splitTkn := strings.Split(tkn, " ")

	jwtTknType := splitTkn[0]

	// Check token is a bearer token
	if strings.Compare(jwtTknType, "Bearer") != 0 {
		logrus.Info("Not a `Bearer` token")
		return false
	}

	return true
}
