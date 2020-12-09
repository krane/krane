package proxy

import (
	"bytes"
	"fmt"

	"github.com/biensupernice/krane/internal/deployment/config"
	"github.com/biensupernice/krane/internal/docker"
	"github.com/biensupernice/krane/internal/proxy/middlewares"
)

type TraefikLabel struct {
	Key   string
	Value string
}

func CreateTraefikContainerLabels(config config.DeploymentConfig) map[string]string {
	labels := make(map[string]string, 0)

	// default labels
	labels["traefik.enable"] = "true"
	labels["traefik.docker.network"] = docker.KraneNetworkName

	// router labels
	for k, v := range traefikRouterLabels(config.Name, config.Alias, config.Secured) {
		labels[k] = v
	}

	// middleware labels
	for k, v := range traefikMiddlewareLabels(config.Name, config.Secured) {
		labels[k] = v
	}

	// service labels
	for k, v := range traefikServiceLabels(config.Name, config.Ports) {
		labels[k] = v
	}

	return labels
}

func traefikRouterLabels(namespace string, aliases []string, secured bool) map[string]string {
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

func traefikServiceLabels(namespace string, ports map[string]string) map[string]string {
	labels := make(map[string]string, 0)
	for _, containerPort := range ports {
		labels[fmt.Sprintf("traefik.http.services.%s-%s.loadbalancer.server.port", namespace, containerPort)] = containerPort
		labels[fmt.Sprintf("traefik.http.services.%s-%s.loadbalancer.server.scheme", namespace, containerPort)] = "http"
	}
	return labels
}

func traefikMiddlewareLabels(namespace string, secured bool) map[string]string {
	labels := make(map[string]string, 0)
	if secured {
		// applies http redirect labels to all secure deployments
		for k, v := range middlewares.RedirectToHTTPSLabels(namespace) {
			labels[k] = v
		}
	}
	return labels
}
