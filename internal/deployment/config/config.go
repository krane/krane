package config

import (
	"errors"

	"github.com/biensupernice/krane/internal/constants"
	"github.com/biensupernice/krane/internal/store"
)

type Config struct {
	Name     string            `json:"name" binding:"required"`
	Registry string            `json:"registry"`
	Image    string            `json:"image" binding:"required"`
	Tag      string            `json:"tag"`
	Alias    []string          `json:"alias"`
	Env      map[string]string `json:"env"`
	Ports    map[string]string `json:"ports"`
	Secrets  map[string]string `json:"secrets"`
	Volumes  map[string]string `json:"volumes"`
}

func (cfg *Config) Save() error {
	if err := cfg.validate(); err != nil {
		return err
	}

	cfg.applyDefaults()

	bytes, _ := cfg.Serialize()
	return store.Instance().Put(constants.DeploymentsCollectionName, cfg.Name, bytes)
}

func GetConfigByDeploymentByName(deploymentName string) (Config, error) {
	bytes, err := store.Instance().Get(constants.DeploymentsCollectionName, deploymentName)
	if err != nil {
		return Config{}, err
	}

	if bytes == nil {
		return Config{}, errors.New("Deployment not found")
	}

	var cfg Config
	err = store.Deserialize(bytes, &cfg)
	if err != nil {
		return Config{}, err
	}

	return cfg, nil
}

func GetAllDeploymentConfigs() ([]Config, error) {
	bytes, err := store.Instance().GetAll(constants.DeploymentsCollectionName)
	if err != nil {
		return make([]Config, 0), err
	}

	cfgs := make([]Config, 0)
	for _, b := range bytes {
		var cfg Config
		_ = store.Deserialize(b, &cfg)
		cfgs = append(cfgs, cfg)
	}

	return cfgs, nil
}
