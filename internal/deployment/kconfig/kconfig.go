package kconfig

import (
	"errors"

	"github.com/biensupernice/krane/internal/constants"
	"github.com/biensupernice/krane/internal/store"
)

type Kconfig struct {
	Name     string            `json:"name" binding:"required"`
	Registry string            `json:"registry"`
	Image    string            `json:"image" binding:"required"`
	Tag      string            `json:"tag"`   // docker image tag
	Alias    []string          `json:"alias"` // custom domain aliases
	Env      map[string]string `json:"env"`
	Ports    map[string]string `json:"ports"`
	Secrets  map[string]string `json:"secrets"`
	Volumes  map[string]string `json:"volumes"`
	Command  string            `json:"command"`
	Scale    int               `json:"scale"`   // number of containers for a deployment
	Secured  bool              `json:"secured"` // enable/disable secure communication over HTTPS/TLS
}

// Apply :
func (cfg *Kconfig) Apply() error {
	if err := cfg.isValid(); err != nil {
		return err
	}

	cfg.applyDefaults()

	bytes, _ := cfg.Serialize()
	return store.Instance().Put(constants.DeploymentsCollectionName, cfg.Name, bytes)
}

// Delete :
func Delete(deploymentName string) error {
	return store.Instance().Remove(constants.DeploymentsCollectionName, deploymentName)
}

// GetConfigByDeploymentByName :
func GetConfigByDeploymentByName(deploymentName string) (Kconfig, error) {
	bytes, err := store.Instance().Get(constants.DeploymentsCollectionName, deploymentName)
	if err != nil {
		return Kconfig{}, err
	}

	if bytes == nil {
		return Kconfig{}, errors.New("Deployment not found")
	}

	var cfg Kconfig
	err = store.Deserialize(bytes, &cfg)
	if err != nil {
		return Kconfig{}, err
	}

	return cfg, nil
}

// GetAllDeploymentConfigs :
func GetAllDeploymentConfigs() ([]Kconfig, error) {
	bytes, err := store.Instance().GetAll(constants.DeploymentsCollectionName)
	if err != nil {
		return make([]Kconfig, 0), err
	}

	cfgs := make([]Kconfig, 0)
	for _, b := range bytes {
		var cfg Kconfig
		_ = store.Deserialize(b, &cfg)
		cfgs = append(cfgs, cfg)
	}

	return cfgs, nil
}
