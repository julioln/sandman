package config

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/user"

	"github.com/julioln/sandman/constants"

	"dario.cat/mergo"
	"github.com/BurntSushi/toml"
	nettypes "github.com/containers/common/libnetwork/types"
	"github.com/containers/storage/pkg/archive"
	specs "github.com/opencontainers/runtime-spec/specs-go"
)

type ContainerConfigBuild struct {
	Instructions         string
	ContextDirectory     string
	Compression          archive.Compression
	AdditionalImageNames []string
	Limits               ContainerConfigBuildLimits
}

type ContainerConfigBuildLimits struct {
	Ulimit []string
}

type ContainerConfigRun struct {
	X11          bool
	Wayland      bool
	Dri          bool
	Ipc          bool
	Gpu          bool
	Pulseaudio   bool
	Pipewire     bool
	Dbus         bool
	Net          bool
	Uidmap       bool
	Home         bool
	HomePath     string
	Fonts        bool
	Network      string
	Name         string
	CgroupParent string
	Priviledged  bool
	Volumes      []string
	Env          []string
	Devices      []string
	Ports        []string
	UsbDevices   []string
	RawMounts    []specs.Mount
	RawPorts     []nettypes.PortMapping
	RawDevices   []specs.LinuxDevice
	Limits       ContainerConfigRunLimits
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

type SandmanConfig struct {
	Defaults SandmanConfigDefaults
}

type SandmanConfigDefaults struct {
	Build ContainerConfigBuild
	Run   ContainerConfigRun
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

func GetOldSandmanConfigDir() string {
	return fmt.Sprintf("%s/%s", getHomeDir(), constants.OLD_SANDMAN_DIR)
}

func GetSandmanConfigDir() string {
	return fmt.Sprintf("%s/%s", getHomeDir(), constants.SANDMAN_DIR)
}

func GetSandmanConfigFilename() string {
	return fmt.Sprintf("%s/%s", getHomeDir(), constants.SANDMAN_CONF)
}

func Setup() error {
	if err := os.MkdirAll(GetHomeStorageDir(), 0755); err != nil {
		fmt.Println("Error: ", err)
		return err
	}
	if err := os.MkdirAll(GetSandmanConfigDir(), 0755); err != nil {
		fmt.Println("Error: ", err)
		return err
	}

	_, err := os.Stat(GetSandmanConfigFilename())
	if os.IsNotExist(err) {
		file, err := os.Create(GetSandmanConfigFilename())
		if err != nil {
			fmt.Println("Error: ", err)
			return err
		}
		defer file.Close()
	}

	return nil
}

func CheckConfig() {
	stat, err := os.Stat(GetHomeStorageDir())
	if err != nil {
		fmt.Println("Storage Directory does not exist: ", GetHomeStorageDir())
		fmt.Println("  -> Run `sandman setup` to create it")
	} else if !stat.IsDir() {
		fmt.Println("Storage Directory exists but is not a dir: ", GetHomeStorageDir())
		fmt.Println("  -> Delete it and run `sandman setup` to create it as a directory")
	} else {
		fmt.Println("Storage Directory exists: ", GetHomeStorageDir())
	}

	_, err = os.Stat(GetOldSandmanConfigDir())
	if err == nil {
		fmt.Println("Legacy Sandman Configuration Directory exists: ", GetOldSandmanConfigDir())
		fmt.Println("  -> Rename it or link it to: ", GetSandmanConfigDir())
	}

	stat, err = os.Stat(GetSandmanConfigDir())
	if err != nil {
		fmt.Println("Container Configuration Directory does not exist: ", GetSandmanConfigDir())
		fmt.Println("  -> Run `sandman setup` to create it")
	} else if !stat.IsDir() {
		fmt.Println("Container Configuration Directory exists but is not a dir: ", GetSandmanConfigDir())
		fmt.Println("  -> Delete it and run `sandman setup` to create it as a directory")
	} else {
		fmt.Println("Container Configuration Directory exists: ", GetSandmanConfigDir())
	}

	stat, err = os.Stat(GetSandmanConfigFilename())
	if err != nil {
		fmt.Println("Sandman Configuration File does not exist: ", GetSandmanConfigFilename())
		fmt.Println("  -> Run `sandman setup` to create it")
	} else if stat.IsDir() {
		fmt.Println("Sandman Configuration File exists but is a dir: ", GetSandmanConfigFilename())
		fmt.Println("  -> Delete it and run `sandman setup` to create it as a regular file")
	} else {
		fmt.Println("Sandman Configuration File exists: ", GetSandmanConfigFilename())
	}
}

func LoadSandmanConfig() SandmanConfig {
	var config SandmanConfig
	var config_file_content []byte
	var config_file_path string = GetSandmanConfigFilename()

	config_file_content, err := ioutil.ReadFile(config_file_path)

	if err != nil {
		fmt.Printf("Can't read sandman configuration file at %s", config_file_path)
		fmt.Println("Error: ", err)
		os.Exit(1)
	}

	_, err = toml.Decode(string(config_file_content), &config)

	if err != nil {
		fmt.Printf("Can't decode sandman configuration file at %s", config_file_path)
		fmt.Println("Error: ", err)
		os.Exit(1)
	}

	return config
}

func LoadContainerConfig(container_name string) ContainerConfig {
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

func LoadConfig(container_name string) ContainerConfig {
	var sandmanConfig SandmanConfig = LoadSandmanConfig()
	var containerConfig ContainerConfig = LoadContainerConfig(container_name)

	if err := mergo.Merge(&containerConfig.Build, sandmanConfig.Defaults.Build); err != nil {
		fmt.Printf("Can't decode merge Build configuration file with defaults")
		fmt.Println("Error: ", err)
		os.Exit(1)
	}

	if err := mergo.Merge(&containerConfig.Run, sandmanConfig.Defaults.Run); err != nil {
		fmt.Printf("Can't decode merge Run configuration file with defaults")
		fmt.Println("Error: ", err)
		os.Exit(1)
	}

	return containerConfig
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
