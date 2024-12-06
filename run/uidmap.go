package run

import (
	"fmt"
	"os"

	"github.com/containers/podman/v5/pkg/specgen"
	"github.com/containers/storage/pkg/idtools"
	"github.com/containers/storage/types"
	"github.com/julioln/sandman/config"
)

func Uidmap(spec *specgen.SpecGenerator, containerConfig config.ContainerConfig) {
	if containerConfig.Run.Uidmap {
		idMaps := []idtools.IDMap{
			{
				ContainerID: os.Getuid(),
				HostID:      0,
				Size:        1,
			},
			{
				ContainerID: 0,
				HostID:      1,
				Size:        os.Getuid(),
			},
			{
				ContainerID: os.Getuid() + 1,
				HostID:      os.Getuid() + 1,
				Size:        65536 - os.Getuid(),
			},
		}
		var idMappingOptions types.IDMappingOptions
		idMappingOptions.HostUIDMapping = true
		idMappingOptions.HostGIDMapping = true
		idMappingOptions.UIDMap = append(idMappingOptions.UIDMap, idMaps...)
		idMappingOptions.GIDMap = append(idMappingOptions.GIDMap, idMaps...)
		spec.IDMappings = &idMappingOptions
		spec.UserNS.NSMode = specgen.Private
		spec.User = fmt.Sprint(os.Getuid())
	}

}
