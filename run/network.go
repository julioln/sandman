package run

import (
	"fmt"

	"github.com/containers/podman/v4/pkg/specgen"
	"github.com/julioln/sandman/config"
)

func Network(spec *specgen.SpecGenerator, containerConfig config.ContainerConfig) {
	// Configure network namespace
	var networkNS specgen.Namespace
	if containerConfig.Run.Network == "" && !containerConfig.Run.Net {
		networkNS.NSMode = specgen.None
	} else if containerConfig.Run.Net {
		// Backwards compatibility
		networkNS.NSMode = specgen.Slirp
	} else {
		var err error
		if networkNS, _, _, err = specgen.ParseNetworkFlag([]string{containerConfig.Run.Network}, true); err != nil {
			fmt.Println("Error parsing network, defaulting to none: ", err)
			networkNS.NSMode = specgen.None
		}
	}
	spec.NetNS = networkNS
}
