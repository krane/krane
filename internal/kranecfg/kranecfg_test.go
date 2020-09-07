package kranecfg

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	validKConfig1 = KraneConfig{
		Name:  "example",
		Image: "biensupernice/krane",
	}
	validKConfig2 = KraneConfig{
		Name:     "example_hello-world_deployment",
		Registry: "github.io",
		Image:    "biensupernice/krane",
	}
	invalidKConfig1 = KraneConfig{
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
	assert.True(t, KraneConfig{Name: "example"}.validateName())
	assert.True(t, KraneConfig{Name: "example-_hello-world_deployment"}.validateName())
	assert.True(t, KraneConfig{Name: "example-_ac913f19-aa6b-4887-92e9-6102a9fa171b"}.validateName())
	assert.True(t, KraneConfig{Name: "example_-f9186962-c575-11ea-87d0-0242ac130003"}.validateName())
	assert.True(t, KraneConfig{Name: "small"}.validateName())
	assert.True(t, KraneConfig{Name: "tiny"}.validateName())
	assert.True(t, KraneConfig{Name: "api"}.validateName())
	assert.True(t, KraneConfig{Name: "up"}.validateName())

	assert.False(t, KraneConfig{Name: "x"}.validateName())
	assert.False(t, KraneConfig{Name: "-example"}.validateName())
	assert.False(t, KraneConfig{Name: "_-example"}.validateName())
	assert.False(t, KraneConfig{Name: "example_-"}.validateName())
	assert.False(t, KraneConfig{Name: "example-_"}.validateName())
	assert.False(t, KraneConfig{Name: "example-_$$2901"}.validateName())
	assert.False(t, KraneConfig{Name: "example#-_&2901"}.validateName())
	assert.False(t, KraneConfig{Name: "example^-_!2901"}.validateName())
}
