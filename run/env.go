package run

import (
	"os"
	"strings"

	"github.com/containers/podman/v6/pkg/specgen"
	"github.com/julioln/sandman/config"
)

func Env(spec *specgen.SpecGenerator, containerConfig config.ContainerConfig) {
	for _, env := range containerConfig.Run.Env {
		e := strings.Split(env, "=")
		var k string = e[0]
		var v string
		if len(e) == 1 {
			// Same behavior as command line
			v = os.Getenv(k)
		} else {
			v = e[1]
		}
		spec.Env[k] = v
	}
}
