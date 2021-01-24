package middlewares

import (
	"context"
	"errors"
	"net/http"

	"github.com/krane/krane/internal/api/response"
	"github.com/krane/krane/internal/auth"
	"github.com/krane/krane/internal/logger"
	"github.com/krane/krane/internal/session"
)

// ValidateSessionMiddleware middleware to authenticate a client token against an active session
func ValidateSessionMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// grab token from headers
		tkn := r.Header.Get("Authorization")

		// then check if its a valid Bearer token
		isValidToken := session.IsValidTokenFormat(tkn)
		if !isValidToken {
			logger.Info("Invalid token provided")
			response.HTTPBad(w, errors.New("invalid token"))
			r.Context().Done()
			return
		}

		// if its a valid token, decode the token using server private key
		pk := auth.GetServerPrivateKey()
		_, tknValue := session.ParseTokenTypeAndValue(tkn)
		decodedTkn, err := session.DecodeJWTToken(pk, tknValue)
		if err != nil {
			logger.Infof("Unable to decode token %s", err.Error())
			response.HTTPBad(w, err)
			r.Context().Done()
			return
		}

		// once token is decoded, parse the session token from the JWT claims
		sessionTkn, err := session.ParseSessionTokenFromJWTClaims(decodedTkn)
		if err != nil {
			logger.Infof("Unable to parse token claims %s", err.Error())
			response.HTTPBad(w, err)
			r.Context().Done()
			return
		}

		// find the session by the id, the id is inside the session token we just decoded
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
