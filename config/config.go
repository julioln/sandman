package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/user"

	"github.com/BurntSushi/toml"
)

const SANDMAN_DIR = ".config/sandman"

type ContainerConfigBuild struct {
	Instructions string
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
	Network    string
	Name       string
	Persistent []string
	Volumes    []string
	Devices    []string
	Args       []string
	Env        []string
	Ports      []string
}

type ContainerConfig struct {
	Build ContainerConfigBuild
	Run   ContainerConfigRun
	Name  string
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

	return config
}
