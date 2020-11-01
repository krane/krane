package auth

import (
	"io/ioutil"
	"os"
	"strings"

	"github.com/sirupsen/logrus"

	"github.com/biensupernice/krane/internal/constants"
)

func GetServerPrivateKey() string {
	return os.Getenv(constants.EnvKranePrivateKey)
}

func GetAuthorizeKeys() []string {
	homeDir, _ := os.UserHomeDir()
	authKeysDir := homeDir + "/.ssh/authorized_keys"

	logrus.Debugf("Reading auth keys from %s", authKeysDir)

	kBytes, err := ioutil.ReadFile(authKeysDir)
	if err != nil {
		logrus.Debugf("unable to read auth keys from %s, %s", authKeysDir, err.Error())
		return make([]string, 0)
	}

	// Remove trailing new line from authorized_keys file
	authKeys := strings.TrimSuffix(string(kBytes), "\n")

	// split the keys file on every new line
	return strings.Split(authKeys, "\n")
}
