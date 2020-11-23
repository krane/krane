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

	// TODO: theres a bug where it only applies a single label if aliases is greater than 1.
	// This is because the labels key get overwritten
	for _, alias := range cfg.Alias {
		routingLabels := proxy.MakeContainerRoutingLabels(cfg.Name, alias)
		for _, label := range routingLabels {
			labels[label.Label] = label.Value
		}
	}

	return labels
}