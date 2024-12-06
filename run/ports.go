package run

import (
	"fmt"
	"strconv"
	"strings"

	nettypes "github.com/containers/common/libnetwork/types"
	"github.com/containers/podman/v5/pkg/specgen"
	"github.com/julioln/sandman/config"
)

func Ports(spec *specgen.SpecGenerator, containerConfig config.ContainerConfig) {
	for _, ports := range containerConfig.Run.Ports {
		p := strings.Split(ports, ":")
		if len(p) < 2 {
			fmt.Println("Invalid port configuration, ignoring: ", p)
			continue
		}

		containerPort, _ := strconv.Atoi(p[0])
		hostPort, _ := strconv.Atoi(p[1])

		spec.PortMappings = append(spec.PortMappings, nettypes.PortMapping{
			ContainerPort: uint16(containerPort),
			HostPort:      uint16(hostPort),
		})
	}
}
