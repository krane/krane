package kconfig

import (
	"errors"

	"github.com/biensupernice/krane/internal/constants"
	"github.com/biensupernice/krane/internal/store"
)

type Kconfig struct {
	Name       string            `json:"name" binding:"required"`
	Registry   string            `json:"registry"`
	Image      string            `json:"image" binding:"required"`
	Tag        string            `json:"tag"`   // docker image tag
	Alias      []string          `json:"alias"` // custom domain aliases
	Labels     map[string]string `json:"labels"`
	Env        map[string]string `json:"env"`
	Secrets    map[string]string `json:"secrets"`
	Ports      map[string]string `json:"ports"`
	Volumes    map[string]string `json:"volumes"`
	Command    string            `json:"command"`
	Entrypoint string            `json:"entrypoint"`
	Scale      int               `json:"scale"`   // number of containers for a deployment
	Secured    bool              `json:"secured"` // enable/disable secure communication over HTTPS/TLS
}

// Save : creates or updates a deployment
func (cfg *Kconfig) Save() error {
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

// GetConfigByDeploymentName :
func GetConfigByDeploymentName(deploymentName string) (Kconfig, error) {
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
