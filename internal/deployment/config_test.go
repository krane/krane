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
