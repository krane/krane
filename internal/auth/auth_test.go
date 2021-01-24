package auth

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/krane/krane/internal/utils/test"
)

func TestMain(m *testing.M) {
	test.SetupDb()

	code := m.Run()

	test.TeardownDb()
	os.Exit(code)
}

func TestGetAuthenticationPhrase(t *testing.T) {
	reqID, phrase, err := CreateAuthenticationPhrase()
	assert.Nil(t, err)

	phrase2, err := GetAuthenticationPhrase(reqID)
	assert.Nil(t, err)
	assert.Equal(t, phrase, phrase2)
}

func TestGetAuthenticationPhraseReturnsErrWhenMissingRequestID(t *testing.T) {
	phrase, err := GetAuthenticationPhrase("0")
	assert.Error(t, err, "invalid request id")
	assert.Empty(t, phrase)
}

func TestCreateAuthenticationPhrase(t *testing.T) {
	reqID, phrase, err := CreateAuthenticationPhrase()
	assert.Nil(t, err)
	assert.NotEmpty(t, reqID)
	assert.NotEmpty(t, phrase)
}

func TestRevokeAuthenticationRequest(t *testing.T) {
	reqID, _, err := CreateAuthenticationPhrase()
	assert.Nil(t, err)

	err = RevokeAuthenticationRequest(reqID)
	assert.Nil(t, err)

	phrase, err := GetAuthenticationPhrase(reqID)
	assert.Error(t, err, "invalid request id")
	assert.Empty(t, phrase)
}
