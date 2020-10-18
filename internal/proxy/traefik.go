package proxy

import (
	"fmt"

	"github.com/biensupernice/krane/internal/docker"
	"github.com/biensupernice/krane/internal/utils"
)

// Flag to indicate containers should be
// labeled with security labels for secure communication over TLS
var Secured = utils.GetBoolEnv("SECURED")

type TraefikLabel struct {
	Label string
	Value string
}

func MakeContainerRoutingLabels(namespace, alias string) []TraefikLabel {
	labels := make([]TraefikLabel, 0)

	labels = append(labels, TraefikLabel{
		Label: "proxy.enabled",
		Value: "true",
	})

	labels = append(labels, TraefikLabel{
		Label: "proxy.docker.network",
		Value: docker.KraneNetworkName,
	})

	labels = append(labels, traefikEntryPointsLabels()...)
	labels = append(labels, traefikMiddlewareLabels(namespace)...)
	labels = append(labels, traefikRouterLabels(namespace, alias)...)
	labels = append(labels, traefikServiceLabels(namespace, alias)...)
	return labels
}

func traefikRouterLabels(namespace, alias string) []TraefikLabel {
	routerLabels := make([]TraefikLabel, 0)

	routerLabels = append(routerLabels, TraefikLabel{
		Label: fmt.Sprintf("proxy.http.routers.%s.rule", namespace),
		Value: fmt.Sprintf("Host(`%s`)", alias),
	})

	if Secured {
		routerLabels = append(routerLabels, TraefikLabel{
			Label: fmt.Sprintf("proxy.http.routers.%s.entrypoints", namespace),
			Value: "websecure",
		})
	}

	return routerLabels
}

func traefikServiceLabels(namespace, alias string) []TraefikLabel {
	serviceLabels := make([]TraefikLabel, 0)
	return serviceLabels
}

func traefikEntryPointsLabels() []TraefikLabel {
	entryPointLabels := make([]TraefikLabel, 0)

	if Secured {
		entryPointLabels = append(entryPointLabels, TraefikLabel{
			Label: "entryPoints.websecure.address",
			Value: "443",
		})
	}

	return entryPointLabels
}

func traefikMiddlewareLabels(namespace string) []TraefikLabel {
	middlewareLabels := make([]TraefikLabel, 0)

	if Secured {
		middlewareLabels = append(middlewareLabels, TraefikLabel{
			Label: fmt.Sprintf("proxy.http.middlewares.%s.redirectscheme.scheme", namespace),
			Value: "https",
		})
	}

	return middlewareLabels
}
