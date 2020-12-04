package secrets

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/biensupernice/krane/internal/constants"
	"github.com/biensupernice/krane/internal/deployment/kconfig"
	"github.com/biensupernice/krane/internal/store"
	"github.com/biensupernice/krane/internal/utils"
)

const boltpath = "./krane.db"
const testNamespace = "krane-test"

func teardown() { os.Remove(boltpath) }

func TestMain(m *testing.M) {
	store.Connect((boltpath))
	defer store.Client().Disconnect()

	// Create deployment (namespace)
	deployment := kconfig.Kconfig{Name: testNamespace}
	bytes, _ := deployment.Serialize()
	store.Client().Put(constants.DeploymentsCollectionName, deployment.Name, bytes)

	code := m.Run()

	teardown()
	os.Exit(code)
}

func TestAddNewSecret(t *testing.T) {
	s1, err := Add(testNamespace, "token", "biensupernice")
	assert.Nil(t, err)
	assert.Equal(t, "token", s1.Key)
	assert.Equal(t, "biensupernice", s1.Value)
	assert.Equal(t, "@TOKEN", s1.Alias)

	s2, err := Add(testNamespace, "api_token", "biensupernice")
	assert.Nil(t, err)
	assert.Equal(t, "@API_TOKEN", s2.Alias)

	s3, err := Add(testNamespace, "api-token", "biensupernice")
	assert.Nil(t, err)
	assert.Equal(t, "@API_TOKEN", s3.Alias)

	s4, err := Add(testNamespace, "api-token123", "biensupernice")
	assert.Nil(t, err)
	assert.Equal(t, "@API_TOKEN123", s4.Alias)

	s5, err := Add(testNamespace, "API_PORT_8080", "8080")
	assert.Nil(t, err)
	assert.Equal(t, "@API_PORT_8080", s5.Alias)

	s6, err := Add(testNamespace, "API-PORT-8080", "8080")
	assert.Nil(t, err)
	assert.Equal(t, "@API_PORT_8080", s6.Alias)

	s7, err := Add(testNamespace, "env", "dev")
	assert.Nil(t, err)
	assert.Equal(t, "@ENV", s7.Alias)

	s8, err := Add(testNamespace, "8080_API_PORT", "8080")
	assert.Nil(t, err)
	assert.Equal(t, "@8080_API_PORT", s8.Alias)

	s9, err := Add(testNamespace, "8080-API-PORT", "8080")
	assert.Nil(t, err)
	assert.Equal(t, "@8080_API_PORT", s9.Alias)

	s10, err := Add(testNamespace, "aPi_ToKeN-1337", "8080")
	assert.Nil(t, err)
	assert.Equal(t, "@API_TOKEN_1337", s10.Alias)

}

func TestErrorWhenAddSecretToNonExistingDeployment(t *testing.T) {
	_, err1 := Add("non-existing-namespace", "TOKEN", "biensupernice")
	assert.Error(t, err1)
	assert.Equal(t, "unable to find namespace non-existing-namespace", err1.Error())

	_, err2 := Add("", "TOKEN", "biensupernice")
	assert.Error(t, err2)
	assert.Equal(t, "unable to find namespace ", err2.Error())

}

func TestFormatSecretCollectionName(t *testing.T) {
	collections := []string{"api", "UI", "api-proxy", "messaging_service", "db-container", "app-123-proxy", "123-proxy_api", "aPi_pR0Xy"}
	for _, collection := range collections {
		expected := fmt.Sprintf("%s-secrets", strings.ToLower(collection)) // lowercase, ending with -secrets
		assert.Equal(t, expected, getNamespaceCollectionName(collection))
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
		Namespace: testNamespace,
		Key:       "SECRET_TOKEN",
		Value:     "biensupernice",
		Alias:     "@SECRET_TOKEN",
	}
	s.Redact()
	assert.Equal(t, "<redacted>", s.Value)
	assert.NotEqual(t, "biensupernice", s.Value)
}

func TestGetSecretsByNamespace(t *testing.T) {
	secretKey := utils.RandomString(20)
	secretValue := utils.RandomString(20)

	secr, err := Add(testNamespace, secretKey, secretValue)
	assert.Nil(t, err)

	secrets, err := GetAll(testNamespace)
	assert.Nil(t, err)
	assert.True(t, len(secrets) > 0)

	var s Secret
	for _, secret := range secrets {
		if secret.Key == secr.Key {
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

	secr, err := Add(testNamespace, secretKey, secretValue)
	assert.Nil(t, err)

	s, err := Get(testNamespace, secr.Key)
	assert.Nil(t, err)

	assert.NotNil(t, s)
	assert.Equal(t, s.Key, secretKey)
	assert.Equal(t, s.Value, secretValue)
}

func TestErrorWhenGetSecretByNonExistingAlias(t *testing.T) {
	_, err := Get(testNamespace, "non-existing-key")
	assert.NotNil(t, err)
	assert.Equal(t, "secret with key non-existing-key not found", err.Error())
}

func TestDeleteSecret(t *testing.T) {
	secretKey := utils.RandomString(20)
	secretValue := utils.RandomString(20)

	// add
	_, err := Add(testNamespace, secretKey, secretValue)
	assert.Nil(t, err)

	// get
	s, err := Get(testNamespace, secretKey)
	assert.Nil(t, err)

	assert.NotNil(t, s)
	assert.Equal(t, s.Key, secretKey)
	assert.Equal(t, s.Value, secretValue)

	// delete
	err = Delete(s.Namespace, s.Key)
	assert.Nil(t, err)

	// get
	_, err = Get(testNamespace, secretKey)
	assert.NotNil(t, err)
	assert.Equal(t, fmt.Sprintf("secret with key %s not found", secretKey), err.Error())
}
