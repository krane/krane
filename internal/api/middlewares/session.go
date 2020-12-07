package middlewares

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"

	"github.com/biensupernice/krane/internal/api/response"
	"github.com/biensupernice/krane/internal/auth"
	"github.com/biensupernice/krane/internal/logger"
	"github.com/biensupernice/krane/internal/session"
)

// ValidateSessionMiddleware : middleware to authenticate a client token against an active session
func ValidateSessionMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// grab token from headers
		tkn := r.Header.Get("Authorization")

		// then check if its a valid token
		isValidToken := isValidTokenFormat(tkn)
		if !isValidToken {
			logger.Info("Invalid token provided")
			response.HTTPBad(w, errors.New("invalid token"))
			r.Context().Done()
			return
		}

		// if its a valid token, decode the token using server private key
		pk := auth.GetServerPrivateKey()
		_, tknValue := parseToken(tkn)
		decodedTkn, err := auth.DecodeJWTToken(pk, tknValue)
		if err != nil {
			logger.Infof("Unable to decode token %s", err.Error())
			response.HTTPBad(w, err)
			r.Context().Done()
			return
		}

		// once token is decoded, parse the session token from the JWT claims
		sessionTkn, err := parseSessionTokenFromJWTClaims(decodedTkn)
		if err != nil {
			logger.Infof("Unable to parse token claims %s", err.Error())
			response.HTTPBad(w, err)
			r.Context().Done()
			return
		}

		// find the session by the id within the session token
		s, err := session.GetSessionByID(sessionTkn.SessionID)
		if err != nil {
			response.HTTPBad(w, err)
			r.Context().Done()
			return
		}

		// add the session as part of the request context
		ctx := context.WithValue(r.Context(), "session", s)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func parseSessionTokenFromJWTClaims(tkn jwt.Token) (session.Token, error) {
	claims, ok := tkn.Claims.(*auth.CustomClaims)
	if !ok {
		return session.Token{}, errors.New("unable to parse the claims for the provided token")
	}

	var sessionTkn session.Token
	bytes, _ := json.Marshal(claims.Data)
	_ = json.Unmarshal(bytes, &sessionTkn)

	return sessionTkn, nil
}

// Parse JWT token. Returns the type and value of the token
func parseToken(tkn string) (string, string) {
	if tkn == "" {
		logger.Error(errors.New("No token provided"))
		return "", ""
	}

	splitTkn := strings.Split(tkn, " ")

	tknType := splitTkn[0]
	tknValue := splitTkn[1]
	return tknType, tknValue
}

// Check if token is a well formatter Bearer token
func isValidTokenFormat(tkn string) bool {
	if tkn == "" {
		return false
	}

	// split on the space of the token ex. Bearer XXXXX
	splitTkn := strings.Split(tkn, " ")

	jwtTknType := splitTkn[0] // Bearer

	// check token is a bearer token
	if strings.Compare(jwtTknType, "Bearer") != 0 {
		logger.Debugf("Not a `Bearer` token, token type is %s", jwtTknType)
		return false
	}

	return true
}
