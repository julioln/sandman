package run

import (
	"github.com/containers/podman/v5/pkg/specgen"
	"github.com/julioln/sandman/config"
)

func Name(spec *specgen.SpecGenerator, containerConfig config.ContainerConfig) {
	if containerConfig.Run.Name != "" {
		spec.Name = containerConfig.Run.Name
		spec.Hostname = containerConfig.Run.Name
	}
}
