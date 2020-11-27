package container

import (
	"net"
	"strconv"

	"github.com/docker/docker/api/types"
	"github.com/docker/go-connections/nat"
	"github.com/sirupsen/logrus"

	"github.com/biensupernice/krane/internal/deployment/kconfig"
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

func fromDockerToKconfigPortMap(pMap nat.PortMap) []Port {
	bindings := make([]Port, 0)

	for container, hostBindings := range pMap {
		for _, hostB := range hostBindings {
			bindings = append(bindings, Port{
				IP:            hostB.HostIP,
				HostPort:      hostB.HostPort,
				Type:          container.Proto(),
				ContainerPort: container.Port(),
			})
		}
	}

	return bindings
}

// from Kconfig to Docker container port map
func fromKconfigToDockerPortMap(cfg kconfig.Kconfig) nat.PortMap {
	bindings := nat.PortMap{}
	for hostPort, containerPort := range cfg.Ports {
		if hostPort == "" {
			// randomly assign a host port if no explicit host port to bind to was provided
			freePort, err := getFreePort()
			if err != nil {
				logrus.Errorf("Error when looking for a free host port %v", err)
				continue
			}
			hostPort = freePort
		}

		hostBinding := nat.PortBinding{HostPort: hostPort}

		// TODO: figure out if we can bind ports of other types besides tcp
		protocol := string(TCP)
		cPort, err := nat.NewPort(protocol, containerPort)
		if err != nil {
			logrus.Errorf("Unable to create new container port %v", err)
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

// asks the kernel for a free open port that is ready to use.
func getFreePort() (string, error) {
	addr, err := net.ResolveTCPAddr(string(TCP), "localhost:0")
	if err != nil {
		return "", err
	}

	listener, err := net.ListenTCP(string(TCP), addr)
	if err != nil {
		return "", err
	}
	defer listener.Close()

	port := listener.Addr().(*net.TCPAddr).Port
	return strconv.Itoa(port), nil
}
