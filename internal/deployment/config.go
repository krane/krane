package deployment

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"

	"github.com/docker/docker/api/types/mount"
	"github.com/docker/go-connections/nat"
	"github.com/lithammer/shortuuid/v3"

	"github.com/krane/krane/internal/constants"
	"github.com/krane/krane/internal/docker"
	"github.com/krane/krane/internal/logger"
	"github.com/krane/krane/internal/proxy"
	"github.com/krane/krane/internal/store"
)

// Config represents a deployment configuration
type Config struct {
	Name       string            `json:"name" binding:"required"`  // deployment name
	Image      string            `json:"image" binding:"required"` // container image
	Registry   string            `json:"registry"`                 // container registry
	Tag        string            `json:"tag"`                      // container image tag
	Alias      []string          `json:"alias"`                    // custom domain aliases (my-app.example.com or my-app.localhost)
	Env        map[string]string `json:"env"`                      // deployment environment variables
	Secrets    map[string]string `json:"secrets"`                  // deployment secrets resolved as environment variables
	Labels     map[string]string `json:"labels"`                   // container labels
	Ports      map[string]string `json:"ports"`                    // container ports to expose from the container to the host
	TargetPort string            `json:"target_port"`              // the target port to load-balance request through
	Volumes    map[string]string `json:"volumes"`                  // container volumes
	Command    string            `json:"command"`                  // container start command
	Entrypoint string            `json:"entrypoint"`               // container entrypoint
	Scale      int               `json:"scale"`                    // number of containers to create for the deployment
	Secure     bool              `json:"secure"`                   // enable/disable secure communication over HTTPS/TLS w/ auto generated certs
	Internal   bool              `json:"internal"`                 // whether a deployment is internal (ie. krane-proxy)
	RateLimit  uint              `json:"rate_limit"`               // requests per second for a given deployment (default 0, which means no rate limit)
}

// SaveConfig a deployment configuration into the db
func SaveConfig(config Config) error {
	config.applyDefaults()

	if err := config.isValid(); err != nil {
		logger.Errorf("deployment config is not valid %v", err)
		return err
	}

	bytes, _ := config.Serialize()
	return store.Client().Put(constants.DeploymentsCollectionName, config.Name, bytes)
}

// Serialize returns the bytes for a deployment config
func (config Config) Serialize() ([]byte, error) {
	return json.Marshal(config)
}

// DeSerialize returns a config from bytes
func DeSerializeConfig(bytes []byte) (Config, error) {
	var config Config
	err := json.Unmarshal(bytes, &config)
	return config, err
}

// applyDefaults applies default deployment configuration values
func (config *Config) applyDefaults() {
	if config.Registry == "" {
		config.Registry = "docker.io"
	}

	if config.Alias == nil {
		config.Alias = make([]string, 0)
	}

	if config.Labels == nil {
		config.Labels = make(map[string]string, 0)
	}

	if config.Secrets == nil {
		config.Secrets = make(map[string]string, 0)
	}

	if config.Env == nil {
		config.Env = make(map[string]string, 0)
	}

	if config.Volumes == nil {
		config.Volumes = make(map[string]string, 0)
	}

	if config.Ports == nil {
		config.Ports = make(map[string]string, 0)
	}

	if config.Tag == "" {
		config.Tag = "latest"
	}

	return
}

// isValid returns an error if a deployment Configuration is not valid
func (config Config) isValid() error {
	isValidName := config.isValidName()
	if !isValidName {
		return fmt.Errorf("invalid name %s in deployment config", config.Name)
	}

	if config.Image == "" {
		return errors.New("image required in deployment config")
	}

	return nil
}

// isValidName return if a deployment name is valid or not
func (config Config) isValidName() bool {
	if len(config.Name) > 50 {
		return false
	}

	startsWithLetter := "[a-z]"
	allowedCharacters := "[a-z0-9_-]"
	endWithLowerCaseAlphanumeric := "[0-9a-z]"
	characterLimit := "{1,}"

	matchers := fmt.Sprintf(`^(%s%s*%s)%s$`, // ^[a-z][a-z0-9_-]*[0-9a-z]{1,48}$
		startsWithLetter,
		allowedCharacters,
		endWithLowerCaseAlphanumeric,
		characterLimit)

	match := regexp.MustCompile(matchers)
	return match.MatchString(config.Name)
}

// GetDeploymentConfig returns a deployments configuration
func GetDeploymentConfig(deployment string) (Config, error) {
	bytes, err := store.Client().Get(constants.DeploymentsCollectionName, deployment)
	if err != nil {
		return Config{}, err
	}

	if bytes == nil {
		return Config{}, fmt.Errorf("deployment %s not found", deployment)
	}

	config, err := DeSerializeConfig(bytes)
	if err != nil {
		return Config{}, err
	}

	return config, nil
}

// GetAllDeploymentConfigs returns a list of all deployment configurations
func GetAllDeploymentConfigs() ([]Config, error) {
	bytes, err := store.Client().GetAll(constants.DeploymentsCollectionName)
	if err != nil {
		return make([]Config, 0), err
	}

	deployments := make([]Config, 0)
	for _, b := range bytes {
		config, _ := DeSerializeConfig(b)
		deployments = append(deployments, config)
	}

	return deployments, nil
}

// DeleteSecret removes a deployment configuration from the db
func DeleteConfig(deployment string) error {
	return store.Client().Remove(constants.DeploymentsCollectionName, deployment)
}

// Empty returns true if a config has not defined a deployment name or image
func (config Config) Empty() bool {
	return config.Name == "" || config.Image == ""
}

// DockerConfig returns the docker configuration for creating a container
func (config Config) DockerConfig() docker.DockerConfig {
	kraneNetwork, err := docker.GetClient().GetNetworkByName(docker.KraneNetworkName)
	if err != nil {
		return docker.DockerConfig{}
	}

	var command []string
	var entrypoint []string

	if config.Command != "" {
		command = append(command, config.Command)
	}

	if config.Entrypoint != "" {
		entrypoint = append(entrypoint, config.Entrypoint)
	}

	containerName := fmt.Sprintf("%s-%s", config.Name, shortuuid.New())
	return docker.DockerConfig{
		ContainerName: containerName,
		Image:         config.Image,
		NetworkID:     kraneNetwork.ID,
		Labels:        config.DockerLabels(),
		Ports:         config.DockerPorts(),
		PortSet:       config.DockerPortSet(),
		VolumeMounts:  config.DockerVolumeMount(),
		VolumeSet:     config.DockerVolumeSet(),
		Env:           config.DockerEnvs(),
		Command:       command,
		Entrypoint:    entrypoint,
	}
}

// DockerEnvs returns a list of formatted Docker environment variables
func (config Config) DockerEnvs() []string {
	envs := make([]string, 0)

	// environment variables sourced from the deployment config
	for k, v := range config.Env {
		envs = append(envs, fmt.Sprintf("%s=%s", k, v))
	}

	// secrets specified in the deployment config which work the same as environment variables
	// but with resolved values located server side
	for key, alias := range config.Secrets {
		secret, err := GetSecret(config.Name, key)
		if err != nil || secret == nil {
			logger.Infof("unable to resolve secret for %s with alias %s", config.Name, alias)
			continue
		}
		envs = append(envs, fmt.Sprintf("%s=%s", key, secret.Value))
	}

	return envs
}

// DockerLabels returns a map of Docker labels that are applied to Krane managed containers
func (config Config) DockerLabels() map[string]string {
	config.Labels[docker.ContainerDeploymentLabel] = config.Name
	config.ApplyProxyLabels()
	return config.Labels
}

// ApplyProxyLabels applies network labels to a deployment config
func (config Config) ApplyProxyLabels() {
	// default traefik labels
	config.Labels["traefik.enable"] = "true"
	config.Labels["traefik.docker.network"] = docker.KraneNetworkName

	// router labels
	for k, v := range proxy.TraefikRouterLabels(config.Name, config.Alias, config.Secure) {
		config.Labels[k] = v
	}

	// middleware labels
	for k, v := range proxy.TraefikMiddlewareLabels(config.Name, config.Secure, config.RateLimit) {
		config.Labels[k] = v
	}

	// service labels
	for k, v := range proxy.TraefikServiceLabels(config.Name, config.Ports, config.TargetPort) {
		config.Labels[k] = v
	}
}

// DockerVolumeMount returns a list of formatted Docker volume mounts
func (config Config) DockerVolumeMount() []mount.Mount {
	volumes := make([]mount.Mount, 0)
	for hostVolume, containerVolume := range config.Volumes {
		volumes = append(volumes, mount.Mount{
			Type:   mount.TypeBind,
			Source: hostVolume,
			Target: containerVolume,
		})
	}
	return volumes
}

// DockerVolumeSet returns a set of Docker formatted volumes
func (config Config) DockerVolumeSet() map[string]struct{} {
	volumes := make(map[string]struct{}, 0)
	for _, containerVolume := range config.Volumes {
		volumes[containerVolume] = struct{}{}
	}
	return volumes
}

// DockerPorts returns Docker formatted port map
func (config Config) DockerPorts() nat.PortMap {
	bindings := nat.PortMap{}
	for hostPort, containerPort := range config.Ports {
		if hostPort == "" {
			// randomly assign a host port if no explicit host port to bind to was provided
			freePort, err := findFreePort()
			if err != nil {
				logger.Errorf("Error looking for a free port on host machine %v", err)
				continue
			}
			hostPort = freePort
		}

		hostBinding := nat.PortBinding{HostPort: hostPort}

		// TODO: expose functionality for binding to other protocols besides tcp
		cPort, err := nat.NewPort(string(TCP), containerPort)
		if err != nil {
			logger.Errorf("Error creating a new container port %v", err)
			continue
		}

		bindings[cPort] = []nat.PortBinding{hostBinding}
	}

	return bindings
}

// DockerPortSet returns Docker formatted port set
func (config Config) DockerPortSet() nat.PortSet {
	bindings := nat.PortSet{}
	for _, containerPort := range config.Ports {
		cPort, err := nat.NewPort(string(TCP), containerPort)
		if err != nil {
			logger.Errorf("Error creating a new container port %v", err)
			continue
		}
		bindings[cPort] = struct{}{}
	}
	return bindings
}
