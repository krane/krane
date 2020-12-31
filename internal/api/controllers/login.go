package controllers

import (
	"fmt"
	"net/http"

	"github.com/docker/distribution/uuid"

	"github.com/biensupernice/krane/internal/api/response"
	"github.com/biensupernice/krane/internal/constants"
	"github.com/biensupernice/krane/internal/logger"
	"github.com/biensupernice/krane/internal/store"
)

// LoginResponse is the response received when you initially want to authenticate.
// The request_id is a uuid stored for future validation and the phrase is a generated phrased
// containing that request_id meant to be signed by the clients private key to later be unsigned
// by the clients public key to establish an authenticated sessions
type LoginResponse struct {
	RequestID string `json:"request_id"`
	Phrase    string `json:"phrase"`
}

// RequestLoginPhrase request a preliminary login request for authentication with the krane server.
// This will return a request id and phrase. The phrase should be encrypted using the clients private key.
// This route does not return a token. You must use /auth and provide the signed phrase.
func RequestLoginPhrase(w http.ResponseWriter, _ *http.Request) {
	reqID := uuid.Generate().String()
	phrase := []byte(fmt.Sprintf("Authenticating with Krane %s", reqID))

	err := store.Client().Put(constants.AuthenticationCollectionName, reqID, phrase)
	if err != nil {
		logger.Error(err)

		err = store.Client().Remove(constants.AuthenticationCollectionName, reqID)
		if err != nil {
			logger.Error(err)
			response.HTTPBad(w, err)
			return
		}
		response.HTTPBad(w, err)
		return
	}

	response.HTTPOk(w, LoginResponse{
		RequestID: reqID,
		Phrase:    string(phrase),
	})
}
