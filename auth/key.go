package auth

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os/user"
	"strings"

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

// GetAuthorizedKeys : get authorized keys from server
func GetAuthorizedKeys(authKeysLocation string) ([]string, error) {
	// sets to ~/.ssh/authorized_keys is not passed in
	if authKeysLocation == "" {
		// Get user home dir
		usrHomeDir := GetHomeDir()
		if usrHomeDir == "" {
			err := fmt.Errorf("Unable to read user home dir when getting public keys")
			return nil, err
		}

		// Format location of authorized_keys using users home directory as base
		authKeysLocation = fmt.Sprintf("%s/.ssh/authorized_keys", usrHomeDir)
	}

	// Read authorized_keys from authKeysLocation
	authKeysBytes, err := ioutil.ReadFile(authKeysLocation)
	if err != nil {
		msg := fmt.Errorf("Failed to load authorized_keys - %v", err)
		return nil, msg
	}

	// Remove trailing new line from authorized_keys file
	authKeys := strings.TrimSuffix(string(authKeysBytes), "\n")

	// Every token is a single line, split the authorized_keys file on every new line returning array of authorized_keys
	authKeysArr := strings.Split(authKeys, "\n")

	return authKeysArr, nil
}

// GetHomeDir : get user home dir
func GetHomeDir() string {
	usr, err := user.Current()
	if err != nil {
		return ""
	}

	return usr.HomeDir
}

// VerifyAuthTokenWithAuthorizedKeys : get auth claims from jwt token using an authorized key from server
func VerifyAuthTokenWithAuthorizedKeys(authorizedKeys []string, authTkn string) (*AuthClaims, error) {
	// Validate if pub key can parse incoming token
	var authClaims *AuthClaims
	for currKey := 0; currKey < len(authorizedKeys); currKey++ {
		// Parse token against curr key
		c, err := ParseToken(authorizedKeys[currKey], authTkn)
		if err != nil {
			continue
		}

		// If parsing was succesful, map jwt claims into authclaims
		jwtClaims, ok := c.(*AuthClaims)
		if !ok {
			continue
		}

		authClaims = jwtClaims
		break
	}

	// Veirfy a token was found and authClaims is not empty, auth claims should have server token
	if authClaims == nil {
		msg := "Unable to verify with public key, make sure to have your public key in authorized_keys on the server"
		return nil, errors.New(msg)
	}

	return authClaims, nil
}
