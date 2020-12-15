package proxy

import (
	"bytes"
	"fmt"

	"github.com/biensupernice/krane/internal/proxy/middlewares"
)

type TraefikLabel struct {
	Key   string
	Value string
}

func TraefikRouterLabels(namespace string, aliases []string, secured bool) map[string]string {
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
		labels[fmt.Sprintf("traefik.http.routers.%s-insecure.rule", namespace)] = hostRules.String()
	}
	labels[fmt.Sprintf("traefik.http.routers.%s-insecure.entrypoints", namespace)] = "web"

	if secured {
		// https
		labels[fmt.Sprintf("traefik.http.routers.%s-secure.tls", namespace)] = "true"
		labels[fmt.Sprintf("traefik.http.routers.%s-secure.entrypoints", namespace)] = "web-secure"
		labels[fmt.Sprintf("traefik.http.routers.%s-secure.tls.certresolver", namespace)] = "lets-encrypt"
		if hostRules.String() != "" {
			labels[fmt.Sprintf("traefik.http.routers.%s-secure.rule", namespace)] = hostRules.String()
		}
	}

	return labels
}

func TraefikServiceLabels(namespace string, ports map[string]string) map[string]string {
	labels := make(map[string]string, 0)
	for _, containerPort := range ports {
		labels[fmt.Sprintf("traefik.http.services.%s-%s.loadbalancer.server.port", namespace, containerPort)] = containerPort
		labels[fmt.Sprintf("traefik.http.services.%s-%s.loadbalancer.server.scheme", namespace, containerPort)] = "http"
	}
	return labels
}

func TraefikMiddlewareLabels(namespace string, secured bool) map[string]string {
	labels := make(map[string]string, 0)
	if secured {
		// applies http redirect labels to all secure deployments
		for k, v := range middlewares.RedirectToHTTPSLabels(namespace) {
			labels[k] = v
		}
	}
	return labels
}
