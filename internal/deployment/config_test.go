package deployment

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMinimalDeploymentConfig(t *testing.T) {
	assert.Nil(t, Config{Name: "example-deployment", Image: "biensupernice/krane"}.isValid())
}

func TestInvalidDeployment(t *testing.T) {
	assert.Error(t, Config{Name: "missing-image"}.isValid())
	assert.Error(t, Config{Name: "e", Image: "biensupernice/krane"}.isValid())
	assert.Error(t, Config{Name: "$example-123", Image: "biensupernice/krane"}.isValid())
	assert.Error(t, Config{Name: "example-$123", Image: "biensupernice/krane"}.isValid())
}

func TestValidDeploymentNames(t *testing.T) {
	assert.True(t, Config{Name: "example"}.isValidName())
	assert.True(t, Config{Name: "example-_hello-world_deployment"}.isValidName())
	assert.True(t, Config{Name: "example-_ac913f19-aa6b-4887-92e9-6102a9fa171b"}.isValidName())
	assert.True(t, Config{Name: "example_-f9186962-c575-11ea-87d0-0242ac130003"}.isValidName())
	assert.True(t, Config{Name: "small"}.isValidName())
	assert.True(t, Config{Name: "tiny"}.isValidName())
	assert.True(t, Config{Name: "api"}.isValidName())
	assert.True(t, Config{Name: "up"}.isValidName())
}

func TestInvalidDeploymentNames(t *testing.T) {
	assert.False(t, Config{Name: "x"}.isValidName())
	assert.False(t, Config{Name: "-example"}.isValidName())
	assert.False(t, Config{Name: "_-example"}.isValidName())
	assert.False(t, Config{Name: "example_-"}.isValidName())
	assert.False(t, Config{Name: "example-_"}.isValidName())
	assert.False(t, Config{Name: "example-_$$2901"}.isValidName())
	assert.False(t, Config{Name: "example#-_&2901"}.isValidName())
	assert.False(t, Config{Name: "example^-_!2901"}.isValidName())
	assert.False(t, Config{Name: "7-example"}.isValidName())
	assert.False(t, Config{Name: "7example"}.isValidName())
	assert.False(t, Config{Name: strings.Repeat("a", 51)}.isValidName()) // max deployment name chars is 50 chars
}

func TestResolveRegistryCredentials_WithSecrets(t *testing.T) {
	config := Config{
		Name:  "test-name",
		Image: "test-image",
		Registry: Registry{
			URL:      "@TEST_URL",
			Username: "@TEST_USERNAME",
			Password: "@TEST_PASSWORD",
		},
	}

	// setup
	urlSecret, err := AddSecret(config.Name, "TEST_URL", "test-url")
	assert.Nil(t, err)
	usernameSecret, err := AddSecret(config.Name, "TEST_USERNAME", "test")
	assert.Nil(t, err)
	passwordSecret, err := AddSecret(config.Name, "TEST_PASSWORD", "123")
	assert.Nil(t, err)

	// act
	err = config.ResolveRegistryCredentials()
	assert.Nil(t, err)

	// assert
	assert.Equal(t, urlSecret.Value, config.Registry.URL)
	assert.Equal(t, usernameSecret.Value, config.Registry.Username)
	assert.Equal(t, passwordSecret.Value, config.Registry.Password)
}

func TestResolveRegistryCredentials_NoSecrets(t *testing.T) {
	config := Config{
		Name:  "test-name",
		Image: "test-image",
		Registry: Registry{
			URL:      "test.io",
			Username: "username",
			Password: "password",
		},
	}

	// act
	err := config.ResolveRegistryCredentials()
	assert.Nil(t, err)

	// assert
	assert.Equal(t, "test.io", config.Registry.URL)
	assert.Equal(t, "username", config.Registry.Username)
	assert.Equal(t, "password", config.Registry.Password)
}

func TestResolveRegistryCredentials_NoRegistry(t *testing.T) {
	config := Config{
		Name:  "test-name",
		Image: "test-image",
	}

	// act
	err := config.ResolveRegistryCredentials()
	assert.Nil(t, err)

	// assert
	assert.Equal(t, "", config.Registry.URL)
	assert.Equal(t, "", config.Registry.Username)
	assert.Equal(t, "", config.Registry.Password)
}

func TestResolveRegistryCredentials_SecretNotFound(t *testing.T) {
	config := Config{
		Name:  "test-name",
		Image: "test-image",
		Registry: Registry{
			URL:      "@TEST_URL_NOTFOUND",
			Username: "@TEST_USERNAME_NOTFOUND",
			Password: "@TEST_PASSWORD_NOTFOUND",
		},
	}
	err := config.ResolveRegistryCredentials()
	assert.Error(t, err, "secret \"@TEST_URL\" not found")
}
