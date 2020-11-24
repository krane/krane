package proxy

import (
	"fmt"

	"github.com/biensupernice/krane/internal/docker"
	"github.com/biensupernice/krane/internal/utils"
)

// Flag to indicate containers should be
// labeled with security labels for secure communication over TLS
var Secured = utils.GetBoolEnv("SECURED")

func MakeContainerRoutingLabels(namespace, alias string) []ProxyLabel {
	labels := make([]ProxyLabel, 0)

	labels = append(labels, ProxyLabel{
		Label: "traefik.docker.network",
		Value: docker.KraneNetworkName,
	})

	// labels = append(labels, traefikEntryPointsLabels()...)
	// labels = append(labels, traefikMiddlewareLabels(namespace)...)
	labels = append(labels, traefikRouterLabels(namespace, alias)...)
	// labels = append(labels, traefikServiceLabels(namespace, alias)...)
	return labels
}

func traefikRouterLabels(namespace, alias string) []ProxyLabel {
	routerLabels := make([]ProxyLabel, 0)

	routerLabels = append(routerLabels, ProxyLabel{
		Label: fmt.Sprintf("traefik.http.routers.%s.rule", namespace),
		Value: fmt.Sprintf("Host(`%s`)", alias),
	})

	if Secured {
		routerLabels = append(routerLabels, ProxyLabel{
			Label: fmt.Sprintf("traefik.http.routers.%s.entrypoints", namespace),
			Value: "web-secure",
		})
		routerLabels = append(routerLabels, ProxyLabel{
			Label: fmt.Sprintf("traefik.http.routers.%s.tls", namespace),
			Value: "true",
		})
		routerLabels = append(routerLabels, ProxyLabel{
			Label: fmt.Sprintf("traefik.http.routers.%s.tls.certresolver", namespace),
			Value: "lets-encrypt",
		})
	}

	return routerLabels
}

func traefikServiceLabels(namespace, alias string) []ProxyLabel {
	serviceLabels := make([]ProxyLabel, 0)

	// serviceLabels = append(serviceLabels, ProxyLabel{
	// 	Label: "traefik.http.services.myservice.loadbalancer.server.port",
	// 	Value: "TODO",
	// })

	if Secured {
		serviceLabels = append(serviceLabels, ProxyLabel{
			Label: fmt.Sprintf("traefik.http.services.%s.loadbalancer.server.scheme", namespace),
			Value: "https",
		})
	}

	return serviceLabels
}

func traefikEntryPointsLabels() []ProxyLabel {
	entryPointLabels := make([]ProxyLabel, 0)

	if Secured {
		entryPointLabels = append(entryPointLabels, ProxyLabel{
			Label: "entryPoints.https.address",
			Value: "443",
		})
	}

	return entryPointLabels
}

func traefikMiddlewareLabels(namespace string) []ProxyLabel {
	middlewareLabels := make([]ProxyLabel, 0)

	if Secured {
		middlewareLabels = append(middlewareLabels, ProxyLabel{
			Label: fmt.Sprintf("traefik.http.middlewares.%s.redirectscheme.scheme", namespace),
			Value: "https",
		})

		middlewareLabels = append(middlewareLabels, ProxyLabel{
			Label: fmt.Sprintf("traefik.http.middlewares.%s.redirectscheme.port", namespace),
			Value: "443",
		})

		middlewareLabels = append(middlewareLabels, ProxyLabel{
			Label: fmt.Sprintf("traefik.http.middlewares.%s.redirectscheme.permanent", namespace),
			Value: "true",
		})
	}

	return middlewareLabels
}
