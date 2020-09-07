package kranecfg

import (
	"errors"

	"github.com/biensupernice/krane/internal/collection"
	"github.com/biensupernice/krane/internal/store"
)

type KraneConfig struct {
	Name     string            `json:"name" binding:"required"`  // Deployment name
	Registry string            `json:"registry"`                 // Docker registry (default to `docker.io`)
	Image    string            `json:"image" binding:"required"` // Docker Image
	Tag      string            `json:"tag"`                      // Docker tag (defaults to `latest`)
	Alias    []string          `json:"alias"`                    //  alias or aliases (ex: [`my-app.com`, `api.my-app.com`]
	Env      map[string]string `json:"env"`
	Secrets  map[string]string `json:"secrets"`
	Volumes  map[string]string `json:"volumes"`
}

// Save : a Krane Config
func (cfg *KraneConfig) Save() error {
	err := cfg.validate()
	if err != nil {
		return err
	}

	cfg.applyDefaults()

	bytes, err := store.Serialize(cfg)
	err = store.Instance().Put(collection.Deployments, cfg.Name, bytes)
	if err != nil {
		return err
	}

	return nil
}

// Delete : a Krane Config
func Delete(name string) error {
	err := store.Instance().Remove(collection.Deployments, name)
	if err != nil {
		return err
	}

	return nil
}

// Get : returns a KraneConfig
func Get(deploymentName string) (KraneConfig, error) {
	bytes, err := store.Instance().Get(collection.Deployments, deploymentName)
	if err != nil {
		return KraneConfig{}, err
	}

	if bytes == nil {
		return KraneConfig{}, errors.New("Deployment not found")
	}

	var cfg KraneConfig
	err = store.Deserialize(bytes, &cfg)
	if err != nil {
		return KraneConfig{}, err
	}

	return cfg, nil
}

// GetAll: returns a list of KraneConfigs
func GetAll() ([]KraneConfig, error) {
	bytes, err := store.Instance().GetAll(collection.Deployments)
	if err != nil {
		return make([]KraneConfig, 0), err
	}

	cfgs := make([]KraneConfig, 0)
	for _, b := range bytes {
		var cfg KraneConfig
		store.Deserialize(b, &cfg)
		cfgs = append(cfgs, cfg)
	}

	return cfgs, nil
}
