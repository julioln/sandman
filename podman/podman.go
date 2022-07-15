package podman

import (
	"context"
	"fmt"
	"os"

	"github.com/containers/podman/v4/pkg/bindings"
)

func defaultSocket() string {
	var base string
	var has_runtime_dir bool

	if base, has_runtime_dir = os.LookupEnv("XDG_RUNTIME_DIR"); !has_runtime_dir {
		base = fmt.Sprintf("/run/user/%d", os.Getuid())
	}
	return fmt.Sprintf("unix://%s/podman/podman.sock", base)
}

func InitializePodman(socket string) context.Context {
	if socket == "" {
		socket = defaultSocket()
	}

	conn, err := bindings.NewConnection(context.Background(), socket)

	if err != nil {
		fmt.Printf("Can't connect to podman socket at %s. Is it active and running?", socket)
		fmt.Println("Error: ", err)
		os.Exit(1)
	}

	return conn
}
