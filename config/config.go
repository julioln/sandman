package config

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/user"

	"github.com/julioln/sandman/constants"

	"github.com/BurntSushi/toml"
	nettypes "github.com/containers/common/libnetwork/types"
	specs "github.com/opencontainers/runtime-spec/specs-go"
)

type ContainerConfigBuild struct {
	Instructions         string
	ContextDirectory     string
	AdditionalImageNames []string
	Limits               ContainerConfigBuildLimits
}

type ContainerConfigBuildLimits struct {
	Ulimit []string
}

type ContainerConfigRun struct {
	X11        bool
	Wayland    bool
	Dri        bool
	Ipc        bool
	Gpu        bool
	Pulseaudio bool
	Pipewire   bool
	Dbus       bool
	Net        bool
	Uidmap     bool
	Home       bool
	Fonts      bool
	Network    string
	Name       string
	Volumes    []string
	Env        []string
	Devices    []string
	Ports      []string
	UsbDevices []string
	RawMounts  []specs.Mount
	RawPorts   []nettypes.PortMapping
	RawDevices []specs.LinuxDevice
	Limits     ContainerConfigRunLimits
}

type ContainerConfigRunLimits struct {
	CPU        specs.LinuxCPU
	Memory     specs.LinuxMemory
	Pids       specs.LinuxPids
	Rlimits    []specs.POSIXRlimit
	CgroupConf map[string]string
}

type ContainerConfig struct {
	Name       string
	ImageName  string
	ConfigFile string
	Build      ContainerConfigBuild
	Run        ContainerConfigRun
}

func getHomeDir() string {
	homedir, err := os.UserHomeDir()
	if err != nil {
		user, err := user.Current()

		if err != nil {
			os.Exit(1)
		}

		homedir = fmt.Sprintf("/home/%s", user.Username)
	}

	return homedir
}

func GetHomeStorageDir() string {
	return fmt.Sprintf("%s/%s", getHomeDir(), constants.SANDMAN_LOCAL_STORAGE)
}

func GetSandmanConfigDir() string {
	return fmt.Sprintf("%s/%s", getHomeDir(), constants.SANDMAN_DIR)
}

func LoadConfig(container_name string) ContainerConfig {
	var config ContainerConfig
	var config_file_content []byte
	var config_file_path string = fmt.Sprintf("%s/%s.toml", GetSandmanConfigDir(), container_name)

	config_file_content, err := ioutil.ReadFile(config_file_path)

	if err != nil {
		fmt.Printf("Can't read container configuration file at %s", config_file_path)
		fmt.Println("Error: ", err)
		os.Exit(1)
	}

	_, err = toml.Decode(string(config_file_content), &config)

	if err != nil {
		fmt.Printf("Can't decode container configuration file at %s", config_file_path)
		fmt.Println("Error: ", err)
		os.Exit(1)
	}

	config.Name = container_name
	config.ConfigFile = config_file_path
	config.ImageName = fmt.Sprintf("sandman/%s", container_name)

	return config
}

func Scaffold() string {
	buf := new(bytes.Buffer)
	err := toml.NewEncoder(buf).Encode(map[string]interface{}{
		"Build": new(ContainerConfigBuild),
		"Run":   new(ContainerConfigRun),
	})

	if err != nil {
		fmt.Printf("Failed to encode TOML")
		fmt.Println("Error: ", err)
		os.Exit(1)
	}

	return buf.String()
}
