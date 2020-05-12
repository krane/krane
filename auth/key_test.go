package auth

import (
	"testing"
)

func TestPubKeyReadVerify(t *testing.T) {
	// Get public key
	bKey, err := ReadPubKeyFile("")
	if err != nil {
		t.Errorf("Error getting PubKey - %s", err.Error())
	}

	// Parse public key
	_, err = PubKey(bKey)
	if err != nil {
		t.Errorf("Error parsing PubKey - %s", err.Error())
	}
}
