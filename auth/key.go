package auth

import (
	"fmt"
	"io/ioutil"
	"os/user"

	"golang.org/x/crypto/ssh"
)

// ParsePubKey : parse public
func ParsePubKey(key []byte) (ssh.PublicKey, error) {
	newAuthorizedKey, _, _, _, err := ssh.ParseAuthorizedKey(key)
	if err != nil {
		return nil, err
	}

	return newAuthorizedKey, nil
}

// ReadPubKeyFile : read public keys from file
func ReadPubKeyFile(dir string) ([]byte, error) {
	if dir == "" {
		homeDir := getHomeDir()
		if homeDir == "" {
			err := fmt.Errorf("Unable to read user home dir when getting public key")
			return nil, err
		}

		dir = fmt.Sprintf("%s/.ssh/id_rsa.pub", homeDir) // Set default dir
	}

	byteKey, err := ioutil.ReadFile(dir)
	if err != nil {
		return nil, err
	}

	return byteKey, nil
}

func getHomeDir() string {
	usr, err := user.Current()
	if err != nil {
		return ""
	}

	return usr.HomeDir
}
