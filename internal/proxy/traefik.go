package proxy

import (
	"bytes"
	"fmt"

	"github.com/biensupernice/krane/internal/deployment/kconfig"
	"github.com/biensupernice/krane/internal/docker"
)

type TraefikLabel struct {
	Key   string
	Value string
}

func MakeContainerRoutingLabels(config kconfig.Kconfig) []TraefikLabel {
	labels := make([]TraefikLabel, 0)

	labels = append(labels, TraefikLabel{
		Key:   "traefik.enable",
		Value: "true",
	})

	labels = append(labels, TraefikLabel{
		Key:   "traefik.docker.network",
		Value: docker.KraneNetworkName,
	})

	labels = append(labels, traefikRouterLabels(config.Name, config.Alias, config.Secured)...)
	labels = append(labels, traefikMiddlewareLabels(config.Name, config.Secured)...)
	labels = append(labels, traefikServiceLabels(config.Name, config.Ports)...)
	return labels
}

func traefikRouterLabels(namespace string, aliases []string, secured bool) []TraefikLabel {
	routerLabels := make([]TraefikLabel, 0)

	// configure aliases as Host labels
	var hostLabels bytes.Buffer
	for i, alias := range aliases {
		if i == len(aliases)-1 {
			// if last or only alias, just append the host with no OR operator
			hostLabels.WriteString(fmt.Sprintf("Host(`%s`)", alias))
		} else {
			// append OR operator if not the last alias
			hostLabels.WriteString(fmt.Sprintf("Host(`%s`) ||", alias))
		}
	}

	// http
	routerLabels = append(routerLabels, TraefikLabel{
		Key:   fmt.Sprintf("traefik.http.routers.%s-insecure.rule", namespace),
		Value: hostLabels.String(),
	})

	routerLabels = append(routerLabels, TraefikLabel{
		Key:   fmt.Sprintf("traefik.http.routers.%s-insecure.entrypoints", namespace),
		Value: "web",
	})

	// https
	if secured {
		routerLabels = append(routerLabels, TraefikLabel{
			Key:   fmt.Sprintf("traefik.http.routers.%s-secure.rule", namespace),
			Value: hostLabels.String(),
		})

		routerLabels = append(routerLabels, TraefikLabel{
			Key:   fmt.Sprintf("traefik.http.routers.%s-secure.tls", namespace),
			Value: "true",
		})

		routerLabels = append(routerLabels, TraefikLabel{
			Key:   fmt.Sprintf("traefik.http.routers.%s-secure.entrypoints", namespace),
			Value: "web-secure",
		})

		routerLabels = append(routerLabels, TraefikLabel{
			Key:   fmt.Sprintf("traefik.http.routers.%s-secure.tls.certresolver", namespace),
			Value: "lets-encrypt",
		})
	}

	return routerLabels
}

func traefikServiceLabels(namespace string, ports map[string]string) []TraefikLabel {
	serviceLabels := make([]TraefikLabel, 0)

	i := 0
	for _, containerPort := range ports {
		serviceLabels = append(serviceLabels, TraefikLabel{
			Key:   fmt.Sprintf("traefik.http.services.%s-%d.loadbalancer.server.port", namespace, i),
			Value: containerPort,
		})

		serviceLabels = append(serviceLabels, TraefikLabel{
			Key:   fmt.Sprintf("traefik.http.services.%s-%d.loadbalancer.server.scheme", namespace, i),
			Value: "http",
		})
		i++
	}

	return serviceLabels
}

func traefikMiddlewareLabels(namespace string, secured bool) []TraefikLabel {
	middlewareLabels := make([]TraefikLabel, 0)

	// TODO: could potentially expose a flag to enable these individually
	// - redirect
	// - compress
	// - ratelimit
	if secured {
		middlewareLabels = append(middlewareLabels, TraefikLabel{
			Key:   fmt.Sprintf("traefik.http.routers.%s-insecure.middlewares", namespace),
			Value: "redirect-to-https@docker",
		})

		middlewareLabels = append(middlewareLabels, TraefikLabel{
			Key:   "traefik.http.middlewares.redirect-to-https.redirectscheme.scheme",
			Value: "https",
		})

		middlewareLabels = append(middlewareLabels, TraefikLabel{
			Key:   "traefik.http.middlewares.redirect-to-https.redirectscheme.port",
			Value: "443",
		})

		middlewareLabels = append(middlewareLabels, TraefikLabel{
			Key:   "traefik.http.middlewares.redirect-to-https.redirectscheme.permanent",
			Value: "true",
		})
	}

	return middlewareLabels
}
