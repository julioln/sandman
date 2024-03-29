package cmd

import (
	"fmt"

	"github.com/julioln/sandman/build"
	"github.com/julioln/sandman/config"
	"github.com/julioln/sandman/constants"
	"github.com/julioln/sandman/podman"
	"github.com/julioln/sandman/run"

	"github.com/spf13/cobra"
)

var (
	Verbose bool   = false
	Keep    bool   = false
	Layers  bool   = false
	Socket  string = ""

	rootCmd = &cobra.Command{
		Use:     "sandman",
		Short:   "sandman: Sandboxes with Podman",
		Long:    "sandman: Build and run sandboxes with Podman",
		Version: constants.VERSION,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			if Verbose {
				fmt.Println("Args: ", args)
			}
		},
	}

	buildCmd = &cobra.Command{
		Use:     "build [container_name]",
		Short:   "Build an image",
		Long:    "Build an image",
		Aliases: []string{"b"},
		Args:    cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			build.CmdExecute(Socket, Verbose, Layers, args)
		},
	}

	runCmd = &cobra.Command{
		Use:     "run [container_image]",
		Short:   "Starts an attached sandboxed container",
		Long:    "Starts an attached sandboxed container",
		Aliases: []string{"r"},
		Args:    cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			run.CmdExecuteRun(Socket, Verbose, Keep, args)
		},
	}

	startCmd = &cobra.Command{
		Use:     "start [container_image]",
		Short:   "Start a detached sandboxed container",
		Long:    "Start a detached sandboxed container",
		Aliases: []string{"s"},
		Args:    cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			run.CmdExecuteStart(Socket, Verbose, Keep, args)
		},
	}

	scaffoldCmd = &cobra.Command{
		Use:     "sample",
		Short:   "Prints a sample configuration file",
		Long:    "Print a sample configuration file",
		Aliases: []string{"scaffold"},
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Print(config.Scaffold())
		},
	}

	testCmd = &cobra.Command{
		Use:   "test",
		Short: "Tests the connection to the socket",
		Long:  "Tests the connection to the socket",
		Run: func(cmd *cobra.Command, args []string) {
			config.CheckConfig()
			podman.TestConnection(Socket)
		},
	}

	setupCmd = &cobra.Command{
		Use:   "setup",
		Short: "Creates Sandman directories and configuration files",
		Long:  "Creates Sandman directories and configuration files",
		Run: func(cmd *cobra.Command, args []string) {
			config.Setup()
		},
	}
)

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.AddCommand(buildCmd)
	rootCmd.AddCommand(runCmd)
	rootCmd.AddCommand(scaffoldCmd)
	rootCmd.AddCommand(startCmd)
	rootCmd.AddCommand(testCmd)
	rootCmd.AddCommand(setupCmd)

	rootCmd.PersistentFlags().BoolVarP(&Verbose, "verbose", "v", false, "Verbose mode (log debug). Defaults to false")
	rootCmd.PersistentFlags().StringVarP(&Socket, "socket", "", "", fmt.Sprintf("Specify podman socket. Defaults to %s", podman.DefaultSocket()))

	buildCmd.Flags().BoolVarP(&Layers, "layers", "l", false, "Use layers for building (default docker behavior)")
	runCmd.Flags().BoolVarP(&Keep, "keep", "k", false, "Keep container after exit (omit --rm)")
	startCmd.Flags().BoolVarP(&Keep, "keep", "k", false, "Keep container after exit (omit --rm)")
}
