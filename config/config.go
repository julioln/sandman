package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/user"

	"github.com/BurntSushi/toml"
	nettypes "github.com/containers/common/libnetwork/types"
	specs "github.com/opencontainers/runtime-spec/specs-go"
)

const SANDMAN_DIR = ".config/sandman"
const SANDMAN_LOCAL_STORAGE = ".local/share/sandman"

type ContainerConfigBuild struct {
	Instructions     string
	ContextDirectory string
}

type ContainerConfigRun struct {
	X11        bool
	Wayland    bool
	Dri        bool
	Ipc        bool
	Gpu        bool
	Pulseaudio bool
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
	Ports      []nettypes.PortMapping
	Mounts     []specs.Mount
}

type ContainerConfig struct {
	Name       string
	ImageName  string
	ConfigFile string
	Build      ContainerConfigBuild
	Run        ContainerConfigRun
	Limits     specs.LinuxResources
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
	return fmt.Sprintf("%s/%s", getHomeDir(), SANDMAN_LOCAL_STORAGE)
}

func LoadConfig(container_name string) ContainerConfig {
	var config ContainerConfig
	var config_file_content []byte

	var homedir string = getHomeDir()
	var config_file_path string = fmt.Sprintf("%s/%s/%s.toml", homedir, SANDMAN_DIR, container_name)

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
