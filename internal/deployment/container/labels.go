package container

import (
	"github.com/biensupernice/krane/internal/deployment/kconfig"
	"github.com/biensupernice/krane/internal/proxy"
)

// from Kconfig to Docker label map
func fromKconfigToDockerLabelMap(cfg kconfig.Kconfig) map[string]string {
	labels := map[string]string{
		KraneContainerNamespaceLabel: cfg.Name,
	}

	routingLabels := proxy.MakeContainerRoutingLabels(cfg)
	for _, label := range routingLabels {
		labels[label.Key] = label.Value
	}

	return labels
}
