package run

import (
	"github.com/containers/podman/v5/pkg/specgen"
	"github.com/julioln/sandman/config"
	specs "github.com/opencontainers/runtime-spec/specs-go"
)

func Fonts(spec *specgen.SpecGenerator, containerConfig config.ContainerConfig) {
	if containerConfig.Run.Fonts {
		spec.Mounts = append(spec.Mounts, specs.Mount{
			Destination: "/usr/share/fonts",
			Source:      "/usr/share/fonts",
			Type:        "bind",
			Options:     []string{"ro"},
		})
	}
}
