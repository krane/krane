package spec

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	validSpec1 = Spec{
		Name: "example",
		Config: Config{
			Image: "biensupernice/krane",
		},
	}
	validSpec2 = Spec{
		Name: "example_hello-world_deployment",
		Config: Config{
			Registry:      "github.io",
			Image:         "biensupernice/krane",
			ContainerPort: "8080",
			HostPort:      "80",
		},
	}
)

var (
	invalidSpec1 = Spec{
		Name: "example-$123",
		Config: Config{
			Image: "biensupernice/krane",
		},
	}
)

func TestSpecValidation(t *testing.T) {
	assert.Nil(t, validSpec1.Validate())
	assert.Nil(t, validSpec2.Validate())

	assert.NotNil(t, invalidSpec1.Validate())
}

func TestSpecNameValidation(t *testing.T) {
	assert.True(t, Spec{Name: "example"}.isValidSpecName())
	assert.True(t, Spec{Name: "example-_hello-world_deployment"}.isValidSpecName())
	assert.True(t, Spec{Name: "example-_ac913f19-aa6b-4887-92e9-6102a9fa171b"}.isValidSpecName())
	assert.True(t, Spec{Name: "example_-f9186962-c575-11ea-87d0-0242ac130003"}.isValidSpecName())
	assert.True(t, Spec{Name: "small"}.isValidSpecName())
	assert.True(t, Spec{Name: "tiny"}.isValidSpecName())

	assert.False(t, Spec{Name: "-example"}.isValidSpecName())
	assert.False(t, Spec{Name: "_-example"}.isValidSpecName())
	assert.False(t, Spec{Name: "example_-"}.isValidSpecName())
	assert.False(t, Spec{Name: "example-_"}.isValidSpecName())
	assert.False(t, Spec{Name: "example-_$$2901"}.isValidSpecName())
	assert.False(t, Spec{Name: "example#-_&2901"}.isValidSpecName())
	assert.False(t, Spec{Name: "example^-_!2901"}.isValidSpecName())
}
