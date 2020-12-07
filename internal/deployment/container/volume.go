package container

import (
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/mount"

	"github.com/biensupernice/krane/internal/deployment/config"
)

type Volume struct {
	HostVolume      string `json:"host_volume"`
	ContainerVolume string `json:"container_volume"`
}

// fromKcontainerToDockerVolumeMount :
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

// fromMountPointToKconfigVolumes :
func fromMountPointToKconfigVolumes(mounts []types.MountPoint) []Volume {
	volumes := make([]Volume, 0)
	for _, m := range mounts {
		volumes = append(volumes, Volume{
			HostVolume:      m.Source,
			ContainerVolume: m.Destination,
		})
	}
	return volumes
}

// fromKconfigToDockerVolumeMount :
func fromKconfigToDockerVolumeMount(cfg config.DeploymentConfig) []mount.Mount {
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
