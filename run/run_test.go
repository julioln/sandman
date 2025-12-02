package run

import (
	"fmt"
	"os"
	"reflect"
	"testing"

	"github.com/containers/podman/v6/pkg/specgen"
	"github.com/julioln/sandman/config"
	"github.com/julioln/sandman/constants"
	specs "github.com/opencontainers/runtime-spec/specs-go"
	"go.podman.io/common/libnetwork/types"
)

func compareMountPoints(m1 specs.Mount, m2 specs.Mount) bool {
	return reflect.DeepEqual(m1, m2)
}

func compareDevices(d1 specs.LinuxDevice, d2 specs.LinuxDevice) bool {
	return reflect.DeepEqual(d1, d2)
}

func comparePorts(p1 types.PortMapping, p2 types.PortMapping) bool {
	return reflect.DeepEqual(p1, p2)
}

func testMaps(t *testing.T, expected map[string]string, existing map[string]string) {
	for k, l1 := range expected {
		l2, exists := existing[k]
		if !exists || l1 != l2 {
			t.Errorf("expected value but couldn't find it: %s = %s", k, l1)
		}
	}
}

func testMountPoints(t *testing.T, spec *specgen.SpecGenerator, mounts []specs.Mount) {
	for _, expected_mount := range mounts {
		found := false
		for _, existing_mount := range spec.Mounts {
			if compareMountPoints(expected_mount, existing_mount) {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("expected mountpoint but couldn't find it: %#v", expected_mount)
		}
	}
}

func testDevices(t *testing.T, spec *specgen.SpecGenerator, devices []specs.LinuxDevice) {
	for _, expected_device := range devices {
		found := false
		for _, existing_device := range spec.Devices {
			if compareDevices(expected_device, existing_device) {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("expected device but couldn't find it: %#v", expected_device)
		}
	}
}

func testPorts(t *testing.T, spec *specgen.SpecGenerator, ports []types.PortMapping) {
	for _, expected_port := range ports {
		found := false
		for _, existing_port := range spec.PortMappings {
			if comparePorts(expected_port, existing_port) {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("expected port but couldn't find it: %#v", expected_port)
		}
	}
}

func TestCreateSpec(t *testing.T) {
	testConfig := new(config.ContainerConfig)
	testConfig.Name = "name"
	testConfig.ImageName = "image/name"
	spec := CreateSpec(*testConfig)

	if spec.Hostname != "image_name" {
		t.Errorf("container hostname incorrect")
	}

	labels := map[string]string{
		"sandman_container_name": "name",
		"sandman_image_name":     "image/name",
		"sandman_version":        constants.VERSION,
	}
	testMaps(t, labels, spec.Labels)

	if !*spec.Terminal {
		t.Errorf("terminal incorrect, expected true, got false")
	}
	if !*spec.Stdin {
		t.Errorf("stdin incorrect, expected true, got false")
	}
	if !*spec.Remove {
		t.Errorf("remove incorrect, expected true, got false")
	}
	if spec.Umask != "0022" {
		t.Errorf("umask incorrect, expected %s, got %s", "0022", spec.Umask)
	}
}

func TestDbus(t *testing.T) {
	testConfig := new(config.ContainerConfig)
	testConfig.Run.Dbus = true
	spec := CreateSpec(*testConfig)

	vars := map[string]string{
		"DBUS_SESSION_BUS_ADDRESS": fmt.Sprintf("unix:path=%s/bus", os.Getenv("XDG_RUNTIME_DIR")),
	}
	testMaps(t, spec.Env, vars)

	mountPoints := []specs.Mount{
		{
			Destination: fmt.Sprintf("%s/bus", os.Getenv("XDG_RUNTIME_DIR")),
			Source:      fmt.Sprintf("%s/bus", os.Getenv("XDG_RUNTIME_DIR")),
			Type:        "bind",
		},
	}
	testMountPoints(t, spec, mountPoints)
}

func TestDevices(t *testing.T) {
	testConfig := new(config.ContainerConfig)
	testConfig.Run.Devices = []string{"/dev/test0", "/dev/test1"}
	spec := CreateSpec(*testConfig)

	devices := []specs.LinuxDevice{
		{
			Path: "/dev/test0",
		},
		{
			Path: "/dev/test1",
		},
	}
	testDevices(t, spec, devices)
}

func TestEnv(t *testing.T) {
	testConfig := new(config.ContainerConfig)
	testConfig.Run.Env = []string{"PWD", "TEST1=value1"}
	spec := CreateSpec(*testConfig)

	vars := map[string]string{
		"PWD":   os.Getenv("PWD"),
		"TEST1": "value1",
	}
	testMaps(t, spec.Env, vars)
}

func TestFonts(t *testing.T) {
	testConfig := new(config.ContainerConfig)
	testConfig.Run.Fonts = true
	spec := CreateSpec(*testConfig)

	mounts := []specs.Mount{
		{
			Destination: "/usr/share/fonts",
			Source:      "/usr/share/fonts",
			Type:        "bind",
			Options:     []string{"ro"},
		},
	}
	testMountPoints(t, spec, mounts)
}

func TestGpu(t *testing.T) {
	devices := []specs.LinuxDevice{
		{
			Path: "/dev/dri",
		},
	}

	testConfig := new(config.ContainerConfig)
	testConfig.Run.Dri = true
	spec := CreateSpec(*testConfig)
	testDevices(t, spec, devices)

	testConfig = new(config.ContainerConfig)
	testConfig.Run.Gpu = true
	spec = CreateSpec(*testConfig)
	testDevices(t, spec, devices)
}

func TestIpc(t *testing.T) {
	testConfig := new(config.ContainerConfig)
	testConfig.Run.Ipc = true
	spec := CreateSpec(*testConfig)
	nsMode := specgen.NamespaceMode("host")

	if spec.IpcNS.NSMode != nsMode {
		t.Errorf("Ipc incorrect, expected %s, got %s", nsMode, spec.IpcNS.NSMode)
	}
}

func TestName(t *testing.T) {
	name := "testing_name_override"
	testConfig := new(config.ContainerConfig)
	testConfig.Name = "original_name"
	testConfig.Run.Name = name
	spec := CreateSpec(*testConfig)

	if spec.Name != name {
		t.Errorf("Name incorrect, expected %s, got %s", spec.Name, name)
	}
	if spec.Hostname != name {
		t.Errorf("Hostname incorrect, expected %s, got %s", spec.Hostname, name)
	}
}

func TestPipewire(t *testing.T) {
	testConfig := new(config.ContainerConfig)
	testConfig.Run.Pipewire = true
	spec := CreateSpec(*testConfig)

	vars := map[string]string{
		"XDG_RUNTIME_DIR": os.Getenv("XDG_RUNTIME_DIR"),
	}
	testMaps(t, spec.Env, vars)

	mountPoints := []specs.Mount{
		{
			Destination: fmt.Sprintf("%s/pipewire-0", os.Getenv("XDG_RUNTIME_DIR")),
			Source:      fmt.Sprintf("%s/pipewire-0", os.Getenv("XDG_RUNTIME_DIR")),
			Type:        "bind",
		},
	}
	testMountPoints(t, spec, mountPoints)
}

func TestPorts(t *testing.T) {
	testConfig := new(config.ContainerConfig)
	testConfig.Run.Ports = []string{"3000:4000", "invalid"}
	spec := CreateSpec(*testConfig)
	ports := []types.PortMapping{
		{
			ContainerPort: 3000,
			HostPort:      4000,
		},
	}
	testPorts(t, spec, ports)
}

func TestPulseaudio(t *testing.T) {
	testConfig := new(config.ContainerConfig)
	testConfig.Run.Pulseaudio = true
	spec := CreateSpec(*testConfig)

	vars := map[string]string{
		"XDG_RUNTIME_DIR": os.Getenv("XDG_RUNTIME_DIR"),
	}
	testMaps(t, spec.Env, vars)

	mountPoints := []specs.Mount{
		{
			Destination: "/etc/machine-id",
			Source:      "/etc/machine-id",
			Type:        "bind",
			Options:     []string{"ro"},
		},
		{
			Destination: fmt.Sprintf("%s/pulse/native", os.Getenv("XDG_RUNTIME_DIR")),
			Source:      fmt.Sprintf("%s/pulse/native", os.Getenv("XDG_RUNTIME_DIR")),
			Type:        "bind",
		},
	}
	testMountPoints(t, spec, mountPoints)
}

func TestNetwork(t *testing.T) {
	testConfig := new(config.ContainerConfig)
	var ns specgen.Namespace

	testConfig.Run.Net = true
	testConfig.Run.Network = ""
	spec := CreateSpec(*testConfig)
	ns.NSMode = specgen.Slirp
	if spec.NetNS.NSMode != ns.NSMode {
		t.Errorf("Network namespace incorrect, expected %#v, got %#v", ns.NSMode, spec.NetNS.NSMode)
	}

	testConfig.Run.Net = false
	testConfig.Run.Network = ""
	spec = CreateSpec(*testConfig)
	ns.NSMode = specgen.None
	if spec.NetNS.NSMode != ns.NSMode {
		t.Errorf("Network namespace incorrect, expected %#v, got %#v", ns.NSMode, spec.NetNS.NSMode)
	}

	testConfig.Run.Net = false
	testConfig.Run.Network = "host"
	spec = CreateSpec(*testConfig)
	ns, _, _, _ = specgen.ParseNetworkFlag([]string{"host"})
	if spec.NetNS.NSMode != ns.NSMode {
		t.Errorf("Network namespace incorrect, expected %#v, got %#v", ns.NSMode, spec.NetNS.NSMode)
	}

	testConfig.Run.Net = false
	testConfig.Run.Network = "other"
	spec = CreateSpec(*testConfig)
	ns.NSMode = specgen.Bridge
	if spec.NetNS.NSMode != ns.NSMode {
		t.Errorf("Network namespace incorrect, expected %#v, got %#v", ns.NSMode, spec.NetNS.NSMode)
	}
}

func TestUidmap(t *testing.T) {
	testConfig := new(config.ContainerConfig)
	testConfig.Run.Uidmap = true
	spec := CreateSpec(*testConfig)

	if spec.UserNS.NSMode != specgen.Private {
		t.Errorf("User namespace incorrect, expected %#v, got %#v", specgen.Private, spec.UserNS.NSMode)
	}
	if !spec.IDMappings.HostGIDMapping {
		t.Errorf("Host GID mapping incorrect, expected %#v, got %#v", true, spec.IDMappings.HostGIDMapping)
	}
	if !spec.IDMappings.HostUIDMapping {
		t.Errorf("Host UID mapping incorrect, expected %#v, got %#v", true, spec.IDMappings.HostUIDMapping)
	}
	if spec.User != fmt.Sprint(os.Getuid()) {
		t.Errorf("User incorrect, expected %#v, got %#v", os.Getuid(), spec.User)
	}
}

func TestVolumes(t *testing.T) {
	testConfig := new(config.ContainerConfig)
	testConfig.Run.Volumes = []string{"/vol1", "/vol2:/vol3", "/vol4:/vol5:ro,atime"}
	spec := CreateSpec(*testConfig)

	mountPoints := []specs.Mount{
		{
			Destination: "/vol1",
			Source:      "/vol1",
			Type:        "bind",
		},
		{
			Destination: "/vol3",
			Source:      "/vol2",
			Type:        "bind",
		},
		{
			Destination: "/vol5",
			Source:      "/vol4",
			Type:        "bind",
			Options:     []string{"ro", "atime"},
		},
	}
	testMountPoints(t, spec, mountPoints)
}

func TestWayland(t *testing.T) {
	testConfig := new(config.ContainerConfig)
	testConfig.Run.Wayland = true
	spec := CreateSpec(*testConfig)

	vars := map[string]string{
		"XDG_SESSION_TYPE": os.Getenv("XDG_SESSION_TYPE"),
		"WAYLAND_DISPLAY":  os.Getenv("WAYLAND_DISPLAY"),
	}
	testMaps(t, spec.Env, vars)

	mountPoints := []specs.Mount{
		{
			Destination: fmt.Sprintf("%s/%s", os.Getenv("XDG_RUNTIME_DIR"), os.Getenv("WAYLAND_DISPLAY")),
			Source:      fmt.Sprintf("%s/%s", os.Getenv("XDG_RUNTIME_DIR"), os.Getenv("WAYLAND_DISPLAY")),
			Type:        "bind",
		},
	}
	testMountPoints(t, spec, mountPoints)
}

func TestX11(t *testing.T) {
	testConfig := new(config.ContainerConfig)
	testConfig.Run.X11 = true
	spec := CreateSpec(*testConfig)

	vars := map[string]string{
		"XDG_SESSION_TYPE": os.Getenv("XDG_SESSION_TYPE"),
		"DISPLAY":          os.Getenv("DISPLAY"),
		"XCURSOR_THEME":    os.Getenv("XCURSOR_THEME"),
		"XCURSOR_SIZE":     os.Getenv("XCURSOR_SIZE"),
	}
	testMaps(t, spec.Env, vars)

	mountPoints := []specs.Mount{
		{
			Destination: "/tmp/.X11-unix",
			Source:      "/tmp/.X11-unix",
			Type:        "bind",
		},
	}
	testMountPoints(t, spec, mountPoints)
}
