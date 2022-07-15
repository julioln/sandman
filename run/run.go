package run

import (
	"context"
	"fmt"
	"os"

	"github.com/julioln/sandman/config"
	"github.com/julioln/sandman/podman"

	"github.com/containers/podman/v4/pkg/bindings/containers"
	"github.com/containers/podman/v4/pkg/specgen"
)

func Start(socket string, containerConfig config.ContainerConfig, keep bool, verbose bool, runCmd []string) {
	var conn context.Context = podman.InitializePodman(socket)
	var spec = CreateSpec(containerConfig)

	// Check overrides
	if keep {
		spec.Remove = false
	}
	if len(runCmd) > 0 {
		spec.Entrypoint = runCmd
	}

	if verbose {
		fmt.Println("Container Config: ", containerConfig)
		fmt.Println("Connection: ", conn)
	}

	var createOptions containers.CreateOptions
	container, err := containers.CreateWithSpec(conn, spec, &createOptions)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if verbose {
		fmt.Println("Container: ", container)
	}

	err = containers.Start(conn, container.ID, nil)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	//var stopOptions containers.StopOptions
	//defer containers.Stop(conn, container.ID, &stopOptions)

	if verbose {
		var inspectOptions containers.InspectOptions
		containerData, err := containers.Inspect(conn, container.ID, &inspectOptions)
		if err == nil {
			fmt.Println("Container data: ", containerData)
			fmt.Println("Exec IDs: ", containerData.ExecIDs)
			fmt.Println("Container Config: ", containerData.Config)
			fmt.Println("Container Config TTY: ", containerData.Config.Tty)
		}
	}

	// TODO: Work on auto-attaching

	//var attachOptions containers.AttachOptions
	//var attachReady chan bool
	//attachReady <- true
	//err = containers.Attach(
	//	conn,
	//	container.ID,
	//	io.Reader(os.Stdin),
	//	io.Writer(os.Stdout),
	//	io.Writer(os.Stderr),
	//	attachReady,
	//	&attachOptions,
	//)
	//if err != nil {
	//	fmt.Println(err)
	//	os.Exit(1)
	//}

	//_, err = containers.Wait(conn, container.ID, &waitOptions)
}

func CreateSpec(containerConfig config.ContainerConfig) *specgen.SpecGenerator {
	var imageName string = fmt.Sprintf("localhost/sandman/%s", containerConfig.Name)
	spec := specgen.NewSpecGenerator(imageName, false)

	// Default configuration
	spec.Terminal = true
	spec.Stdin = true
	spec.Remove = true

	// TODO: Expand container config

	return spec
}

func CmdExecute(socket string, verbose bool, keep bool, args []string) {
	var container_name string = args[0]
	var runCmd []string = args[1:]
	Start(socket, config.LoadConfig(container_name), keep, verbose, runCmd)
}
