package container

import (
	"github.com/biensupernice/krane/internal/deployment/config"
	"github.com/biensupernice/krane/internal/proxy"
)

// KraneContainerLabel : container label to identify krane managed containers
const KraneContainerLabel = "krane.deployment"

// fromKconfigToDockerLabelMap :
func fromKconfigToDockerLabelMap(cfg config.DeploymentConfig) map[string]string {
	labels := make(map[string]string, 0)
	labels[KraneContainerLabel] = cfg.Name

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
