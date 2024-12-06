package run

import (
	"os"

	"github.com/containers/podman/v5/pkg/specgen"
	"github.com/julioln/sandman/config"
	specs "github.com/opencontainers/runtime-spec/specs-go"
)

func X11(spec *specgen.SpecGenerator, containerConfig config.ContainerConfig) {
	if containerConfig.Run.X11 {
		spec.Env["DISPLAY"] = os.Getenv("DISPLAY")
		spec.Env["XDG_SESSION_TYPE"] = os.Getenv("XDG_SESSION_TYPE")
		spec.Env["XCURSOR_THEME"] = os.Getenv("XCURSOR_THEME")
		spec.Env["XCURSOR_SIZE"] = os.Getenv("XCURSOR_SIZE")
		spec.Mounts = append(spec.Mounts, specs.Mount{
			Destination: "/tmp/.X11-unix",
			Source:      "/tmp/.X11-unix",
			Type:        "bind",
		})
	}
}
