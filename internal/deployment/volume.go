package deployment

import "github.com/docker/docker/api/types"

type Volume struct {
	HostVolume      string `json:"host_volume"`
	ContainerVolume string `json:"container_volume"`
}

// fromMountPointToVolumeList converts a list of volume MountPoints into a list of formatted Krane Volumes
func fromMountPointToVolumeList(mounts []types.MountPoint) []Volume {
	volumes := make([]Volume, 0)
	for _, m := range mounts {
		volumes = append(volumes, Volume{
			HostVolume:      m.Source,
			ContainerVolume: m.Destination,
		})
	}
	return volumes
}
