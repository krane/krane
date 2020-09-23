package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	validKConfig1 = Config{
		Name:  "example",
		Image: "biensupernice/krane",
	}
	validKConfig2 = Config{
		Name:     "example_hello-world_deployment",
		Registry: "github.io",
		Image:    "biensupernice/krane",
	}
	invalidKConfig1 = Config{
		Name:  "example-$123",
		Image: "biensupernice/krane",
	}
)

func TestKraneConfigValidation(t *testing.T) {
	assert.Nil(t, validKConfig1.validate())
	assert.Nil(t, validKConfig2.validate())
	assert.NotNil(t, invalidKConfig1.validate())
}

func TestKraneConfigNameValidation(t *testing.T) {
	assert.True(t, Config{Name: "example"}.validateName())
	assert.True(t, Config{Name: "example-_hello-world_deployment"}.validateName())
	assert.True(t, Config{Name: "example-_ac913f19-aa6b-4887-92e9-6102a9fa171b"}.validateName())
	assert.True(t, Config{Name: "example_-f9186962-c575-11ea-87d0-0242ac130003"}.validateName())
	assert.True(t, Config{Name: "small"}.validateName())
	assert.True(t, Config{Name: "tiny"}.validateName())
	assert.True(t, Config{Name: "api"}.validateName())
	assert.True(t, Config{Name: "up"}.validateName())

	assert.False(t, Config{Name: "x"}.validateName())
	assert.False(t, Config{Name: "-example"}.validateName())
	assert.False(t, Config{Name: "_-example"}.validateName())
	assert.False(t, Config{Name: "example_-"}.validateName())
	assert.False(t, Config{Name: "example-_"}.validateName())
	assert.False(t, Config{Name: "example-_$$2901"}.validateName())
	assert.False(t, Config{Name: "example#-_&2901"}.validateName())
	assert.False(t, Config{Name: "example^-_!2901"}.validateName())
}
