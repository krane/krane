package config

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"regexp"
	"strconv"

	"github.com/docker/docker/api/types/mount"
	"github.com/docker/go-connections/nat"
	"github.com/lithammer/shortuuid/v3"
	"github.com/pkg/errors"

	"github.com/biensupernice/krane/internal/constants"
	"github.com/biensupernice/krane/internal/deployment/secrets"
	"github.com/biensupernice/krane/internal/docker"
	"github.com/biensupernice/krane/internal/logger"
	"github.com/biensupernice/krane/internal/proxy"
	"github.com/biensupernice/krane/internal/store"
)

type DeploymentConfig struct {
	Name       string            `json:"name" binding:"required"`  // deployment name
	Image      string            `json:"image" binding:"required"` // container image
	Registry   string            `json:"registry"`                 // container registry
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

// isValid : validate deployment configuration
func (cfg DeploymentConfig) isValid() error {
	isValidName := cfg.validateDeploymentName()
	if !isValidName {
		return fmt.Errorf("invalid name %s in deployment config", cfg.Name)
	}

	if cfg.Image == "" {
		return errors.New("image required in deployment config")
	}

	return nil
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

// DockerConfig: returns the configuration for creating a docker container
func (cfg DeploymentConfig) DockerConfig() docker.CreateContainerConfig {
	ctx := context.Background()
	defer ctx.Done()

	kraneNetwork, err := docker.GetClient().GetNetworkByName(ctx, docker.KraneNetworkName)
	if err != nil {
		return docker.CreateContainerConfig{}
	}

	var command []string
	var entrypoint []string

	if cfg.Command != "" {
		command = append(command, cfg.Command)
	}

	if cfg.Entrypoint != "" {
		entrypoint = append(entrypoint, cfg.Entrypoint)
	}

	containerName := fmt.Sprintf("%s-%s", cfg.Name, shortuuid.New())
	return docker.CreateContainerConfig{
		ContainerName: containerName,
		Image:         cfg.Image,
		NetworkID:     kraneNetwork.ID,
		Labels:        cfg.DockerLabels(),
		Ports:         cfg.DockerPorts(),
		VolumesMount:  cfg.DockerVolumesMount(),
		VolumesMap:    cfg.DockerVolumesMap(),
		Env:           cfg.DockerEnvs(),
		Command:       command,
		Entrypoint:    entrypoint,
	}
}

// DockerEnvs : returns a list of formatted Docker environment variables
func (cfg DeploymentConfig) DockerEnvs() []string {
	envs := make([]string, 0)

	// environment variables
	for k, v := range cfg.Env {
		envs = append(envs, fmt.Sprintf("%s=%s", k, v))
	}

	// resolve secrets by alias
	for key, alias := range cfg.Secrets {
		secret, err := secrets.Get(cfg.Name, alias)
		if err != nil || secret == nil {
			logger.Debugf("Unable to resolve secret for %s with alias @%s", cfg.Name, alias)
			continue
		}
		envs = append(envs, fmt.Sprintf("%s=%s", key, secret.Value))
	}

	return envs
}

// DockerVolumesMount : returns a list of formatted Docker volume mounts
func (cfg DeploymentConfig) DockerVolumesMount() []mount.Mount {
	volumes := make([]mount.Mount, 0)
	for hostVolume, containerVolume := range cfg.Volumes {
		volumes = append(volumes, mount.Mount{
			Type:   mount.TypeBind,
			Source: hostVolume,
			Target: containerVolume,
		})
	}
	return volumes
}

// DockerVolumesMap : returns a map of a formatted Docker volume map
func (cfg DeploymentConfig) DockerVolumesMap() map[string]struct{} {
	volumes := make(map[string]struct{}, 0)
	for _, containerVolume := range cfg.Volumes {
		volumes[containerVolume] = struct{}{}
	}
	return volumes
}

// DockerLabels : returns a map of Docker labels that should be applied to krane managed containers
func (cfg DeploymentConfig) DockerLabels() map[string]string {
	cfg.Labels[docker.ContainerNamespaceLabel] = cfg.Name
	cfg.ApplyProxyLabels()
	return cfg.Labels
}

// DockerPorts : returns a formatted Docker port map
func (cfg DeploymentConfig) DockerPorts() nat.PortMap {
	bindings := nat.PortMap{}
	for hostPort, containerPort := range cfg.Ports {
		if hostPort == "" {
			// randomly assign a host port if no explicit host port to bind to was provided
			freePort, err := getFreePort()
			if err != nil {
				logger.Errorf("Error when looking for a free host port %v", err)
				continue
			}
			hostPort = freePort
		}

		hostBinding := nat.PortBinding{HostPort: hostPort}

		// TODO: expose functionality for binding to other protocols besides tcp
		cPort, err := nat.NewPort("tcp", containerPort)
		if err != nil {
			logger.Errorf("Unable to create new container port %v", err)
			continue
		}

		bindings[cPort] = []nat.PortBinding{hostBinding}
	}

	return bindings
}

// getFreePort : find a free port on the host machine
func getFreePort() (string, error) {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		return "", err
	}

	listener, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return "", err
	}
	defer listener.Close()

	port := listener.Addr().(*net.TCPAddr).Port
	return strconv.Itoa(port), nil
}

// ApplyProxyLabels : applies network labels to a deployment config; currently only applies Traefik labels
func (cfg DeploymentConfig) ApplyProxyLabels() {
	// default labels
	cfg.Labels["traefik.enable"] = "true"
	cfg.Labels["traefik.docker.network"] = docker.KraneNetworkName

	// router labels
	for k, v := range proxy.TraefikRouterLabels(cfg.Name, cfg.Alias, cfg.Secured) {
		cfg.Labels[k] = v
	}

	// middleware labels
	for k, v := range proxy.TraefikMiddlewareLabels(cfg.Name, cfg.Secured) {
		cfg.Labels[k] = v
	}

	// service labels
	for k, v := range proxy.TraefikServiceLabels(cfg.Name, cfg.Ports) {
		cfg.Labels[k] = v
	}
}
