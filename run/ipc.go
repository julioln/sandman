package run

import (
	"github.com/containers/podman/v4/pkg/specgen"
	"github.com/julioln/sandman/config"
)

func Ipc(spec *specgen.SpecGenerator, containerConfig config.ContainerConfig) {
	// Share IPC with host
	if containerConfig.Run.Ipc {
		spec.IpcNS.NSMode = specgen.NamespaceMode("host")
	}
}
