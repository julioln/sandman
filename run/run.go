package run

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/julioln/sandman/config"
	"github.com/julioln/sandman/podman"

	"github.com/containers/podman/v4/libpod/define"
	"github.com/containers/podman/v4/pkg/bindings/containers"
	"github.com/containers/podman/v4/pkg/specgen"
	"github.com/containers/storage/pkg/idtools"
	"github.com/containers/storage/types"

	specs "github.com/opencontainers/runtime-spec/specs-go"
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

	err = containers.Start(conn, container.ID, nil)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	var waitOptions containers.WaitOptions
	waitOptions.Condition = append(waitOptions.Condition, define.ContainerStateRunning)
	_, err = containers.Wait(conn, container.ID, &waitOptions)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if verbose {
		var inspectOptions containers.InspectOptions
		containerData, err := containers.Inspect(conn, container.ID, &inspectOptions)
		if err == nil {
			fmt.Printf("Container data: %#v\n", containerData)
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
}

func CreateSpec(containerConfig config.ContainerConfig) *specgen.SpecGenerator {
	spec := specgen.NewSpecGenerator(containerConfig.ImageName, false)

	// Default configuration
	spec.Terminal = true
	spec.Stdin = true
	spec.Remove = true
	spec.Hostname = strings.Replace(containerConfig.Name, "/", "-", -1)
	spec.Umask = "0022"
	spec.Env = make(map[string]string)
	spec.Labels = make(map[string]string)
	spec.Labels["sandman_container_name"] = containerConfig.Name
	spec.Labels["sandman_image_name"] = containerConfig.ImageName

	if containerConfig.Run.X11 {
		spec.Env["DISPLAY"] = os.Getenv("DISPLAY")
		spec.Env["XCURSOR_THEME"] = os.Getenv("XCURSOR_THEME")
		spec.Env["XCURSOR_SIZE"] = os.Getenv("XCURSOR_SIZE")
		spec.Mounts = append(spec.Mounts, specs.Mount{
			Destination: "/tmp/.X11-unix",
			Source:      "/tmp/.X11-unix",
			Type:        "bind",
		})
	}

	if containerConfig.Run.Wayland {
		spec.Env["WAYLAND_DISPLAY"] = os.Getenv("WAYLAND_DISPLAY")
		spec.Mounts = append(spec.Mounts, specs.Mount{
			Destination: fmt.Sprintf("%s/%s", os.Getenv("XDG_RUNTIME_DIR"), os.Getenv("WAYLAND_DISPLAY")),
			Source:      fmt.Sprintf("%s/%s", os.Getenv("XDG_RUNTIME_DIR"), os.Getenv("WAYLAND_DISPLAY")),
			Type:        "bind",
		})
	}

	if containerConfig.Run.Dri || containerConfig.Run.Gpu {
		spec.Devices = append(spec.Devices, specs.LinuxDevice{
			Path: "/dev/dri",
		})
	}

	if containerConfig.Run.Ipc {
		spec.IpcNS.NSMode = specgen.NamespaceMode("host")
	}

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

	if containerConfig.Run.Dbus {
		spec.Env["DBUS_SESSION_BUS_ADDRESS"] = fmt.Sprintf("unix:path=%s/bus", os.Getenv("XDG_RUNTIME_DIR"))
		spec.Mounts = append(spec.Mounts, specs.Mount{
			Destination: fmt.Sprintf("%s/bus", os.Getenv("XDG_RUNTIME_DIR")),
			Source:      fmt.Sprintf("%s/bus", os.Getenv("XDG_RUNTIME_DIR")),
			Type:        "bind",
		})
	}

	if containerConfig.Run.Uidmap {
		var idMappingOptions types.IDMappingOptions
		idMappingOptions.HostUIDMapping = true
		idMappingOptions.HostGIDMapping = true
		idMappingOptions.UIDMap = append(idMappingOptions.UIDMap,
			idtools.IDMap{
				ContainerID: os.Getuid(),
				HostID:      0,
				Size:        1,
			},
			idtools.IDMap{
				ContainerID: 0,
				HostID:      1,
				Size:        os.Getuid(),
			},
			idtools.IDMap{
				ContainerID: os.Getuid() + 1,
				HostID:      os.Getuid() + 1,
				Size:        65536 - os.Getuid(),
			},
		)
		idMappingOptions.GIDMap = append(idMappingOptions.GIDMap,
			idtools.IDMap{
				ContainerID: os.Getuid(),
				HostID:      0,
				Size:        1,
			},
			idtools.IDMap{
				ContainerID: 0,
				HostID:      1,
				Size:        os.Getuid(),
			},
			idtools.IDMap{
				ContainerID: os.Getuid() + 1,
				HostID:      os.Getuid() + 1,
				Size:        65536 - os.Getuid(),
			},
		)
		spec.IDMappings = &idMappingOptions
		spec.UserNS.NSMode = specgen.Private
		spec.User = fmt.Sprint(os.Getuid())
	}

	if containerConfig.Run.Name != "" {
		spec.Name = containerConfig.Run.Name
		spec.Hostname = containerConfig.Run.Name
	}

	return spec
}

func CmdExecute(socket string, verbose bool, keep bool, args []string) {
	var container_name string = args[0]
	var runCmd []string = args[1:]
	Start(socket, config.LoadConfig(container_name), keep, verbose, runCmd)
}
