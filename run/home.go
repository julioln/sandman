package run

import (
	"fmt"
	"os"

	"github.com/containers/podman/v4/pkg/specgen"
	"github.com/julioln/sandman/config"
	specs "github.com/opencontainers/runtime-spec/specs-go"
)

func Home(spec *specgen.SpecGenerator, containerConfig config.ContainerConfig) {
	if containerConfig.Run.Home {
		var mountPoint = fmt.Sprintf("%s/%s", config.GetHomeStorageDir(), containerConfig.Name)
		if err := os.MkdirAll(mountPoint, 0755); err == nil {
			// Allow destination to be overriden
			var destination = "/home/user"
			if containerConfig.Run.HomePath != "" {
				destination = containerConfig.Run.HomePath
			}

			spec.Mounts = append(spec.Mounts, specs.Mount{
				Destination: destination,
				Source:      mountPoint,
				Type:        "bind",
			})
		}
	}
}
