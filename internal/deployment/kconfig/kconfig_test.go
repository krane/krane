package kconfig

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	validKConfig1 = Kconfig{
		Name:  "example",
		Image: "biensupernice/krane",
	}
	validKConfig2 = Kconfig{
		Name:     "example_hello-world_deployment",
		Registry: "github.io",
		Image:    "biensupernice/krane",
	}
	invalidKConfig1 = Kconfig{
		Name:  "example-$123",
		Image: "biensupernice/krane",
	}
)

func TestKraneConfigValidation(t *testing.T) {
	assert.Nil(t, validKConfig1.isValid())
	assert.Nil(t, validKConfig2.isValid())
	assert.NotNil(t, invalidKConfig1.isValid())
}

func TestKraneConfigNameValidation(t *testing.T) {
	assert.True(t, Kconfig{Name: "example"}.validateName())
	assert.True(t, Kconfig{Name: "example-_hello-world_deployment"}.validateName())
	assert.True(t, Kconfig{Name: "example-_ac913f19-aa6b-4887-92e9-6102a9fa171b"}.validateName())
	assert.True(t, Kconfig{Name: "example_-f9186962-c575-11ea-87d0-0242ac130003"}.validateName())
	assert.True(t, Kconfig{Name: "small"}.validateName())
	assert.True(t, Kconfig{Name: "tiny"}.validateName())
	assert.True(t, Kconfig{Name: "api"}.validateName())
	assert.True(t, Kconfig{Name: "up"}.validateName())

	assert.False(t, Kconfig{Name: "x"}.validateName())
	assert.False(t, Kconfig{Name: "-example"}.validateName())
	assert.False(t, Kconfig{Name: "_-example"}.validateName())
	assert.False(t, Kconfig{Name: "example_-"}.validateName())
	assert.False(t, Kconfig{Name: "example-_"}.validateName())
	assert.False(t, Kconfig{Name: "example-_$$2901"}.validateName())
	assert.False(t, Kconfig{Name: "example#-_&2901"}.validateName())
	assert.False(t, Kconfig{Name: "example^-_!2901"}.validateName())
}
