package auth

import (
	"io/ioutil"
	"os"
	"strings"

	"github.com/sirupsen/logrus"
)

func GetServerPrivateKey() string {
	return os.Getenv("KRANE_PRIVATE_KEY")
}

func GetAuthorizeKeys() []string {
	homeDir, _ := os.UserHomeDir()
	authKeysDir := homeDir + "/.ssh/authorized_keys"

	logrus.Debugf("Reading auth keys from %s", authKeysDir)

	kBytes, err := ioutil.ReadFile(authKeysDir)
	if err != nil {
		logrus.Debug(err)
		return make([]string, 0)
	}

	// Remove trailing new line from authorized_keys file
	authKeys := strings.TrimSuffix(string(kBytes), "\n")

	// split the keys file on every new line
	return strings.Split(authKeys, "\n")
}
