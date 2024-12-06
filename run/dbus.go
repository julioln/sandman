package run

import (
	"fmt"
	"os"

	"github.com/containers/podman/v5/pkg/specgen"
	"github.com/julioln/sandman/config"
	specs "github.com/opencontainers/runtime-spec/specs-go"
)

func Dbus(spec *specgen.SpecGenerator, containerConfig config.ContainerConfig) {
	if containerConfig.Run.Dbus {
		spec.Env["DBUS_SESSION_BUS_ADDRESS"] = fmt.Sprintf("unix:path=%s/bus", os.Getenv("XDG_RUNTIME_DIR"))
		spec.Mounts = append(spec.Mounts, specs.Mount{
			Destination: fmt.Sprintf("%s/bus", os.Getenv("XDG_RUNTIME_DIR")),
			Source:      fmt.Sprintf("%s/bus", os.Getenv("XDG_RUNTIME_DIR")),
			Type:        "bind",
		})
	}
}
