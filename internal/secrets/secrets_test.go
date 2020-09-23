package secrets

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/biensupernice/krane/internal/constants"
	"github.com/biensupernice/krane/internal/deployment/config"
	"github.com/biensupernice/krane/internal/store"
	"github.com/biensupernice/krane/internal/utils"
)

const boltpath = "./krane.db"
const testNamespace = "krane-test"

func teardown() { os.Remove(boltpath) }

func TestMain(m *testing.M) {
	store.New((boltpath))
	defer store.Instance().Shutdown()

	// Create deployment (namespace)
	deployment := config.Config{Name: testNamespace}
	bytes, _ := deployment.Serialize()
	store.Instance().Put(constants.DeploymentsCollectionName, deployment.Name, bytes)

	code := m.Run()

	teardown()
	os.Exit(code)
}

func TestAddNewSecret(t *testing.T) {
	alias1, err := Add("token", "biensupernice", testNamespace)
	assert.Nil(t, err)
	assert.Equal(t, "@TOKEN", alias1)

	alias2, err := Add("api_token", "biensupernice", testNamespace)
	assert.Nil(t, err)
	assert.Equal(t, "@API_TOKEN", alias2)

	alias3, err := Add("api-token", "biensupernice", testNamespace)
	assert.Nil(t, err)
	assert.Equal(t, "@API_TOKEN", alias3)

	alias4, err := Add("api-token123", "biensupernice", testNamespace)
	assert.Nil(t, err)
	assert.Equal(t, "@API_TOKEN123", alias4)

	alias5, err := Add("API_PORT_8080", "8080", testNamespace)
	assert.Nil(t, err)
	assert.Equal(t, "@API_PORT_8080", alias5)

	alias6, err := Add("API-PORT-8080", "8080", testNamespace)
	assert.Nil(t, err)
	assert.Equal(t, "@API_PORT_8080", alias6)

	alias7, err := Add("env", "dev", testNamespace)
	assert.Nil(t, err)
	assert.Equal(t, "@ENV", alias7)

	alias8, err := Add("8080_API_PORT", "8080", testNamespace)
	assert.Nil(t, err)
	assert.Equal(t, "@8080_API_PORT", alias8)

	alias9, err := Add("8080-API-PORT", "8080", testNamespace)
	assert.Nil(t, err)
	assert.Equal(t, "@8080_API_PORT", alias9)

	alias10, err := Add("aPi_ToKeN-1337", "8080", testNamespace)
	assert.Nil(t, err)
	assert.Equal(t, "@API_TOKEN_1337", alias10)

}

func TestErrorWhenAddSecretToNonExistingNamespace(t *testing.T) {
	_, err1 := Add("TOKEN", "biensupernice", "non-existing-namespace")
	assert.Error(t, err1)
	assert.Equal(t, "unable to find namespace non-existing-namespace", err1.Error())

	_, err2 := Add("TOKEN", "biensupernice", "")
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
	assert.False(t, isValidSecretKey(utils.RandomString(21)))

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

	alias, err := Add(secretKey, secretValue, testNamespace)
	assert.Nil(t, err)

	secrets, err := GetAll(testNamespace)
	assert.Nil(t, err)
	assert.True(t, len(secrets) > 0)

	var s Secret
	for _, secret := range secrets {
		if secret.Alias == alias {
			s = *secret
			break
		}
	}

	assert.NotNil(t, s)
	assert.Equal(t, s.Key, secretKey)
	assert.Equal(t, s.Value, secretValue)
}

func TestGetSecretByAlias(t *testing.T) {
	secretKey := utils.RandomString(20)
	secretValue := utils.RandomString(20)

	alias, err := Add(secretKey, secretValue, testNamespace)
	assert.Nil(t, err)

	s, err := Get(testNamespace, alias)
	assert.Nil(t, err)

	assert.NotNil(t, s)
	assert.Equal(t, s.Key, secretKey)
	assert.Equal(t, s.Value, secretValue)
}

func TestErrorWhenGetSecretByNonExistingAlias(t *testing.T) {
	_, err := Get(testNamespace, "@NON-EXISTING-ALIAS")
	assert.NotNil(t, err)
	assert.Equal(t, "secret with alias @NON-EXISTING-ALIAS not found", err.Error())
}
