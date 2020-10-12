package container

import "github.com/docker/docker/api/types/mount"

type Volume struct {
		HostVolume      string `json:"host_volume"`
		ContainerVolume string `json:"container_volume"`
}

func makeVolumes(volumes []Volume) []mount.Mount {
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
