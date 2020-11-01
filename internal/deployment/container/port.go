package container

import (
	"strconv"

	"github.com/docker/docker/api/types"
	"github.com/docker/go-connections/nat"

	"github.com/biensupernice/krane/internal/deployment/config"
)

type Port struct {
	IP            string `json:"ip"`
	Type          string `json:"type"`
	HostPort      string `json:"host_port"`
	ContainerPort string `json:"container_port"`
}

type PortProtocol string

const (
	TCP PortProtocol = "tcp"
)

// from Kcontainer to Docker container port mapping
func fromKcontainerToDockerPortMap(ports []Port) (nat.PortMap, error) {
	portMap := nat.PortMap{}

	for _, port := range ports {
		hPort := nat.PortBinding{HostPort: port.HostPort}
		cPort, err := nat.NewPort(port.Type, port.ContainerPort)
		if err != nil {
			return nat.PortMap{}, err
		}

		portMap[cPort] = []nat.PortBinding{hPort}
	}

	return portMap, nil
}

func fromDockerToKcontainerPorts(ports []types.Port) []Port {
	kPorts := make([]Port, 0)
	for _, port := range ports {
		kPorts = append(kPorts, Port{
			Type:          port.Type,
			IP:            port.IP,
			HostPort:      strconv.Itoa(int(port.PublicPort)),
			ContainerPort: strconv.Itoa(int(port.PrivatePort)),
		})
	}
	return kPorts
}

// from Kconfig to Docker container port map
func fromKconfigToDockerPortMap(cfg config.Kconfig) nat.PortMap {
	bindings := nat.PortMap{}
	for hostPort, containerPort := range cfg.Ports {
		// host port
		hostBinding := nat.PortBinding{HostPort: hostPort}

		// container port
		// TODO: figure out if we can bind ports of other types besides tcp
		cPort, err := nat.NewPort("tcp", containerPort)
		if err != nil {
			continue
		}

		bindings[cPort] = []nat.PortBinding{hostBinding}
	}

	return bindings
}

// from Docker container port mapping to Kcontainer port map
func fromContainerPortMap(ports types.Port) map[string]string {
	portMap := make(map[string]string)
	return portMap
}
