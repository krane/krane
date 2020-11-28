package proxy

import (
	"fmt"

	"github.com/biensupernice/krane/internal/deployment/kconfig"
	"github.com/biensupernice/krane/internal/docker"
)

type ProxyLabel struct {
	Key   string
	Value string
}

func MakeContainerRoutingLabels(config kconfig.Kconfig) []ProxyLabel {
	labels := make([]ProxyLabel, 0)

	labels = append(labels, ProxyLabel{
		Key:   "traefik.enable",
		Value: "true",
	})

	labels = append(labels, ProxyLabel{
		Key:   "traefik.docker.network",
		Value: docker.KraneNetworkName,
	})

	labels = append(labels, traefikEntryPointsLabels(config.Secured)...)
	labels = append(labels, traefikMiddlewareLabels(config.Name, config.Secured)...)
	labels = append(labels, traefikRouterLabels(config.Name, config.Alias, config.Secured)...)
	labels = append(labels, traefikServiceLabels(config.Name, config.Ports, config.Secured)...)
	return labels
}

func traefikRouterLabels(namespace string, aliases []string, secured bool) []ProxyLabel {
	routerLabels := make([]ProxyLabel, 0)

	for i, alias := range aliases {
		routerLabels = append(routerLabels, ProxyLabel{
			Key:   fmt.Sprintf("traefik.http.routers.%s-%d.rule", namespace, i),
			Value: fmt.Sprintf("Host(`%s`)", alias),
		})
	}

	routerLabels = append(routerLabels, ProxyLabel{
		Key:   fmt.Sprintf("traefik.http.routers.%s.entrypoints", namespace),
		Value: "web",
	})

	if secured {
		routerLabels = append(routerLabels, ProxyLabel{
			Key:   fmt.Sprintf("traefik.http.routers.%s-secured.entrypoints", namespace),
			Value: "web-secure",
		})
		routerLabels = append(routerLabels, ProxyLabel{
			Key:   fmt.Sprintf("traefik.http.routers.%s-secured.tls", namespace),
			Value: "true",
		})
		routerLabels = append(routerLabels, ProxyLabel{
			Key:   fmt.Sprintf("traefik.http.routers.%s-secured.tls.certresolver", namespace),
			Value: "lets-encrypt",
		})
	}

	return routerLabels
}

func traefikServiceLabels(namespace string, ports map[string]string, secured bool) []ProxyLabel {
	serviceLabels := make([]ProxyLabel, 0)

	i := 0
	for _, containerPort := range ports {
		serviceLabels = append(serviceLabels, ProxyLabel{
			Key:   fmt.Sprintf("traefik.http.services.%s-%d.loadbalancer.server.port", namespace, i),
			Value: containerPort,
		})
		i++
	}

	if secured {
		serviceLabels = append(serviceLabels, ProxyLabel{
			Key:   fmt.Sprintf("traefik.http.services.%s-secured.loadbalancer.server.scheme", namespace),
			Value: "https",
		})
	}

	return serviceLabels
}

func traefikEntryPointsLabels(secured bool) []ProxyLabel {
	entryPointLabels := make([]ProxyLabel, 0)

	entryPointLabels = append(entryPointLabels, ProxyLabel{
		Key:   "entryPoints.web.address",
		Value: "80",
	})

	if secured {
		entryPointLabels = append(entryPointLabels, ProxyLabel{
			Key:   "entryPoints.web-secure.address",
			Value: "443",
		})
	}

	return entryPointLabels
}

func traefikMiddlewareLabels(namespace string, secured bool) []ProxyLabel {
	middlewareLabels := make([]ProxyLabel, 0)

	if secured {
		middlewareLabels = append(middlewareLabels, ProxyLabel{
			Key:   fmt.Sprintf("traefik.http.middlewares.%s.redirectscheme.scheme", namespace),
			Value: "https",
		})

		middlewareLabels = append(middlewareLabels, ProxyLabel{
			Key:   fmt.Sprintf("traefik.http.middlewares.%s.redirectscheme.port", namespace),
			Value: "443",
		})

		middlewareLabels = append(middlewareLabels, ProxyLabel{
			Key:   fmt.Sprintf("traefik.http.middlewares.%s.redirectscheme.permanent", namespace),
			Value: "true",
		})
	}

	return middlewareLabels
}
