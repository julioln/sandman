package run

import (
	"github.com/containers/podman/v6/pkg/specgen"
	"github.com/julioln/sandman/config"
)

func Raw(spec *specgen.SpecGenerator, containerConfig config.ContainerConfig) {
	// Add raw configuration
	spec.PortMappings = append(spec.PortMappings, containerConfig.Run.RawPorts...)
	spec.Mounts = append(spec.Mounts, containerConfig.Run.RawMounts...)
	spec.Devices = append(spec.Devices, containerConfig.Run.RawDevices...)
}
