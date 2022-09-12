package run

import (
	"fmt"
	"os"

	"github.com/containers/podman/v4/pkg/specgen"
	"github.com/julioln/sandman/config"
	specs "github.com/opencontainers/runtime-spec/specs-go"
)

func Pipewire(spec *specgen.SpecGenerator, containerConfig config.ContainerConfig) {
	if containerConfig.Run.Pipewire {
		spec.Env["XDG_RUNTIME_DIR"] = os.Getenv("XDG_RUNTIME_DIR")
		spec.Mounts = append(spec.Mounts,
			specs.Mount{
				Destination: fmt.Sprintf("%s/pipewire-0", os.Getenv("XDG_RUNTIME_DIR")),
				Source:      fmt.Sprintf("%s/pipewire-0", os.Getenv("XDG_RUNTIME_DIR")),
				Type:        "bind",
			},
		)
	}
}
