package run

import (
	"context"
	"fmt"

	"github.com/julioln/sandman/config"
	"github.com/julioln/sandman/podman"
	//"github.com/containers/podman/v4/pkg/bindings"
	//"github.com/containers/podman/v4/pkg/bindings/containers"
	//"github.com/containers/podman/v4/pkg/bindings/images"
	//"github.com/containers/podman/v4/pkg/domain/entities"
)

func Run(socket string, containerConfig config.ContainerConfig, keep bool, verbose bool) {
	var conn context.Context = podman.InitializePodman(socket)

	if verbose {
		fmt.Println("Container Config: ", containerConfig)
		fmt.Println("Connection: ", conn)
	}
}

func CmdExecute(socket string, verbose bool, keep bool, args []string) {
	var container_name string = args[0]
	Run(socket, config.LoadConfig(container_name), keep, verbose)
}
