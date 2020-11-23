package container

import (
	"github.com/docker/docker/api/types/mount"

	"github.com/biensupernice/krane/internal/deployment/kconfig"
)

type Volume struct {
	HostVolume      string `json:"host_volume"`
	ContainerVolume string `json:"container_volume"`
}

// from Kcontainer to Docker volume mounts
func fromKcontainerToDockerVolumeMount(volumes []Volume) []mount.Mount {
	vols := make([]mount.Mount, 0)
	for _, v := range volumes {
		vols = append(vols, mount.Mount{
			Type:   mount.TypeBind,
			Source: v.HostVolume,
			Target: v.ContainerVolume,
		})
	}

	return vols
}

// from Kconfig to Docker container volume mounts
func fromKconfigToDockerVolumeMount(cfg kconfig.Kconfig) []mount.Mount {
	volumes := make([]mount.Mount, 0)
	for hostVolume, containerVolume := range cfg.Volumes {
		volumes = append(volumes, mount.Mount{
			Type:   mount.TypeBind,
			Source: hostVolume,
			Target: containerVolume,
		})
	}

	return volumes
}
