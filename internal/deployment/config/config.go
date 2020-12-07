package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"

	"github.com/biensupernice/krane/internal/constants"
	"github.com/biensupernice/krane/internal/store"
)

type DeploymentConfig struct {
	Name       string            `json:"name" binding:"required"`  // deployment name
	Registry   string            `json:"registry"`                 // container registry
	Image      string            `json:"image" binding:"required"` // container image
	Tag        string            `json:"tag"`                      // container image tag
	Alias      []string          `json:"alias"`                    // custom domain aliases (my-app.example.com)
	Env        map[string]string `json:"env"`                      // deployment environment variables
	Secrets    map[string]string `json:"secrets"`                  // deployment secrets
	Labels     map[string]string `json:"labels"`                   // container labels
	Ports      map[string]string `json:"ports"`                    // container ports
	Volumes    map[string]string `json:"volumes"`                  // container volumes
	Command    string            `json:"command"`                  // container start command
	Entrypoint string            `json:"entrypoint"`               // container entrypoint
	Scale      int               `json:"scale"`                    // number of containers
	Secured    bool              `json:"secured"`                  // enable/disable secure communication over HTTPS/TLS
}

// Save : save a deployment
func (cfg *DeploymentConfig) Save() error {
	if err := cfg.isValid(); err != nil {
		return err
	}

	cfg.applyDefaults()
	bytes, _ := cfg.Serialize()
	return store.Client().Put(constants.DeploymentsCollectionName, cfg.Name, bytes)
}

// Delete : delete a deployment
func Delete(deploymentName string) error {
	return store.Client().Remove(constants.DeploymentsCollectionName, deploymentName)
}

// GetDeploymentConfig : get a deployment configuration
func GetDeploymentConfig(deploymentName string) (DeploymentConfig, error) {
	bytes, err := store.Client().Get(constants.DeploymentsCollectionName, deploymentName)
	if err != nil {
		return DeploymentConfig{}, err
	}

	if bytes == nil {
		return DeploymentConfig{}, fmt.Errorf("Deployment %s not found", deploymentName)
	}

	var cfg DeploymentConfig
	err = store.Deserialize(bytes, &cfg)
	if err != nil {
		return DeploymentConfig{}, err
	}

	return cfg, nil
}

// GetAllDeploymentConfigurations : get all deployment configurations
func GetAllDeploymentConfigurations() ([]DeploymentConfig, error) {
	bytes, err := store.Client().GetAll(constants.DeploymentsCollectionName)
	if err != nil {
		return make([]DeploymentConfig, 0), err
	}

	deployments := make([]DeploymentConfig, 0)
	for _, b := range bytes {
		var cfg DeploymentConfig
		_ = store.Deserialize(b, &cfg)
		deployments = append(deployments, cfg)
	}

	return deployments, nil
}

// Serialize : serialize a deployment config into bytes
func (cfg DeploymentConfig) Serialize() ([]byte, error) { return json.Marshal(cfg) }

// isValid : validate deployment configuration
func (cfg DeploymentConfig) isValid() error {
	isValidName := cfg.validateDeploymentName()
	if !isValidName {
		return errors.New("invalid name in deployment config")
	}

	if cfg.Image == "" {
		return errors.New("image required in deployment config")
	}

	return nil
}

// applyDefaults : apply default deployment values
func (cfg *DeploymentConfig) applyDefaults() {
	if cfg.Registry == "" {
		cfg.Registry = "docker.io"
	}

	if cfg.Alias == nil {
		cfg.Alias = make([]string, 0)
	}

	if cfg.Labels == nil {
		cfg.Labels = make(map[string]string, 0)
	}

	if cfg.Secrets == nil {
		cfg.Secrets = make(map[string]string, 0)
	}

	if cfg.Env == nil {
		cfg.Env = make(map[string]string, 0)
	}

	if cfg.Volumes == nil {
		cfg.Volumes = make(map[string]string, 0)
	}

	if cfg.Ports == nil {
		cfg.Ports = make(map[string]string, 0)
	}

	if cfg.Tag == "" {
		cfg.Tag = "latest"
	}

	return
}

// validateDeploymentName : verify no funky business in the deployment name
func (cfg DeploymentConfig) validateDeploymentName() bool {
	startsWithLetter := "[a-z]"
	allowedCharacters := "[a-z0-9_-]"
	endWithLowerCaseAlphanumeric := "[0-9a-z]"
	characterLimit := "{1,}"

	matchers := fmt.Sprintf(`^%s%s*%s%s$`, // ^[a - z][a - z0 - 9_ -]*[0-9a-z]$
		startsWithLetter,
		allowedCharacters,
		endWithLowerCaseAlphanumeric,
		characterLimit)

	match := regexp.MustCompile(matchers)
	return match.MatchString(cfg.Name)
}
