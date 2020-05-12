package auth

import (
	"testing"
)

func TestPubKeyReadVerify(t *testing.T) {
	// Get public key
	pubKey, err := ReadPubKeyFile("")
	if err != nil {
		t.Errorf("Error getting PubKey - %s", err.Error())
	}

	// Parse public key
	_, err = ParsePubKey(pubKey)
	if err != nil {
		t.Errorf("Error parsing PubKey - %s", err.Error())
	}
}
