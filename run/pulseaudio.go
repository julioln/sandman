package run

import (
	"fmt"
	"os"

	"github.com/containers/podman/v5/pkg/specgen"
	"github.com/julioln/sandman/config"
	specs "github.com/opencontainers/runtime-spec/specs-go"
)

func Pulseaudio(spec *specgen.SpecGenerator, containerConfig config.ContainerConfig) {
	if containerConfig.Run.Pulseaudio {
		spec.Env["XDG_RUNTIME_DIR"] = os.Getenv("XDG_RUNTIME_DIR")
		spec.Mounts = append(spec.Mounts,
			specs.Mount{
				Destination: "/etc/machine-id",
				Source:      "/etc/machine-id",
				Type:        "bind",
				Options:     []string{"ro"},
			},
			specs.Mount{
				Destination: fmt.Sprintf("%s/pulse/native", os.Getenv("XDG_RUNTIME_DIR")),
				Source:      fmt.Sprintf("%s/pulse/native", os.Getenv("XDG_RUNTIME_DIR")),
				Type:        "bind",
			},
		)
	}
}
