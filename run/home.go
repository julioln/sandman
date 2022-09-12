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
			spec.Mounts = append(spec.Mounts, specs.Mount{
				Destination: "/home",
				Source:      mountPoint,
				Type:        "bind",
			})
		}
	}
}
