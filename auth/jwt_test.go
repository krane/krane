package auth

import (
	"encoding/json"
	"log"
	"testing"
)

var (
	phrase = "Hello krane tests!"

	sKey = []byte("biensupernice")

	dummyTknNoPayload   = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJwYXlsb2FkIjpudWxsLCJleHAiOjE2MjA3NzA5OTUsImlzcyI6ImtyYW5lLXNlcnZlciJ9.N8-Y5P8loK062zllCUWo7Duq52xx-tJk5ezRPvs7Rmw"
	dummyTknWithPayload = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJwYXlsb2FkIjoiZXlKd2FISmhjMlVpT2lKSVpXeHNieUJyY21GdVpTQjBaWE4wY3lFaWZRPT0iLCJleHAiOjE2MjA3NzQyNTcsImlzcyI6ImtyYW5lLXNlcnZlciJ9.0zufLuG8oC2LmHW3P2COWkP02NE3CO7Q24JfQw0KoG8"
)

type TokenPayload struct {
	Phrase string `json:"phrase"`
}

func TestCreateToken(t *testing.T) {
	var data []byte
	data, _ = json.Marshal(&TokenPayload{Phrase: phrase})

	tkn, err := CreateToken(sKey, data)  // With payload
	tkn2, err2 := CreateToken(sKey, nil) // Without payload
	tkn3, err3 := CreateToken(nil, nil)  // Without key or payload

	log.Printf("My token - %s", tkn)
	log.Printf("Payload - %s", data)

	// Assert no error
	if err != nil {
		t.Errorf("Unable to create token - %s", err.Error())
	}

	if err2 != nil {
		t.Errorf("Unable to create token - %s", err.Error())
	}

	// Assert error
	if err3 == nil {
		t.Errorf("%s", err.Error())
	}

	// Assert token is not empty
	if tkn == "" {
		t.Errorf("Unable to create token - %s", err.Error())
	}

	if tkn2 == "" {
		t.Errorf("Unable to create token - %s", err.Error())
	}

	if tkn3 != "" {
		t.Errorf("Expected token to be empty")
	}
}

func TestCreateTokenFailsWhenEmptyPrivKey(t *testing.T) {
	tkn, err := CreateToken(nil, nil)

	// Assert error
	if err == nil {
		t.Errorf("Expected error when creating token with no sign key")
	}

	// Assert empty token
	if tkn != "" {
		t.Errorf("Expected error when creating token with no sign key")
	}

	expMsg := "Cannot create token - signing key not provided"
	if err.Error() != expMsg {
		t.Errorf("Expected error: `%s` when creating token with no signing key", expMsg)
	}
}

func TestParseTokenWithNoPayload(t *testing.T) {
	bytes, err := ParseToken(sKey, dummyTknNoPayload)

	if err != nil {
		t.Errorf("Unable to validate token")
	}

	payload := getTokenPayload(bytes)

	if payload.Phrase != "" {
		t.Error("Expected empty payload when parsing token")
	}
}

func TestParseTokenWithPayload(t *testing.T) {
	bytes, err := ParseToken(sKey, dummyTknWithPayload)

	if err != nil {
		t.Errorf("Unable to parse token")
	}

	payload := getTokenPayload(bytes)

	if payload.Phrase == "" {
		t.Error("Expected payload to contain values")
	}

	if payload.Phrase != phrase {
		t.Errorf("Expected phrase to be `%s`, instead got `%s`", phrase, payload.Phrase)
	}
}

func getTokenPayload(bytes []byte) TokenPayload {
	data := TokenPayload{}
	json.Unmarshal(bytes, &data)
	return data
}
