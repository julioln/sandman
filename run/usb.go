package run

import (
	"fmt"
	"strings"

	"github.com/citilinkru/libudev"
	"github.com/citilinkru/libudev/matcher"
	"github.com/citilinkru/libudev/types"
	"github.com/containers/podman/v4/pkg/specgen"
	"github.com/julioln/sandman/config"
	specs "github.com/opencontainers/runtime-spec/specs-go"
)

func extractDevicePaths(dev *types.Device) []string {
	var paths []string

	paths = append(paths, fmt.Sprintf("/dev/%s", dev.Env["DEVNAME"]))

	if len(dev.Children) > 0 {
		for _, child := range dev.Children {
			// Recurse
			paths = append(paths, extractDevicePaths(child)...)
		}
	}

	return paths
}

func UsbDevicePaths(idVendor string, idProduct string) []string {
	var paths []string

	devices, err := UsbDevices(idVendor, idProduct)
	if err != nil {
		fmt.Println("Failed to list usb devices, ignoring.")
		return []string{}
	}

	for _, d := range devices {
		paths = append(paths, extractDevicePaths(d)...)
	}

	return paths
}

func UsbDevices(idVendor string, idProduct string) ([]*types.Device, error) {
	sc := libudev.NewScanner()
	err, devices := sc.ScanDevices()

	if err != nil {
		return nil, err
	}

	m := matcher.NewMatcher()
	m.SetStrategy(matcher.StrategyAnd)
	m.AddRule(matcher.NewRuleEnv("DEVNAME", "usb"))

	if idVendor != "" {
		m.AddRule(matcher.NewRuleAttr("idVendor", idVendor))
	}
	if idProduct != "" {
		m.AddRule(matcher.NewRuleAttr("idProduct", idProduct))
	}

	return m.Match(devices), nil
}

func Usb(spec *specgen.SpecGenerator, containerConfig config.ContainerConfig) {
	// Add usb devices
	for _, usbDev := range containerConfig.Run.UsbDevices {
		u := strings.Split(usbDev, ":")
		var vendor string
		var product string
		if len(u) < 2 {
			vendor = u[0]
			// Catch all
			product = ""
		} else {
			vendor = u[0]
			product = u[1]
		}

		for _, devicePath := range UsbDevicePaths(vendor, product) {
			spec.Devices = append(spec.Devices, specs.LinuxDevice{
				Path: devicePath,
			})
		}
	}
}
