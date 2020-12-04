package auth

import (
	"io/ioutil"
	"os"
	"strings"

	"github.com/biensupernice/krane/internal/constants"
	"github.com/biensupernice/krane/internal/logger"
)

// GetServerPrivateKey : get the private key for the Krane server
func GetServerPrivateKey() string {
	return os.Getenv(constants.EnvKranePrivateKey)
}

// GetAuthorizeKeys : get the authorized keys on the machine running Krane
func GetAuthorizeKeys() []string {
	homeDir, _ := os.UserHomeDir()
	authKeysDir := homeDir + "/.ssh/authorized_keys"

	logger.Debugf("Reading auth keys from %s", authKeysDir)

	bytes, err := ioutil.ReadFile(authKeysDir)
	if err != nil {
		logger.Debugf("unable to read auth keys from %s, %s", authKeysDir, err.Error())
		return make([]string, 0)
	}

	// remove trailing new line from authorized_keys file
	authKeys := strings.TrimSuffix(string(bytes), "\n")

	// split the keys on every new line
	return strings.Split(authKeys, "\n")
}
