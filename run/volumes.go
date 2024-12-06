package run

import (
	"strings"

	"github.com/containers/podman/v5/pkg/specgen"
	"github.com/julioln/sandman/config"
	specs "github.com/opencontainers/runtime-spec/specs-go"
)

func Volumes(spec *specgen.SpecGenerator, containerConfig config.ContainerConfig) {
	// Mount additional volumes
	for _, volume := range containerConfig.Run.Volumes {
		v := strings.Split(volume, ":")
		var dest string
		var src string
		var mountOptions []string

		if len(v) < 2 {
			// Shorthand
			src = v[0]
			dest = v[0]
		} else {
			src = v[0]
			dest = v[1]
		}

		if len(v) > 2 {
			// mount -o like arguments
			mountOptions = strings.Split(v[2], ",")
		}

		spec.Mounts = append(spec.Mounts, specs.Mount{
			Destination: dest,
			Source:      src,
			Type:        "bind",
			Options:     mountOptions,
		})
	}
}
