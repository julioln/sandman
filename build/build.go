package build

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/julioln/sandman/config"
	"github.com/julioln/sandman/constants"
	"github.com/julioln/sandman/podman"

	"github.com/containers/buildah/define"
	"github.com/containers/podman/v5/pkg/bindings/images"
	"github.com/containers/podman/v5/pkg/domain/entities"
)

func Build(socket string, containerConfig config.ContainerConfig, layers bool, verbose bool) {
	var conn context.Context = podman.InitializePodman(socket)
	var options entities.BuildOptions
	var commonBuildOptions define.CommonBuildOptions

	if verbose {
		fmt.Printf("Container Config: %#v\n", containerConfig)
		fmt.Printf("Connection: %#v\n", conn)
	}

	// Create temporary Dockerfile
	dockerFile, err := ioutil.TempFile("", fmt.Sprintf("sandman_build_%s", strings.Replace(containerConfig.Name, "/", "_", -1)))
	if err != nil {
		fmt.Println("Failed to write to temp dockerfile")
		fmt.Println("Error: ", err)
		os.Exit(1)
	}
	defer os.Remove(dockerFile.Name())
	defer dockerFile.Close()

	// Write instructions to Dockerfile
	if _, err := dockerFile.Write([]byte(containerConfig.Build.Instructions)); err != nil {
		fmt.Println("Failed to write to temp dockerfile")
		fmt.Println("Error: ", err)
		os.Exit(1)
	}

	// Image paramenters
	options.Layers = layers
	options.Output = containerConfig.ImageName
	if containerConfig.Build.ContextDirectory == "" {
		options.ContextDirectory = config.GetSandmanConfigDir()
	} else {
		options.ContextDirectory = containerConfig.Build.ContextDirectory
	}
	options.AdditionalTags = append(options.AdditionalTags, containerConfig.Build.AdditionalImageNames...)
	options.Labels = append(options.Labels,
		fmt.Sprintf("sandman_version=%s", constants.VERSION),
		fmt.Sprintf("sandman_image_name=%s", containerConfig.ImageName),
		fmt.Sprintf("sandman_container_name=%s", containerConfig.Name),
	)

	// Set building parameters
	commonBuildOptions.Ulimit = containerConfig.Build.Limits.Ulimit
	options.CommonBuildOpts = &commonBuildOptions
	options.Compression = containerConfig.Build.Compression

	if verbose {
		fmt.Printf("Build Options: %#v\n", options)
	}

	buildReport, err := images.Build(conn, []string{dockerFile.Name()}, options)

	if err != nil {
		fmt.Println("Failed to build image")
		fmt.Println("Error: ", err)
	} else if verbose {
		fmt.Println("Build report: ", *buildReport)
	}
}

func CmdExecute(socket string, verbose bool, layers bool, args []string) {
	var container_name string = args[0]
	Build(socket, config.LoadConfig(container_name), layers, verbose)
}
