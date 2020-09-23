package container

import "github.com/docker/go-connections/nat"

type Port struct {
	IP            string `json:"ip"`
	Type          string `json:"type"`
	HostPort      string `json:"container_port"`
	ContainerPort string `json:"container_port"`
}

type PortProtocol string

const (
	TCP PortProtocol = "tcp"
)

func makePortBindings(ports []Port) (nat.PortMap, error) {
	portMap := nat.PortMap{}

	for _, port := range ports {
		hPort := nat.PortBinding{HostPort: port.HostPort}
		cPort, err := nat.NewPort(string(port.Type), port.ContainerPort)
		if err != nil {
			return nat.PortMap{}, err
		}

		portMap[cPort] = []nat.PortBinding{hPort}
	}

	return portMap, nil
}
