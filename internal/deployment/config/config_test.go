package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidDeploymentConfigs(t *testing.T) {
	assert.Nil(t, DeploymentConfig{
		Name:  "example-deployment",
		Image: "biensupernice/krane",
	}.isValid())
	assert.Nil(t, DeploymentConfig{
		Name:     "example-deployment",
		Registry: "github.io",
		Image:    "biensupernice/krane",
	}.isValid())
}

func TestInvalidNameInDeploymentConfig(t *testing.T) {
	assert.NotNil(t, DeploymentConfig{
		Name:  "example-$123",
		Image: "biensupernice/krane",
	}.isValid())
}

func TestMissingImageInDeploymentConfig(t *testing.T) {
	assert.NotNil(t, DeploymentConfig{
		Name: "example-deployment",
	}.isValid())
}

func TestDeploymentConfigNameValidation(t *testing.T) {
	assert.True(t, DeploymentConfig{Name: "example"}.validateDeploymentName())
	assert.True(t, DeploymentConfig{Name: "example-_hello-world_deployment"}.validateDeploymentName())
	assert.True(t, DeploymentConfig{Name: "example-_ac913f19-aa6b-4887-92e9-6102a9fa171b"}.validateDeploymentName())
	assert.True(t, DeploymentConfig{Name: "example_-f9186962-c575-11ea-87d0-0242ac130003"}.validateDeploymentName())
	assert.True(t, DeploymentConfig{Name: "small"}.validateDeploymentName())
	assert.True(t, DeploymentConfig{Name: "tiny"}.validateDeploymentName())
	assert.True(t, DeploymentConfig{Name: "api"}.validateDeploymentName())
	assert.True(t, DeploymentConfig{Name: "up"}.validateDeploymentName())

	assert.False(t, DeploymentConfig{Name: "x"}.validateDeploymentName())
	assert.False(t, DeploymentConfig{Name: "-example"}.validateDeploymentName())
	assert.False(t, DeploymentConfig{Name: "_-example"}.validateDeploymentName())
	assert.False(t, DeploymentConfig{Name: "example_-"}.validateDeploymentName())
	assert.False(t, DeploymentConfig{Name: "example-_"}.validateDeploymentName())
	assert.False(t, DeploymentConfig{Name: "example-_$$2901"}.validateDeploymentName())
	assert.False(t, DeploymentConfig{Name: "example#-_&2901"}.validateDeploymentName())
	assert.False(t, DeploymentConfig{Name: "example^-_!2901"}.validateDeploymentName())
}
