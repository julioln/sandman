package build

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/julioln/sandman/config"
	"github.com/julioln/sandman/podman"

	"github.com/containers/podman/v4/pkg/bindings/images"
	"github.com/containers/podman/v4/pkg/domain/entities"
)

func Build(socket string, containerConfig config.ContainerConfig, layers bool, verbose bool) {
	var conn context.Context = podman.InitializePodman(socket)
	var options entities.BuildOptions

	if verbose {
		fmt.Printf("Container Config: %#v\n", containerConfig)
		fmt.Printf("Connection: %#v\n", conn)
	}

	// Create temporary Dockerfile
	dockerFile, err := ioutil.TempFile("", fmt.Sprintf("sandman_build_%s", containerConfig.Name))
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

	// Build the image
	options.Layers = layers
	options.Output = containerConfig.ImageName
	options.ContextDirectory = containerConfig.Build.ContextDirectory

	if verbose {
		fmt.Printf("Build Options: %#v\n", options)
	}

	buildReport, err := images.Build(conn, []string{dockerFile.Name()}, options)

	if err != nil {
		fmt.Println("Failed to build image")
		fmt.Println("Error: ", err)
		os.Exit(1)
	}

	if verbose {
		fmt.Println("Build report: ", *buildReport)
	}
}

func CmdExecute(socket string, verbose bool, layers bool, args []string) {
	var container_name string = args[0]
	Build(socket, config.LoadConfig(container_name), layers, verbose)
}
