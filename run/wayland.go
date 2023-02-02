package run

import (
	"fmt"
	"os"

	"github.com/containers/podman/v4/pkg/specgen"
	"github.com/julioln/sandman/config"
	specs "github.com/opencontainers/runtime-spec/specs-go"
)

func Wayland(spec *specgen.SpecGenerator, containerConfig config.ContainerConfig) {
	if containerConfig.Run.Wayland {
		spec.Env["WAYLAND_DISPLAY"] = os.Getenv("WAYLAND_DISPLAY")
		spec.Env["XDG_SESSION_TYPE"] = os.Getenv("XDG_SESSION_TYPE")
		spec.Mounts = append(spec.Mounts, specs.Mount{
			Destination: fmt.Sprintf("%s/%s", os.Getenv("XDG_RUNTIME_DIR"), os.Getenv("WAYLAND_DISPLAY")),
			Source:      fmt.Sprintf("%s/%s", os.Getenv("XDG_RUNTIME_DIR"), os.Getenv("WAYLAND_DISPLAY")),
			Type:        "bind",
		})
	}
}
