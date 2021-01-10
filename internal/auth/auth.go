package auth

import (
	"errors"

	"github.com/krane/krane/internal/constants"
	"github.com/krane/krane/internal/logger"
	"github.com/krane/krane/internal/store"
)

func GetClientServerPhrase(requestID string) (string, error) {
	bytes, err := store.Client().Get(constants.AuthenticationCollectionName, requestID)
	if err != nil {
		return "", err
	}

	if bytes == nil || len(bytes) == 0 {
		logger.Warn("Invalid request id")
		return "", errors.New("unable to authenticate, invalid request id")
	}

	return string(bytes), nil
}

func RevokeClientRequestID(requestID string) error {
	return store.Client().Remove(constants.AuthenticationCollectionName, requestID)
}
