package routes

import (
	"fmt"
	"net/http"

	"github.com/docker/distribution/uuid"
	"github.com/sirupsen/logrus"

	"github.com/biensupernice/krane/api/utils"
	"github.com/biensupernice/krane/internal/storage"
)

const (
	AuthCollection = "Auth"
)

// RequestLoginPhrase : request a preliminary login request for authentication with the krane server.
// This will return a request id and phrase. The phrase should be encrypted using the clients private auth.
// This route does not return a token. You must use /auth and provide the signed phrase.
func RequestLoginPhrase(w http.ResponseWriter, r *http.Request) {
	reqID := uuid.Generate().String()
	phrase := []byte(fmt.Sprintf("Authenticating with krane %s", reqID))

	err := storage.Put(AuthCollection, reqID, phrase)
	if err != nil {
		logrus.Error(err)

		err = storage.Remove(AuthCollection, reqID)
		if err != nil {
			logrus.Error(err)
			utils.HTTPBad(w, err)
			return
		}
		utils.HTTPBad(w, err)
		return
	}

	utils.HTTPOk(w, struct {
		RequestID string `json:"request_id"`
		Phrase    string `json:"phrase"`
	}{
		RequestID: reqID,
		Phrase:    string(phrase),
	})
}
