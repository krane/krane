package routes

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/docker/distribution/uuid"
	"github.com/sirupsen/logrus"

	"github.com/biensupernice/krane/api/utils"
	"github.com/biensupernice/krane/internal/auth"
	"github.com/biensupernice/krane/internal/storage"
)

// AuthRequest : the payload expected when authenticating with krane
type AuthRequest struct {
	RequestID string `json:"request_id" binding:"required"`
	Token     string `json:"token" binding:"required"`
}

// AuthenticateClientJWT  : authenticate a client signed jwt token. This usually gets called after a call to /login.
// The login route returns a server phrase which is signed using the clients private auth. To authenticate the clients signed
// token this route is called and the token is validated against the clients public auth on the host machine.
func AuthenticateClientJWT(w http.ResponseWriter, r *http.Request) {
	var body AuthRequest
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		utils.HTTPBad(w, err)
		return
	}

	// Check if request id is valid, get phrase stored on the server
	serverPhrase, err := storage.Get(AuthCollection, body.RequestID)
	if err != nil {
		utils.HTTPBad(w, err)
		return
	}

	if string(serverPhrase) == "" {
		logrus.Debug("Invalid request id")
		utils.HTTPBad(w, errors.New("unable to authenticate"))
		return
	}

	authKeys := auth.GetAuthorizeKeys()
	if len(authKeys) == 0 || authKeys[0] == "" {
		logrus.Info("authorized keys not found on the server")
		utils.HTTPBad(w, errors.New("unable to authenticate"))
		return
	}

	// If any public key can be used to parse the incoming jwt token
	// the decode it, and use the phrase in the token with the one on the server
	claims := auth.VerifyAuthTokenWithAuthorizedKeys(authKeys, body.Token)
	if claims == nil || strings.Compare(string(serverPhrase), claims.Phrase) != 0 {
		utils.HTTPBad(w, errors.New("invalid token"))
		return
	}

	// Create a new token and assign it to a session
	// Remove auth data from auth bucket
	err = storage.Remove(AuthCollection, body.RequestID)
	if err != nil {
		utils.HTTPBad(w, err)
		return
	}

	sessionTkn := &auth.SessionToken{SessionID: uuid.Generate().String()}
	signedTkn, err := auth.CreateSessionToken(auth.GetServerPrivateKey(), *sessionTkn)
	if err != nil {
		utils.HTTPBad(w, err)
		return
	}

	session := &auth.Session{
		ID:        sessionTkn.SessionID,
		Token:     signedTkn,
		ExpiresAt: UnixToDate(auth.OneYear),
		Principal: "root",
	}

	err = auth.SaveSession(*session)
	if err != nil {
		utils.HTTPBad(w, err)
		return
	}

	utils.HTTPOk(w, session)
	return
}

// UnixToDate : format unix date into MM/DD/YYYY
func UnixToDate(u int64) string {
	t := time.Unix(u, 0)
	year, month, day := t.Date()
	return fmt.Sprintf("%02d/%d/%d", int(month), day, year)
}
