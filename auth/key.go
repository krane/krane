package auth

import (
	"fmt"
	"io/ioutil"
	"os/user"

	"golang.org/x/crypto/ssh"
)

// ParsePubKey : parse public
func ParsePubKey(key []byte) (ssh.PublicKey, error) {
	pubKey, err := ssh.ParsePublicKey(key)

	if err != nil {
		return nil, err
	}

	return pubKey, nil
}

// ReadPubKeyFile : read public keys from file
func ReadPubKeyFile(pubKeyLocation string) ([]byte, error) {
	if pubKeyLocation == "" {
		homeDir := GetHomeDir()
		if homeDir == "" {
			err := fmt.Errorf("Unable to read user home dir when getting public key")
			return nil, err
		}

		pubKeyLocation = fmt.Sprintf("%s/.ssh/authorized_keys", homeDir) // Set default dir
	}

	keys, err := ioutil.ReadFile(pubKeyLocation)
	if err != nil {
		msg := fmt.Errorf("Failed to load authorized_keys - %v", err)
		return nil, msg
	}

	return keys, err
}

// GetHomeDir : get user home dir
func GetHomeDir() string {
	usr, err := user.Current()
	if err != nil {
		return ""
	}

	return usr.HomeDir
}
