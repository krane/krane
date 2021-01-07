package proxy

import (
	"bytes"
	"fmt"

	"github.com/krane/krane/internal/proxy/middlewares"
)

type TraefikLabel struct {
	Key   string
	Value string
}

func TraefikRouterLabels(deployment string, aliases []string, secure bool) map[string]string {
	// configure aliases as Host('my-alias.example.com') labels
	var hostRules bytes.Buffer
	for i, alias := range aliases {
		if alias == "" {
			continue
		}

		if i == len(aliases)-1 {
			// if last or only alias, just append the host with no OR operator
			hostRules.WriteString(fmt.Sprintf("Host(`%s`)", alias))
		} else {
			// append OR operator if not the last alias
			hostRules.WriteString(fmt.Sprintf("Host(`%s`) || ", alias))
		}
	}

	labels := make(map[string]string, 0)

	// http
	if hostRules.String() != "" {
		labels[fmt.Sprintf("traefik.http.routers.%s-insecure.rule", deployment)] = hostRules.String()
	}
	labels[fmt.Sprintf("traefik.http.routers.%s-insecure.entrypoints", deployment)] = "web"

	if secure {
		// https
		labels[fmt.Sprintf("traefik.http.routers.%s-secure.tls", deployment)] = "true"
		labels[fmt.Sprintf("traefik.http.routers.%s-secure.entrypoints", deployment)] = "web-secure"
		labels[fmt.Sprintf("traefik.http.routers.%s-secure.tls.certresolver", deployment)] = "lets-encrypt"
		if hostRules.String() != "" {
			labels[fmt.Sprintf("traefik.http.routers.%s-secure.rule", deployment)] = hostRules.String()
		}
	}

	return labels
}

func TraefikServiceLabels(deployment string, ports map[string]string, targetPort string) map[string]string {
	labels := make(map[string]string, 0)

	if targetPort != "" {
		labels[fmt.Sprintf("traefik.http.services.%s.loadbalancer.server.port", deployment)] = targetPort
		labels[fmt.Sprintf("traefik.http.services.%s.loadbalancer.server.scheme", deployment)] = "http"
	} else {
		for _, containerPort := range ports {
			labels[fmt.Sprintf("traefik.http.services.%s-%s.loadbalancer.server.port", deployment, containerPort)] = containerPort
			labels[fmt.Sprintf("traefik.http.services.%s-%s.loadbalancer.server.scheme", deployment, containerPort)] = "http"
		}
	}

	return labels
}

func TraefikMiddlewareLabels(deployment string, secured bool) map[string]string {
	labels := make(map[string]string, 0)
	if secured {
		// applies http redirect labels to all secure deployments
		for k, v := range middlewares.RedirectToHTTPSLabels(deployment) {
			labels[k] = v
		}
	}
	return labels
}
