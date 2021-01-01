package deployment

import (
	"net"
	"strconv"

	"github.com/docker/go-connections/nat"
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

func fromPortMapToPortList(pMap nat.PortMap) []Port {
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

// findFreePort returns a free port on the host machine
func findFreePort() (string, error) {
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
