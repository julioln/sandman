package run

import (
	"context"
	"fmt"
	"os"
	"strings"
	"strconv"

	"github.com/julioln/sandman/config"
	"github.com/julioln/sandman/constants"
	"github.com/julioln/sandman/podman"

	"github.com/containers/podman/v4/libpod/define"
	"github.com/containers/podman/v4/pkg/bindings/containers"
	"github.com/containers/podman/v4/pkg/specgen"
	"github.com/containers/storage/pkg/idtools"
	"github.com/containers/storage/types"

	nettypes "github.com/containers/common/libnetwork/types"
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
	spec.Hostname = strings.Replace(containerConfig.ImageName, "/", "_", -1)
	spec.Umask = "0022"
	spec.Env = make(map[string]string)
	spec.Labels = make(map[string]string)
	spec.Labels["sandman_container_name"] = containerConfig.Name
	spec.Labels["sandman_image_name"] = containerConfig.ImageName
	spec.Labels["sandman_version"] = constants.VERSION

	// X11 Forwarding
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

	// Wayland forwarding
	if containerConfig.Run.Wayland {
		spec.Env["WAYLAND_DISPLAY"] = os.Getenv("WAYLAND_DISPLAY")
		spec.Mounts = append(spec.Mounts, specs.Mount{
			Destination: fmt.Sprintf("%s/%s", os.Getenv("XDG_RUNTIME_DIR"), os.Getenv("WAYLAND_DISPLAY")),
			Source:      fmt.Sprintf("%s/%s", os.Getenv("XDG_RUNTIME_DIR"), os.Getenv("WAYLAND_DISPLAY")),
			Type:        "bind",
		})
	}

	// Expose Direct Render Inteface for GPU acceleration
	if containerConfig.Run.Dri || containerConfig.Run.Gpu {
		spec.Devices = append(spec.Devices, specs.LinuxDevice{
			Path: "/dev/dri",
		})
	}

	// Share IPC with host
	if containerConfig.Run.Ipc {
		spec.IpcNS.NSMode = specgen.NamespaceMode("host")
	}

	// Expose pulseaudio client
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

	// Pipewire client
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

	// Expose host dbus
	if containerConfig.Run.Dbus {
		spec.Env["DBUS_SESSION_BUS_ADDRESS"] = fmt.Sprintf("unix:path=%s/bus", os.Getenv("XDG_RUNTIME_DIR"))
		spec.Mounts = append(spec.Mounts, specs.Mount{
			Destination: fmt.Sprintf("%s/bus", os.Getenv("XDG_RUNTIME_DIR")),
			Source:      fmt.Sprintf("%s/bus", os.Getenv("XDG_RUNTIME_DIR")),
			Type:        "bind",
		})
	}

	// Add host fonts
	if containerConfig.Run.Fonts {
		spec.Mounts = append(spec.Mounts, specs.Mount{
			Destination: "/usr/share/fonts",
			Source:      "/usr/share/fonts",
			Type:        "bind",
			Options:     []string{"ro"},
		})
	}

	// Create custom user namespace inside container with a mapped uid
	if containerConfig.Run.Uidmap {
		idMaps := []idtools.IDMap{
			{
				ContainerID: os.Getuid(),
				HostID:      0,
				Size:        1,
			},
			{
				ContainerID: 0,
				HostID:      1,
				Size:        os.Getuid(),
			},
			{
				ContainerID: os.Getuid() + 1,
				HostID:      os.Getuid() + 1,
				Size:        65536 - os.Getuid(),
			},
		}
		var idMappingOptions types.IDMappingOptions
		idMappingOptions.HostUIDMapping = true
		idMappingOptions.HostGIDMapping = true
		idMappingOptions.UIDMap = append(idMappingOptions.UIDMap, idMaps...)
		idMappingOptions.GIDMap = append(idMappingOptions.GIDMap, idMaps...)
		spec.IDMappings = &idMappingOptions
		spec.UserNS.NSMode = specgen.Private
		spec.User = fmt.Sprint(os.Getuid())
	}

	// Customize container name
	if containerConfig.Run.Name != "" {
		spec.Name = containerConfig.Run.Name
		spec.Hostname = containerConfig.Run.Name
	}

	// Configure network namespace
	var networkNS specgen.Namespace
	if containerConfig.Run.Network == "" && !containerConfig.Run.Net {
		networkNS.NSMode = specgen.None
	} else if containerConfig.Run.Net {
		// Backwards compatibility
		networkNS.NSMode = specgen.Slirp
	} else if networkNS, _, _, err := specgen.ParseNetworkFlag([]string{containerConfig.Run.Network}); err != nil {
		fmt.Println("Error parsing network, defaulting to none: ", err)
		networkNS.NSMode = specgen.None
	}
	spec.NetNS = networkNS

	// Setup port mappings
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
			HostPort: uint16(hostPort),
		})
	}

	// Automatically mount home in a persistent location
	if containerConfig.Run.Home {
		var mountPoint = fmt.Sprintf("%s/%s", config.GetHomeStorageDir(), containerConfig.Name)
		if err := os.MkdirAll(mountPoint, 0755); err == nil {
			spec.Mounts = append(spec.Mounts, specs.Mount{
				Destination: "/home",
				Source:      mountPoint,
				Type:        "bind",
			})
		}
	}

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

	// Export additional environments
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

	// Mount additional linux devices
	for _, dev := range containerConfig.Run.Devices {
		spec.Devices = append(spec.Devices, specs.LinuxDevice{
			Path: dev,
		})
	}

	// Add raw configuration
	spec.PortMappings = append(spec.PortMappings, containerConfig.Run.RawPorts...)
	spec.Mounts = append(spec.Mounts, containerConfig.Run.RawMounts...)
	spec.Devices = append(spec.Devices, containerConfig.Run.RawDevices...)

	// Apply limits
	spec.ResourceLimits = &containerConfig.Limits.LinuxResources

	return spec
}

func CmdExecute(socket string, verbose bool, keep bool, args []string) {
	var container_name string = args[0]
	var runCmd []string = args[1:]
	Start(socket, config.LoadConfig(container_name), keep, verbose, runCmd)
}
