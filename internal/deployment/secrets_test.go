package deployment

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/krane/krane/internal/store"
	"github.com/krane/krane/internal/utils"
)

const boltpath = "./krane.db"
const testDeployment = "krane-test"

func teardown() { os.Remove(boltpath) }

func TestMain(m *testing.M) {
	store.Connect((boltpath))
	defer store.Client().Disconnect()

	code := m.Run()

	teardown()
	os.Exit(code)
}

func TestAddNewSecret(t *testing.T) {
	s1, err := AddSecret(testDeployment, "token", "biensupernice")
	assert.Nil(t, err)
	assert.Equal(t, "token", s1.Key)
	assert.Equal(t, "biensupernice", s1.Value)
	assert.Equal(t, "@TOKEN", s1.Alias)

	s2, err := AddSecret(testDeployment, "api_token", "biensupernice")
	assert.Nil(t, err)
	assert.Equal(t, "@API_TOKEN", s2.Alias)

	s3, err := AddSecret(testDeployment, "api-token", "biensupernice")
	assert.Nil(t, err)
	assert.Equal(t, "@API_TOKEN", s3.Alias)

	s4, err := AddSecret(testDeployment, "api-token123", "biensupernice")
	assert.Nil(t, err)
	assert.Equal(t, "@API_TOKEN123", s4.Alias)

	s5, err := AddSecret(testDeployment, "API_PORT_8080", "8080")
	assert.Nil(t, err)
	assert.Equal(t, "@API_PORT_8080", s5.Alias)

	s6, err := AddSecret(testDeployment, "API-PORT-8080", "8080")
	assert.Nil(t, err)
	assert.Equal(t, "@API_PORT_8080", s6.Alias)

	s7, err := AddSecret(testDeployment, "env", "dev")
	assert.Nil(t, err)
	assert.Equal(t, "@ENV", s7.Alias)

	s8, err := AddSecret(testDeployment, "8080_API_PORT", "8080")
	assert.Nil(t, err)
	assert.Equal(t, "@8080_API_PORT", s8.Alias)

	s9, err := AddSecret(testDeployment, "8080-API-PORT", "8080")
	assert.Nil(t, err)
	assert.Equal(t, "@8080_API_PORT", s9.Alias)

	s10, err := AddSecret(testDeployment, "aPi_ToKeN-1337", "8080")
	assert.Nil(t, err)
	assert.Equal(t, "@API_TOKEN_1337", s10.Alias)
}

func TestFormatSecretCollectionName(t *testing.T) {
	collections := []string{"api", "UI", "api-proxy", "messaging_service", "db-container", "app-123-proxy", "123-proxy_api", "aPi_pR0Xy"}
	for _, collection := range collections {
		expected := fmt.Sprintf("%s-secrets", strings.ToLower(collection)) // lowercase, ending with -secrets
		assert.Equal(t, expected, getSecretsCollectionName(collection))
	}
}

func TestErrorWhenInvalidSecretKey(t *testing.T) {
	assert.False(t, isValidSecretKey(""))
	assert.False(t, isValidSecretKey("X"))
	assert.False(t, isValidSecretKey("X"))
	assert.False(t, isValidSecretKey(utils.RandomString(51)))

	illegalChars := []string{"!", "@", "#", "$", "%", "&", "*", ",", ".", "'", ":", "/", "\"", "=", "+", "?", ">", "<", "|", "}", "{", "-", "_", "^", "(", ")", "[", "]"}
	for _, char := range illegalChars {
		assert.False(t, isValidSecretKey(fmt.Sprintf("%sTOKEN", char)), char)
		assert.False(t, isValidSecretKey(fmt.Sprintf("TOKEN%s", char)), char)
	}
}

func TestSuccessWhenValidSecretKey(t *testing.T) {
	assert.True(t, isValidSecretKey("OS"))
	assert.True(t, isValidSecretKey(utils.RandomString(2)))
	assert.True(t, isValidSecretKey(utils.RandomString(20)))
}

func TestRedactSecret(t *testing.T) {
	s := Secret{
		Deployment: testDeployment,
		Key:        "SECRET_TOKEN",
		Value:      "biensupernice",
		Alias:      "@SECRET_TOKEN",
	}
	s.Redact()
	assert.Equal(t, "<redacted>", s.Value)
	assert.NotEqual(t, "biensupernice", s.Value)
}

func TestGetSecretsByDeployment(t *testing.T) {
	secretKey := utils.RandomString(20)
	secretValue := utils.RandomString(20)

	newSecret, err := AddSecret(testDeployment, secretKey, secretValue)
	assert.Nil(t, err)

	secrets, err := GetAllSecrets(testDeployment)
	assert.Nil(t, err)
	assert.True(t, len(secrets) > 0)

	var s Secret
	for _, secret := range secrets {
		if secret.Key == newSecret.Key {
			s = *secret
			break
		}
	}

	assert.NotNil(t, s)
	assert.Equal(t, s.Key, secretKey)
	assert.Equal(t, s.Value, secretValue)
}

func TestGetSecret(t *testing.T) {
	secretKey := utils.RandomString(20)
	secretValue := utils.RandomString(20)

	secr, err := AddSecret(testDeployment, secretKey, secretValue)
	assert.Nil(t, err)

	s, err := GetSecret(testDeployment, secr.Key)
	assert.Nil(t, err)

	assert.NotNil(t, s)
	assert.Equal(t, s.Key, secretKey)
	assert.Equal(t, s.Value, secretValue)
}

func TestErrorWhenGetSecretByNonExistingAlias(t *testing.T) {
	_, err := GetSecret(testDeployment, "non-existing-key")
	assert.NotNil(t, err)
	assert.Equal(t, "secret with key non-existing-key not found for deployment krane-test", err.Error())
}

func TestDeleteSecret(t *testing.T) {
	secretKey := utils.RandomString(20)
	secretValue := utils.RandomString(20)

	// add
	_, err := AddSecret(testDeployment, secretKey, secretValue)
	assert.Nil(t, err)

	// get
	s, err := GetSecret(testDeployment, secretKey)
	assert.Nil(t, err)

	assert.NotNil(t, s)
	assert.Equal(t, s.Key, secretKey)
	assert.Equal(t, s.Value, secretValue)

	// delete
	err = DeleteSecret(s.Deployment, s.Key)
	assert.Nil(t, err)

	// get
	_, err = GetSecret(testDeployment, secretKey)
	assert.NotNil(t, err)
	assert.Equal(t, fmt.Sprintf("secret with key %s not found for deployment %s", secretKey, testDeployment), err.Error())
}
