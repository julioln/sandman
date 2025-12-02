package run

import (
	"github.com/containers/podman/v6/pkg/specgen"
	"github.com/julioln/sandman/config"
	specs "github.com/opencontainers/runtime-spec/specs-go"
)

func Limits(spec *specgen.SpecGenerator, containerConfig config.ContainerConfig) {
	var limits specs.LinuxResources
	limits.CPU = &containerConfig.Run.Limits.CPU
	limits.Memory = &containerConfig.Run.Limits.Memory
	limits.Pids = &containerConfig.Run.Limits.Pids
	spec.ResourceLimits = &limits
	spec.Rlimits = containerConfig.Run.Limits.Rlimits
	spec.CgroupConf = containerConfig.Run.Limits.CgroupConf
}
