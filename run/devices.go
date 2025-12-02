package run

import (
	"github.com/containers/podman/v6/pkg/specgen"
	"github.com/julioln/sandman/config"
	specs "github.com/opencontainers/runtime-spec/specs-go"
)

func Devices(spec *specgen.SpecGenerator, containerConfig config.ContainerConfig) {
	for _, dev := range containerConfig.Run.Devices {
		spec.Devices = append(spec.Devices, specs.LinuxDevice{
			Path: dev,
		})
	}
}
