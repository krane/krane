package auth

import (
	"errors"
	"fmt"

	"github.com/docker/distribution/uuid"

	"github.com/krane/krane/internal/constants"
	"github.com/krane/krane/internal/store"
)

// GetAuthenticationPhrase returns the generate phrase for a given request id
func GetAuthenticationPhrase(requestID string) (string, error) {
	bytes, err := store.Client().Get(constants.AuthenticationCollectionName, requestID)
	if err != nil {
		return "", err
	}

	if bytes == nil || len(bytes) == 0 {
		return "", errors.New("invalid request id")
	}

	return string(bytes), nil
}

// CreateAuthenticationPhrase returns a request id (uuid) and a phrase used by the client for authentication
func CreateAuthenticationPhrase() (string, string, error) {
	reqID := uuid.Generate().String()
	phrase := []byte(fmt.Sprintf("Krane authentication request id: %s", reqID))

	if err := store.Client().Put(constants.AuthenticationCollectionName, reqID, phrase); err != nil {
		if err := store.Client().Remove(constants.AuthenticationCollectionName, reqID); err != nil {
			return "", "", err
		}
		return "", "", err
	}

	return reqID, string(phrase), nil
}

// RevokeAuthenticationRequest removes the request from the authentication collection
func RevokeAuthenticationRequest(requestID string) error {
	return store.Client().Remove(constants.AuthenticationCollectionName, requestID)
}
