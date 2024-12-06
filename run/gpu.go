package run

import (
	"github.com/containers/podman/v5/pkg/specgen"
	"github.com/julioln/sandman/config"
	specs "github.com/opencontainers/runtime-spec/specs-go"
)

func Gpu(spec *specgen.SpecGenerator, containerConfig config.ContainerConfig) {
	// Expose Direct Render Inteface for GPU acceleration
	if containerConfig.Run.Dri || containerConfig.Run.Gpu {
		spec.Devices = append(spec.Devices, specs.LinuxDevice{
			Path: "/dev/dri",
		})
	}
}
