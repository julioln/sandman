package run

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/julioln/sandman/config"
	"github.com/julioln/sandman/constants"
	"github.com/julioln/sandman/podman"

	"github.com/containers/podman/v4/libpod/define"
	"github.com/containers/podman/v4/pkg/bindings/containers"
	"github.com/containers/podman/v4/pkg/specgen"
)

var (
	configFunctions = []func(spec *specgen.SpecGenerator, config config.ContainerConfig){
		Dbus,
		Devices,
		Env,
		Fonts,
		Gpu,
		Home,
		Ipc,
		Limits,
		Name,
		Network,
		Pipewire,
		Ports,
		Pulseaudio,
		Raw,
		Uidmap,
		Usb,
		Volumes,
		Wayland,
		X11,
	}
)

func Start(socket string, containerConfig config.ContainerConfig, attach bool, keep bool, verbose bool, runCmd []string) {
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
		fmt.Printf("Container Config: %#v\n", containerConfig)
		fmt.Printf("Connection: %#v\n", conn)
		fmt.Printf("Container Spec: %#v\n", spec)
	}

	var createOptions containers.CreateOptions
	container, err := containers.CreateWithSpec(conn, spec, &createOptions)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if verbose {
		fmt.Printf("Container: %#v\n", container)
	}

	if err = containers.Start(conn, container.ID, nil); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	var waitOptions containers.WaitOptions
	waitOptions.Condition = append(waitOptions.Condition, define.ContainerStateRunning)
	if _, err = containers.Wait(conn, container.ID, &waitOptions); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if verbose {
		var inspectOptions containers.InspectOptions
		if containerData, err := containers.Inspect(conn, container.ID, &inspectOptions); err == nil {
			fmt.Printf("Container data: %#v\n", containerData)
		}
	}

	if attach {
		attachOptions := new(containers.AttachOptions)
		if err = containers.Attach(conn, container.ID, os.Stdin, os.Stdout, os.Stderr, nil, attachOptions); err != nil {
			fmt.Println("Failed to attach to container: ", err)
		}
	}
}

func CreateSpec(containerConfig config.ContainerConfig) *specgen.SpecGenerator {
	spec := specgen.NewSpecGenerator(containerConfig.ImageName, false)

	// Default configuration
	spec.Terminal = true
	spec.Stdin = true
	spec.Remove = true
	spec.Hostname = strings.Replace(containerConfig.ImageName, "/", "_", -1)
	spec.Umask = "0022"
	spec.Env = make(map[string]string)
	spec.Labels = make(map[string]string)
	spec.Labels["sandman_container_name"] = containerConfig.Name
	spec.Labels["sandman_image_name"] = containerConfig.ImageName
	spec.Labels["sandman_version"] = constants.VERSION

	// Apply all configurators
	for _, f := range configFunctions {
		f(spec, containerConfig)
	}

	return spec
}

func CmdExecuteStart(socket string, verbose bool, keep bool, args []string) {
	var container_name string = args[0]
	var runCmd []string = args[1:]
	Start(socket, config.LoadConfig(container_name), false, keep, verbose, runCmd)
}

func CmdExecuteRun(socket string, verbose bool, keep bool, args []string) {
	var container_name string = args[0]
	var runCmd []string = args[1:]
	Start(socket, config.LoadConfig(container_name), true, keep, verbose, runCmd)
}
