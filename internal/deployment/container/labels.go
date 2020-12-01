package container

import (
	"github.com/biensupernice/krane/internal/deployment/kconfig"
	"github.com/biensupernice/krane/internal/proxy"
)

// KraneContainerNamespaceLabel : container label for krane managed containers
const KraneContainerNamespaceLabel = "krane.deployment.namespace"

// fromKconfigToDockerLabelMap :
func fromKconfigToDockerLabelMap(cfg kconfig.Kconfig) map[string]string {
	labels := make(map[string]string, 0)
	labels[KraneContainerNamespaceLabel] = cfg.Name

	traefikLabels := proxy.CreateTraefikContainerLabels(cfg)
	for k, v := range traefikLabels {
		labels[k] = v
	}

	// combine custom labels with Config provided labels
	// Note: config labels have higher priority and will override any custom label
	for k, v := range cfg.Labels {
		labels[k] = v
	}

	return labels
}
